package utils

import (
	"fmt"
	"strings"

	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"gitlab.com/dkub/ssmparams/types"
	"gopkg.in/yaml.v3"
)

type pathNode *map[string]interface{}

func makePathNode() pathNode {
	node := make(map[string]interface{})
	return &node
}

func BuildYamlFromParams(parameters []types.ParameterPackage) ([]byte, error) {
	root := makePathNode()

	for _, parameter := range parameters {
		pathElements := strings.Split(*parameter.Parameter.Name, "/")
		if pathElements[0] == "" {
			pathElements = pathElements[1:]
		}
		fmt.Println("Got path elements:", pathElements)
		param := makeOrFindParam(root, pathElements)
		(*param).Value = *parameter.Parameter.Value
		(*param).Type = parameter.Parameter.Type
		if (*param).Type == ssmTypes.ParameterTypeSecureString && parameter.Metadata.KeyId != nil {
			(*param).Key = *parameter.Metadata.KeyId
		}
	}

	d, err := yaml.Marshal(root)
	return string(d), err
}

func makeOrFindParam(root pathNode, pathElements []string) *types.SerialParameter {
	length := len(pathElements)

	if length == 0 {
		panic("makeOrFindLeaf given pathElements with a length of 0.")
	}

	element := pathElements[0]
	if length == 1 { //Base case, retrieve or create leafNode.
		rawParam := (*root)[element]
		if rawParam == nil {
			param := &types.SerialParameter{}
			(*root)[element] = param
			return param
		} else {
			param := rawParam.(*types.SerialParameter)
			return param
		}
	}
	// Navigate down instead
	nextNode := (*root)[element]
	if nextNode == nil {
		newNode := makePathNode()
		(*root)[element] = newNode
		return makeOrFindParam(newNode, pathElements[1:])
	} else {
		castNode := nextNode.(pathNode)
		return makeOrFindParam(castNode, pathElements[1:])
	}
}
