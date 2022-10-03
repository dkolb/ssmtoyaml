package app

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"gopkg.in/yaml.v3"

	"gitlab.com/dkub/ssmparams/types"
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
	var params []types.AwsParameterPackage
	params, err = e.gatherParameters()
	if err != nil {
		return err
	}
	var yaml []byte
	yaml, err = e.BuildYamlFromParamPackages(params)
	if err != nil {
		panic(err)
	}
	if err = os.WriteFile(e.ExportFile, yaml, 0666); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (e *ExportApp) BuildYamlFromParamPackages(params []types.AwsParameterPackage) ([]byte, error) {
	paramTree := types.NewParameterTree()
	for _, param := range params {
		paramTree.AddParamFromPackage(param)
	}
	return yaml.Marshal(paramTree)
}

func (e *ExportApp) gatherParameters() ([]types.AwsParameterPackage, error) {
	params := &ssm.GetParametersByPathInput{
		Path:           &e.SsmPathRoot,
		Recursive:      aws.Bool(true),
		WithDecryption: new(bool),
	}

	paginator := ssm.NewGetParametersByPathPaginator(e.client, params)

	var parameters = []types.AwsParameterPackage{}

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

func (e *ExportApp) generatePackage(parameter ssmTypes.Parameter) (types.AwsParameterPackage, error) {
	pkg := types.AwsParameterPackage{Parameter: parameter}
	var err error
	if parameter.Type == ssmTypes.ParameterTypeSecureString {
		pkg.Metadata, err = e.describeParameter(parameter)
		if err != nil {
			return types.AwsParameterPackage{}, err
		}
	}

	pkg.Tags, err = e.getTags(parameter)
	if err != nil {
		return types.AwsParameterPackage{}, err
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

func (e *ExportApp) getTags(parameter ssmTypes.Parameter) ([]ssmTypes.Tag, error) {
	tagResponse, err := e.client.ListTagsForResource(
		context.TODO(),
		&ssm.ListTagsForResourceInput{
			ResourceId:   parameter.Name,
			ResourceType: ssmTypes.ResourceTypeForTaggingParameter,
		},
	)
	return tagResponse.TagList, err
}
