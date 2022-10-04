package types

import (
	"fmt"
	"log"
	"strings"

	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

// ParameterTree functions as the root of our YAML document format.
type ParameterTree struct {
	root *ParameterTreeNode
}

// Explicit interface implementation
var _ yaml.Marshaler = (*ParameterTree)(nil)
var _ yaml.Unmarshaler = (*ParameterTree)(nil)

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

func (p *ParameterTree) ForAllValues(
	visit func(path string, value *ParameterTreeValue) error,
) error {
	return p.root.depthFirstWalk([]string{}, visit)
}

func (p *ParameterTree) MarshalYAML() (interface{}, error) {
	return p.root.MarshalYAML()
}

func (p *ParameterTree) UnmarshalYAML(value *yaml.Node) error {
	p.root = NewParameterTreeNodePath()
	if value.Kind == yaml.MappingNode {
		return p.root.UnmarshalYAML(value)
	} else {
		return fmt.Errorf("root node must be !!map")
	}
}

type ParameterTreeTags map[string]string

type ParameterTreeValue struct {
	Type        ssmTypes.ParameterType `yaml:"_type"`
	Value       string                 `yaml:"_value"`
	Key         string                 `yaml:"_key,omitempty"`
	Tags        ParameterTreeTags      `yaml:"_tags,omitempty"`
	Description string                 `yaml:"_desc,omitempty"`
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
var _ yaml.Unmarshaler = (*ParameterTreeNode)(nil)

func NewParameterTreeNodeValue(value *ParameterTreeValue) *ParameterTreeNode {
	return &ParameterTreeNode{
		nodeType: ParameterTreeNodeTypeValue,
		value:    value,
	}
}

func NewParameterTreeNodePath() *ParameterTreeNode {
	return &ParameterTreeNode{
		nodeType: ParameterTreeNodeTypePath,
		children: map[string]*ParameterTreeNode{},
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

func (n *ParameterTreeNode) addBlankChildNode(name string) (*ParameterTreeNode, error) {
	err := n.canAddChild(name)
	if err != nil {
		return nil, err
	}
	newChild := &ParameterTreeNode{
		children: map[string]*ParameterTreeNode{},
	}
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

func (n *ParameterTreeNode) depthFirstWalk(
	pathElements []string,
	visit func(path string, value *ParameterTreeValue) error,
) error {
	if n.nodeType == ParameterTreeNodeTypeValue {
		//Base case
		return visit("/"+strings.Join(pathElements, "/"), n.value)
	} else if n.nodeType == ParameterTreeNodeTypePath {
		var errs error
		for name, child := range n.children {
			err := child.depthFirstWalk(append(pathElements, name), visit)
			multierror.Append(errs, err)
		}
		return errs
	} else {
		return fmt.Errorf("Unknown node type")
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

func (n *ParameterTreeNode) UnmarshalYAML(value *yaml.Node) error {
	// Content is always scalar node then map node for path nodes,
	// or it's an even number of scalars for a value node.
	// We don't have any arrays in our schema.
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("node '%s' is not a map", value.Value)
	}

	// First attempt to decode myself as a ParameterTreeNodeTypeValue.
	if isValueNode(value) {
		n.value = &ParameterTreeValue{}
		err := value.Decode(n.value)
		if err != nil {
			return err
		}
		n.nodeType = ParameterTreeNodeTypeValue
		return nil //We're done here
	}

	// The remaining case is this is a Path node. Work our way through each
	// content Node pair depth first.
	n.nodeType = ParameterTreeNodeTypePath

	var errs error

	for i := 0; i < len(value.Content); i += 2 {
		nameNode := value.Content[i]
		mapNode := value.Content[i+1]

		if nameNode.Kind != yaml.ScalarNode {
			errs = multierror.Append(
				errs,
				fmt.Errorf("node '%s' is not a scalar", nameNode.Value),
			)
		}

		if mapNode.Kind != yaml.MappingNode {
			errs = multierror.Append(
				errs,
				fmt.Errorf("node '%s' is not a map", mapNode.Value),
			)
		}

		if nameNode.Kind == yaml.ScalarNode && mapNode.Kind == yaml.MappingNode {
			// Create a child node
			childNode, err := n.addBlankChildNode(nameNode.Value)
			if err != nil {
				errs = multierror.Append(errs, err)
			}
			//Decode the mapping node into the child node.
			//Essentially calling childNode.UnmarshalYAML(mapNode) but with
			//some extra checking. Might be better to just call directly?
			derr := mapNode.Decode(childNode)
			if derr != nil {
				errs = multierror.Append(errs, derr)
			}
		}
	}
	return errs
}

func isValueNode(mapping *yaml.Node) bool {
	if mapping.Kind != yaml.MappingNode {
		log.Printf("tryValueNode: %s not a mapping node", mapping.Value)
		return false
	}

	if len(mapping.Content)%2 != 0 {
		log.Printf("tryValueNode: %s has an odd number of content nodes", mapping.Value)
		return false
	}

	//Checking both that all names start with '_' and the nodes are the right
	//types at the same time so we loop through one time, short circuiting
	//as soon as we find a problem.
	for i := 0; i < len(mapping.Content); i += 2 {
		name := mapping.Content[i]
		value := mapping.Content[i+1]
		if name.Value == "_tags" && name.Kind == yaml.ScalarNode && value.Kind == yaml.MappingNode {
			continue
		} else if name.Value[0] == '_' && name.Kind == yaml.ScalarNode && value.Kind == yaml.ScalarNode {
			continue
		} else {
			return false
		}
	}
	return true
}
