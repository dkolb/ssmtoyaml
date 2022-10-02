package utils

import (
	"reflect"
)

type LeafVisitor func(nodeNames []string, value string) error
type NodeVisitor func(nodeNames []string, node map[string]interface{}) (bool, error)

func WalkMapTree(root map[string]interface{}, visitNode NodeVisitor, visitLeaf LeafVisitor) error {
	for nodeName, nodeChild := range root {
		nodeNames := []string{nodeName}
		err := walkNode(nodeNames, nodeChild, visitNode, visitLeaf)
		if err != nil {
			return err
		}
	}
	return nil
}

func walkNode(nodeNames []string, node interface{}, visitNode NodeVisitor, visitLeaf LeafVisitor) error {
	nodeValue := reflect.ValueOf(node)
	if nodeValue.Kind() == reflect.Map {
		// This node has children.  Visit it, then walk it's children.
		castNode := node.(map[string]interface{})
		descend, err := visitNode(nodeNames, castNode)

		// Check for an err
		if err != nil {
			return err
		}

		if descend {
			for childNodeName, childNode := range castNode {
				childNodeNames := append(nodeNames, childNodeName)
				return walkNode(childNodeNames, childNode, visitNode, visitLeaf)
				if err != nil {
					return err
				}
			}
		}
	} else if nodeValue.Kind() == reflect.String {
		castLeaf := node.(string)
		err := visitLeaf(nodeNames, castLeaf)
		if err != nil {
			return err
		}
	}
	return nil
}
