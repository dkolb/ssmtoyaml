package types

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type AwsParameterPackage struct {
	Parameter types.Parameter
	Metadata  types.ParameterMetadata
	Tags      []types.Tag
}

func NewAwsParameterPackage() *AwsParameterPackage {
	return &AwsParameterPackage{}
}

func (pkg *AwsParameterPackage) GetPathElements() []string {
	pathElements := strings.Split(*pkg.Parameter.Name, "/")
	if len(pathElements) > 0 && pathElements[0] == "" {
		return pathElements[1:]
	} else {
		return pathElements
	}
}
