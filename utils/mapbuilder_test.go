package utils_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"gitlab.com/dkub/ssmparams/types"
	"gitlab.com/dkub/ssmparams/utils"
)

func TestMapBuilder(t *testing.T) {
	result, _ := utils.BuildYamlFromParams(generateTestData())
	fmt.Println(string(result))
}

func generateTestData() []types.ParameterPackage {
	return []types.ParameterPackage{
		{
			Parameter: ssmTypes.Parameter{
				ARN:              aws.String("arn:aws:ssm:us-east-2:111122223333:parameter/Application/dev/GithubPassword"),
				DataType:         aws.String("text"),
				LastModifiedDate: &time.Time{},
				Name:             aws.String("/Application/dev/GithubPassword"),
				Type:             ssmTypes.ParameterTypeSecureString,
				Value:            aws.String("Stuff"),
				Version:          1,
			},
			Metadata: ssmTypes.ParameterMetadata{
				Name:             aws.String("/Application/dev/GithubPassword"),
				KeyId:            aws.String("alias/basic-data-symmetric"),
				LastModifiedDate: &time.Time{},
				LastModifiedUser: aws.String("arn:aws:iam::111122223333:root"),
				Description:      aws.String("A description"),
				Version:          *aws.Int64(3),
				Tier:             ssmTypes.ParameterTierStandard,
				Policies:         []ssmTypes.ParameterInlinePolicy{},
				DataType:         aws.String("text"),
			},
		},
		{
			Parameter: ssmTypes.Parameter{
				ARN:              aws.String("arn:aws:ssm:us-east-2:111122223333:parameter/Application/dev/GithubUsername"),
				DataType:         aws.String("text"),
				LastModifiedDate: &time.Time{},
				Name:             aws.String("/Application/dev/GithubUsername"),
				Type:             ssmTypes.ParameterTypeString,
				Value:            aws.String("Stuff"),
				Version:          1,
			},
		},
		{
			Parameter: ssmTypes.Parameter{
				ARN:              aws.String("arn:aws:ssm:us-east-2:111122223333:parameter/Application/prod/GithubPassword"),
				DataType:         aws.String("text"),
				LastModifiedDate: &time.Time{},
				Name:             aws.String("/Application/prod/GithubPassword"),
				Type:             ssmTypes.ParameterTypeSecureString,
				Value:            aws.String("Stuff"),
				Version:          1,
			},
		},
		{
			Parameter: ssmTypes.Parameter{
				ARN:              aws.String("arn:aws:ssm:us-east-2:111122223333:parameter/Application/prod/GithubUsername"),
				DataType:         aws.String("text"),
				LastModifiedDate: &time.Time{},
				Name:             aws.String("/Application/prod/GithubUsername"),
				Type:             ssmTypes.ParameterTypeString,
				Value:            aws.String("Stuff"),
				Version:          1,
			},
			Metadata: ssmTypes.ParameterMetadata{
				Name:             aws.String("/Application/dev/GithubUsername"),
				KeyId:            aws.String("alias/basic-data-symmetric"),
				LastModifiedDate: &time.Time{},
				LastModifiedUser: aws.String("arn:aws:iam::111122223333:root"),
				Description:      aws.String("A description"),
				Version:          *aws.Int64(3),
				Tier:             ssmTypes.ParameterTierStandard,
				Policies:         []ssmTypes.ParameterInlinePolicy{},
				DataType:         aws.String("text"),
			},
		},
		{
			Parameter: ssmTypes.Parameter{
				ARN:              aws.String("arn:aws:ssm:us-east-2:111122223333:parameter/Application/dev/GithubUsername"),
				DataType:         aws.String("text"),
				LastModifiedDate: &time.Time{},
				Name:             aws.String("/AnotherApplication/dev/SomeSetting"),
				Type:             ssmTypes.ParameterTypeString,
				Value:            aws.String("AnotherSetting"),
				Version:          1,
			},
		},
	}
}
