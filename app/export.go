package app

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"gitlab.com/dkub/ssmparams/types"
	"gitlab.com/dkub/ssmparams/utils"
)

type ExportApp struct {
	SsmPathRoot    string
	ExportFile     string
	Decrypt        bool
	ForceOverwrite bool
	Region         string
	client         *ssm.Client
}

func (e *ExportApp) Exec() error {
	config, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(e.Region),
	)
	if err != nil {
		log.Fatal("Error loading AWS config...", err)
	}
	e.client = ssm.NewFromConfig(config)
	var params []types.ParameterPackage
	params, err = e.gatherParameters()
	if err != nil {
		return err
	}
	var yaml []byte
	yaml, err = utils.BuildYamlFromParams(params)
	if err = os.WriteFile(e.ExportFile, yaml, 0666); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (e *ExportApp) gatherParameters() ([]types.ParameterPackage, error) {
	params := &ssm.GetParametersByPathInput{
		Path:           &e.SsmPathRoot,
		Recursive:      aws.Bool(true),
		WithDecryption: new(bool),
	}

	paginator := ssm.NewGetParametersByPathPaginator(e.client, params)

	var parameters = []types.ParameterPackage{}

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatal("Error paging through parameters", err)
			return nil, err
		}

		for _, param := range output.Parameters {
			pkg, err := e.generatePackage(param)
			if err != nil {
				return parameters, err
			}
			parameters = append(parameters, pkg)
		}
	}
	return parameters, nil
}

func (e *ExportApp) generatePackage(parameter ssmTypes.Parameter) (types.ParameterPackage, error) {
	pkg := types.ParameterPackage{Parameter: parameter}
	if parameter.Type == ssmTypes.ParameterTypeSecureString {
		var err error
		pkg.Metadata, err = e.describeParameter(parameter)
		if err != nil {
			return types.ParameterPackage{}, err
		}
	}
	return pkg, nil
}

func (e *ExportApp) describeParameter(parameter ssmTypes.Parameter) (ssmTypes.ParameterMetadata, error) {
	paramFilter := ssmTypes.ParameterStringFilter{
		Key:    aws.String("Name"),
		Values: []string{*parameter.Name},
	}
	describeResponse, err := e.client.DescribeParameters(
		context.TODO(),
		&ssm.DescribeParametersInput{
			ParameterFilters: []ssmTypes.ParameterStringFilter{paramFilter},
		},
	)

	if err != nil {
		log.Fatalf("Failed to describe %s", *parameter.Name)
		return ssmTypes.ParameterMetadata{}, err
	}

	return describeResponse.Parameters[0], nil
}
