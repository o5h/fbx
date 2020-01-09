package fbx

import (
	"fmt"
	"strings"
)

type Node struct {
	EndOffset       uint64
	NumProperties   uint64
	PropertyListLen uint64
	NameLen         uint8
	Name            string
	Properties      []*Property
	NestedNodes     []*Node
}

func (node *Node) IsEmpty() bool {
	if node.EndOffset == 0 &&
		node.NumProperties == 0 &&
		node.PropertyListLen == 0 &&
		node.NameLen == 0 {
		return true
	}
	return false
}

func (node *Node) FilterName(name string) []*Node {
	return node.Filter(FilterName(name))
}

func (node *Node) Filter(f NodeFilter) (nodes []*Node) {
	if f(node) {
		nodes = append(nodes, node)
	}
	for _, sub := range node.NestedNodes {
		subNodes := sub.Filter(f)
		nodes = append(nodes, subNodes...)
	}
	return
}

func (node *Node) Int32Slice(name string) ([]int32, bool) {
	nodes := node.FilterName(name)
	if len(nodes) != 1 {
		return nil, false
	}
	properties := nodes[0].Properties
	if len(properties) != 1 {
		return nil, false
	}
	return properties[0].AsInt32Slice()
}

func (node *Node) Float64Slice(name string) ([]float64, bool) {
	nodes := node.FilterName(name)
	if len(nodes) != 1 {
		return nil, false
	}
	properties := nodes[0].Properties
	if len(properties) != 1 {
		return nil, false
	}
	return properties[0].AsFloat64Slice()
}

func (n *Node) String() string {
	b := strings.Builder{}
	b.WriteString(n.Name)
	b.WriteString(":")
	if len(n.Properties) > 0 {
		b.WriteString(fmt.Sprint("", n.Properties, ""))
	}
	if len(n.NestedNodes) > 0 {
		b.WriteString(fmt.Sprint("{", n.NestedNodes, "}"))
	}
	b.WriteString("\n")
	return b.String()
}
