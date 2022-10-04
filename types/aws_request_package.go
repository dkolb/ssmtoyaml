package types

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type AwsRequestPackage struct {
	exists   bool
	PutParam ssm.PutParameterInput
	AddTags  ssm.AddTagsToResourceInput
}

func (r *AwsRequestPackage) SetExists(exists bool) {
	if exists {
		r.exists = true
		r.PutParam.Overwrite = aws.Bool(true)
	} else {
		r.exists = false
		r.PutParam.Overwrite = aws.Bool(false)
		r.PutParam.Tags = r.AddTags.Tags
	}
}

func (r *AwsRequestPackage) GetExists() bool {
	return r.exists
}

func AwsRequestPackagesFromParameterTree(
	paramTree *ParameterTree,
	globalTags *ParameterTreeTags,
) ([]*AwsRequestPackage, error) {
	requests := []*AwsRequestPackage{}
	overrideTags := false
	var globalTagList []ssmTypes.Tag
	if globalTags != nil {
		overrideTags = true
		globalTagList = tagList(*globalTags)
	}

	err := paramTree.ForAllValues(func(path string, value *ParameterTreeValue) error {
		putReq := ssm.PutParameterInput{
			Name:        aws.String(path),
			Value:       aws.String(value.Value),
			Description: aws.String(value.Description),
			Type:        value.Type,
		}

		if value.Type == ssmTypes.ParameterTypeSecureString {
			putReq.KeyId = aws.String(value.Key)
		}

		var tags []ssmTypes.Tag
		if overrideTags {
			tags = globalTagList
		} else {
			tags = []ssmTypes.Tag{}
			for name, value := range value.Tags {
				tags = append(tags, ssmTypes.Tag{
					Key:   aws.String(name),
					Value: aws.String(value),
				})
			}
		}

		requests = append(requests, &AwsRequestPackage{
			PutParam: putReq,
			AddTags: ssm.AddTagsToResourceInput{
				ResourceId:   aws.String(path),
				ResourceType: ssmTypes.ResourceTypeForTaggingParameter,
				Tags:         tags,
			},
		})
		return nil
	})

	return requests, err
}

func tagList(tagMap ParameterTreeTags) []ssmTypes.Tag {
	tags := []ssmTypes.Tag{}
	for name, value := range tagMap {
		tags = append(tags, ssmTypes.Tag{
			Key:   aws.String(name),
			Value: aws.String(value),
		})
	}
	return tags
}
