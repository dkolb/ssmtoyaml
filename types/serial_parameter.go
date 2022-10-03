package types

import (
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type SerialParameter struct {
	Value string
	Type  ssmTypes.ParameterType
	Key   string            `yaml:",omitempty"`
	Tags  map[string]string `yaml:",omitempty"`
}

func IsSerialParameter(theMap map[string]interface{}) bool {
	_, hasValueKey := theMap["value"]
	_, hasTypeKey := theMap["type"]
	return hasValueKey && hasTypeKey
}

func NewSerialParameterFromMap(theMap map[string]interface{}) *SerialParameter {
	key, hasKey := theMap["key"]
	tags, hasTags := theMap["tags"]

	returnVal := SerialParameter{
		Value: theMap["value"].(string),
		Type:  ssmTypes.ParameterType(theMap["type"].(string)),
	}

	if hasKey {
		returnVal.Key = key.(string)
	}

	if hasTags {
		returnVal.Tags = tags.(map[string]string)
	}

	return &returnVal
}

func NewSerialParameterFromPackage(parameter AwsParameterPackage) *SerialParameter {
	param := SerialParameter{
		Value: *parameter.Parameter.Value,
		Type:  parameter.Parameter.Type,
	}
	if param.Type == ssmTypes.ParameterTypeSecureString && parameter.Metadata.KeyId != nil {
		param.Key = *parameter.Metadata.KeyId
	}
	if len(parameter.Tags) > 0 {
		tagMap := make(map[string]string)
		for _, tag := range parameter.Tags {
			tagMap[*tag.Key] = *tag.Value
		}
		param.Tags = tagMap
	}
	return &param
}
