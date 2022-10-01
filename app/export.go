package app

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type ExportApp struct {
	SsmPathRoot    string
	ExportFile     string
	Decrypt        bool
	ForceOverwrite bool
	Region         string
}

func (e *ExportApp) Exec() error {
	config, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(e.Region),
	)
	if err != nil {
		log.Fatal("Error loading AWS config...", err)
	}
	client := ssm.NewFromConfig(config)

	params := &ssm.GetParametersByPathInput{
		Path:           &e.SsmPathRoot,
		Recursive:      aws.Bool(true),
		WithDecryption: new(bool),
	}

	paginator := ssm.NewGetParametersByPathPaginator(client, params)

	var parameters = []types.Parameter{}

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatal("Error paging through parameters", err)
			return err
		}
		parameters = append(parameters, output.Parameters...)
	}

	return nil
}
