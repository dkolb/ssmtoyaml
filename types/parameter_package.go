package types

import "github.com/aws/aws-sdk-go-v2/service/ssm/types"

type ParameterPackage struct {
	Parameter types.Parameter
	Metadata  types.ParameterMetadata
}
