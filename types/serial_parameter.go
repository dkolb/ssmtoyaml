package types

import (
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type SerialParameter struct {
	Value string
	Type  ssmTypes.ParameterType
	Key   string `yaml:",omitempty"`
}
