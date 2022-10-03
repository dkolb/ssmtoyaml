package types

import (
	"fmt"

	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"gopkg.in/yaml.v3"
)

// ParameterTree functions as the root of our YAML document format.
type ParameterTree struct {
	root *ParameterTreeNode
}

// Explicit interface implementation
var _ yaml.Marshaler = (*ParameterTree)(nil)

func NewParameterTree() *ParameterTree {
	return &ParameterTree{
		root: NewParameterTreeNodePath(),
	}
}

func (p *ParameterTree) AddParamFromPackage(pkg AwsParameterPackage) {
	pathElements := pkg.GetPathElements()
	value := NewParameterTreeValueFromPackage(pkg)
	p.root.recursiveAddValueNode(pathElements, value)
}

func (p *ParameterTree) MarshalYAML() (interface{}, error) {
	return p.root.MarshalYAML()
}

type ParameterTreeTags map[string]string

type ParameterTreeValue struct {
	Type  ssmTypes.ParameterType `yaml:"_type"`
	Value string                 `yaml:"_value"`
	Key   string                 `yaml:"_key,omitempty"`
	Tags  ParameterTreeTags      `yaml:"_tags,omitempty"`
}

func NewParameterTreeValueFromPackage(pkg AwsParameterPackage) *ParameterTreeValue {
	paramValue := ParameterTreeValue{
		Value: *pkg.Parameter.Value,
		Type:  pkg.Parameter.Type,
	}
	if paramValue.Type == ssmTypes.ParameterTypeSecureString && pkg.Metadata.KeyId != nil {
		paramValue.Key = *pkg.Metadata.KeyId
	}
	if len(pkg.Tags) > 0 {
		tagMap := make(map[string]string)
		for _, tag := range pkg.Tags {
			tagMap[*tag.Key] = *tag.Value
		}
		paramValue.Tags = tagMap
	}
	return &paramValue
}

type ParameterTreeNodeType int

const (
	ParameterTreeNodeTypePath ParameterTreeNodeType = iota
	ParameterTreeNodeTypeValue
)

type ParameterTreeNode struct {
	nodeType ParameterTreeNodeType
	value    *ParameterTreeValue
	children map[string]*ParameterTreeNode
}

// Explicit interface implementation check
var _ yaml.Marshaler = (*ParameterTreeNode)(nil)

func NewParameterTreeNodeValue(value *ParameterTreeValue) *ParameterTreeNode {
	return &ParameterTreeNode{
		nodeType: ParameterTreeNodeTypeValue,
		value:    value,
	}
}

func NewParameterTreeNodePath() *ParameterTreeNode {
	return &ParameterTreeNode{
		nodeType: ParameterTreeNodeTypePath,
		children: make(map[string]*ParameterTreeNode),
	}
}

func (n *ParameterTree) GetRoot() *ParameterTreeNode {
	return n.root
}

func (n *ParameterTreeNode) GetNodeType() ParameterTreeNodeType {
	return n.nodeType
}

func (n *ParameterTreeNode) GetValue() (*ParameterTreeValue, error) {
	if n.nodeType == ParameterTreeNodeTypeValue {
		return n.value, nil
	} else {
		return nil, fmt.Errorf("called getValue on a node that isn't a Value node")
	}
}

func (n *ParameterTreeNode) SetValue(value *ParameterTreeValue) {
	n.value = value
}

func (n *ParameterTreeNode) GetChildNode(name string) (*ParameterTreeNode, bool, error) {
	if n.nodeType != ParameterTreeNodeTypePath {
		return nil, false, fmt.Errorf("node with value cannot have children path nodes")
	}
	child, hasChild := n.children[name]
	return child, hasChild, nil
}

func (n *ParameterTreeNode) AddChildPathNode(name string) (*ParameterTreeNode, error) {
	err := n.canAddChild(name)
	if err != nil {
		return nil, err
	}
	newChild := NewParameterTreeNodePath()
	n.children[name] = newChild
	return newChild, nil
}

func (n *ParameterTreeNode) AddChildValueNode(name string, value *ParameterTreeValue) (*ParameterTreeNode, error) {
	err := n.canAddChild(name)
	if err != nil {
		return nil, err
	}
	newChild := NewParameterTreeNodeValue(value)
	n.children[name] = newChild
	return newChild, nil
}

func (n *ParameterTreeNode) canAddChild(name string) error {
	if n.nodeType != ParameterTreeNodeTypePath {
		return fmt.Errorf("node with value cannot have children path nodes")
	}
	_, hasChild := n.children[name]
	if hasChild {
		return fmt.Errorf("node already has child '%s'", name)
	}
	return nil
}

func (n *ParameterTreeNode) recursiveAddValueNode(nodeNames []string, value *ParameterTreeValue) error {
	nodeNameLength := len(nodeNames)
	if nodeNameLength < 0 {
		//"Impossible" case
		return fmt.Errorf("nodeName array 0 length")
	} else if nodeNameLength == 1 {
		//Base case
		n.AddChildValueNode(nodeNames[0], value)
		return nil //exit
	} else {
		//Navigate down one node. Create if necessary.
		child, hasChild, err := n.GetChildNode(nodeNames[0])
		if err != nil {
			//This is not a node that can have children. Panick!
			return err
		} else if hasChild {
			//Child path node exists.
			return child.recursiveAddValueNode(nodeNames[1:], value)
		} else {
			//Create child path node.
			node, err := n.AddChildPathNode(nodeNames[0])
			if err != nil {
				//Massive wtf moment if you hit this.
				return err
			} else {
				//Navigate to new node.
				return node.recursiveAddValueNode(nodeNames[1:], value)
			}
		}
	}
}

func (n *ParameterTreeNode) MarshalYAML() (interface{}, error) {
	if n.nodeType == ParameterTreeNodeTypePath {
		return n.children, nil
	} else if n.nodeType == ParameterTreeNodeTypeValue {
		return n.value, nil
	} else {
		return nil, nil
	}
}
