package exporter

import (
	"fmt"
	"strings"

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
		err := makePathNodesAndParam(root, pathElements, parameter)
		if err != nil {
			return []byte{}, err
		}
	}

	return yaml.Marshal(root)
}

func makePathNodesAndParam(root pathNode, pathElements []string, parameter types.ParameterPackage) error {
	length := len(pathElements)

	if length == 0 {
		panic("makePathNodesAndParam given pathElements with a length of 0.")
	}

	element := pathElements[0]
	if length == 1 { //Base case, retrieve or create leafNode.
		rawParam := (*root)[element]
		if rawParam == nil {
			param := types.NewSerialParameterFromPackage(parameter)
			(*root)[element] = param
			return nil
		} else {
			return fmt.Errorf("Duplicate SSM parameter \"%s\" in export map.", *parameter.Parameter.Name)
		}
	}
	// Navigate down instead
	nextNode := (*root)[element]
	if nextNode == nil {
		newNode := makePathNode()
		(*root)[element] = newNode
		return makePathNodesAndParam(newNode, pathElements[1:], parameter)
	} else {
		castNode := nextNode.(pathNode)
		return makePathNodesAndParam(castNode, pathElements[1:], parameter)
	}
}
