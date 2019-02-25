package fbx

import (
	"fmt"
	"strings"
)

type Node struct {
	EndOffset       uint32
	NumProperties   uint32
	PropertyListLen uint32
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
