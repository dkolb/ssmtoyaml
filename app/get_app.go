package app

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"gopkg.in/yaml.v3"

	"github.com/dkolb/ssmtoyaml/types"
	"github.com/dkolb/ssmtoyaml/utils"
)

type GetApp struct {
	SsmPathRoot    string
	ExportFile     string
	Decrypt        bool
	ForceOverwrite bool
	Region         string
	IgnoreTags     bool
	client         *ssm.Client
}

func (g *GetApp) Init() error {
	config, err := utils.AwsLoadConfig(&g.Region)
	if err != nil {
		log.Printf("Failed to load AWS config: %v", err)
		return err
	}
	g.client = ssm.NewFromConfig(*config)
	return nil
}

func (g *GetApp) Exec() error {
	var err error

	if g.client == nil {
		err = g.Init()
		if err != nil {
			return err
		}
	}

	//Validate file doesn't exist or force overwrite is set
	//so we don't waste time crawling the SSM API.
	if !g.ForceOverwrite {
		_, err = os.Stat(g.ExportFile)
		fileExists := !errors.Is(err, os.ErrNotExist)
		if fileExists {
			return errors.New("cannot overwrite file provide --overwrite")
		}
	}

	//Gather parameters under path.
	var params []types.AwsParameterPackage
	params, err = g.gatherParameters()
	if err != nil {
		return err
	}

	//Generate YAML.
	var yaml []byte
	yaml, err = g.BuildYamlFromParamPackages(params)
	if err != nil {
		return err
	}

	//Write out to file.
	_, err = os.Stat(g.ExportFile)
	fileExists := !errors.Is(err, os.ErrNotExist)
	if fileExists && !g.ForceOverwrite {
		return errors.New("cannot overwrite file provide --overwrite")
	}
	if err = os.WriteFile(g.ExportFile, yaml, 0666); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (g *GetApp) BuildYamlFromParamPackages(params []types.AwsParameterPackage) ([]byte, error) {
	paramTree := types.NewParameterTree()
	for _, param := range params {
		paramTree.AddParamFromPackage(param)
	}
	return yaml.Marshal(paramTree)
}

func (g *GetApp) gatherParameters() ([]types.AwsParameterPackage, error) {
	params := &ssm.GetParametersByPathInput{
		Path:           &g.SsmPathRoot,
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(g.Decrypt),
	}

	paginator := ssm.NewGetParametersByPathPaginator(g.client, params)

	var parameters = []types.AwsParameterPackage{}

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatal("Error paging through parameters", err)
			return nil, err
		}

		for _, param := range output.Parameters {
			pkg, err := g.generatePackage(param)
			if err != nil {
				return parameters, err
			}
			parameters = append(parameters, pkg)
		}
	}
	return parameters, nil
}

func (g *GetApp) generatePackage(parameter ssmTypes.Parameter) (types.AwsParameterPackage, error) {
	pkg := types.AwsParameterPackage{Parameter: parameter}
	var err error
	if parameter.Type == ssmTypes.ParameterTypeSecureString {
		pkg.Metadata, err = g.describeParameter(parameter)
		if err != nil {
			return types.AwsParameterPackage{}, err
		}
	}

	if !g.IgnoreTags {
		pkg.Tags, err = g.getTags(parameter)
	}
	if err != nil {
		return types.AwsParameterPackage{}, err
	}

	return pkg, nil
}

func (g *GetApp) describeParameter(parameter ssmTypes.Parameter) (ssmTypes.ParameterMetadata, error) {
	paramFilter := ssmTypes.ParameterStringFilter{
		Key:    aws.String("Name"),
		Values: []string{*parameter.Name},
	}
	describeResponse, err := g.client.DescribeParameters(
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

func (g *GetApp) getTags(parameter ssmTypes.Parameter) ([]ssmTypes.Tag, error) {
	tagResponse, err := g.client.ListTagsForResource(
		context.TODO(),
		&ssm.ListTagsForResourceInput{
			ResourceId:   parameter.Name,
			ResourceType: ssmTypes.ResourceTypeForTaggingParameter,
		},
	)
	if tagResponse == nil {
		return nil, err
	} else {
		return tagResponse.TagList, err
	}
}
