package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aquasecurity/table"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/go-multierror"
	"github.com/liamg/tml"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v3"

	"gitlab.com/dkub/ssmparams/types"
	"gitlab.com/dkub/ssmparams/utils"
)

type PutApp struct {
	File        string
	TagFile     string
	Interactive bool
	RetryLimit  int
	Region      string
	ssm         *ssm.Client
	sts         *sts.Client
	iam         *iam.Client
}

func (im *PutApp) Init() error {
	var (
		conf *aws.Config
		err  error
	)

	conf, err = utils.AwsLoadConfig(&im.Region)
	if err != nil {
		log.Printf("Failed to load AWS default config with %s region: %v", im.Region, err)
		return err
	}

	im.ssm = ssm.NewFromConfig(*conf)
	im.sts = sts.NewFromConfig(*conf)
	im.iam = iam.NewFromConfig(*conf)

	return nil
}

func (im *PutApp) Exec() error {
	var err error

	if im.ssm == nil {
		err = im.Init()
		if err != nil {
			log.Println("failed to init app object")
			return err
		}
	}

	data, err := os.ReadFile(im.File)
	if err != nil {
		log.Printf("failed to open file %s: %v", im.File, err)
		return err
	}

	var paramTree *types.ParameterTree
	paramTree, err = im.UnmarshalData(data)
	if err != nil {
		log.Printf("failed to parse file %s: %v", im.File, err)
		return err
	}

	var tagMap types.ParameterTreeTags = make(map[string]string)
	if len(im.TagFile) > 0 {
		data, err = os.ReadFile(im.TagFile)
		if err != nil {
			log.Printf("failed to open file %s: %v", im.TagFile, err)
			return err
		}
		err = yaml.Unmarshal(data, &tagMap)
		if err != nil {
			log.Printf("failed to parse file %s: %v", im.TagFile, err)
			return err
		}
	}

	var requests []*types.AwsRequestPackage

	requests, err = types.AwsRequestPackagesFromParameterTree(paramTree, nil)

	if err != nil {
		log.Println("failed to convert param tree to requests", err)
		return err
	}

	im.PrepareRequests(requests)
	im.Report(requests)

	if im.Interactive {
		confirmed, err := im.Confirm(requests)
		if err != nil {
			return err
		} else if !confirmed {
			return errors.New("confirmation rejected")
		}
	}

	return im.MakeRequests(requests)
}

func (im *PutApp) UnmarshalData(data []byte) (*types.ParameterTree, error) {
	paramTree := types.NewParameterTree()
	err := yaml.Unmarshal(data, paramTree)
	return paramTree, err
}

func (im *PutApp) PrepareRequests(requestBatch []*types.AwsRequestPackage) {
	for _, request := range requestBatch {
		exists := im.checkParamExists(*request)
		request.SetExists(exists)
	}
}

func (im *PutApp) MakeRequests(requestBatch []*types.AwsRequestPackage) error {
	var errs *multierror.Error
	for _, request := range requestBatch {
		_, err := im.ssm.PutParameter(context.TODO(), request.PutParam)
		errs = multierror.Append(errs, err)
		if err != nil {
			continue //Don't try to add tags to a param that returned an error
		}
		_, err = im.ssm.AddTagsToResource(context.TODO(), request.AddTags)
		errs = multierror.Append(errs, err)
	}
	return errs.ErrorOrNil()
}

func (im *PutApp) Report(requestBatch []*types.AwsRequestPackage) error {
	var (
		callerIdent    *sts.GetCallerIdentityOutput
		accountAliases *iam.ListAccountAliasesOutput
		err            error
	)
	if im.sts == nil || im.iam == nil {
		err = im.Init()
		if err != nil {
			log.Printf("failed to init import app: %v\n", err)
			return err
		}
	}
	callerIdent, err = im.sts.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Printf("failed to get caller identity: %v\n", err)
		return err
	}
	accountAliases, err = im.iam.ListAccountAliases(context.TODO(), &iam.ListAccountAliasesInput{})
	if err != nil {
		log.Printf("failed to list account aliases: %v\n", err)
		return err
	}

	t := table.New(os.Stdout)
	t.SetBorders(true)
	t.SetHeaders("Path", "Type", "Value", "Overwrite?")
	for _, request := range requestBatch {
		var overwrite string
		if request.GetExists() {
			overwrite = tml.Sprintf("<yellow>Yes</yellow>")
		} else {
			overwrite = tml.Sprintf("<green>No</green>")
		}
		t.AddRow(
			*request.PutParam.Name,
			string(request.PutParam.Type),
			*request.PutParam.Value,
			overwrite,
		)
	}
	t.Render()

	fmt.Printf(
		"These operations will be applied to AWS account %s -- %s.\n",
		*callerIdent.Account,
		strings.Join(accountAliases.AccountAliases, ", "),
	)

	return nil
}

func (im *PutApp) Confirm(requestBatch []*types.AwsRequestPackage) (bool, error) {
	prompt := promptui.Select{
		Label: "Proceed with operations? [Yes/No]",
		Items: []string{"No", "Yes"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, fmt.Errorf("prompt run failed with %v", err)
	}
	return result == "Yes", nil
}

func (im *PutApp) checkParamExists(req types.AwsRequestPackage) bool {
	_, err := im.ssm.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name: req.PutParam.Name,
	})

	var nferr *ssmTypes.ParameterNotFound

	return err == nil || !errors.As(err, &nferr)
}
