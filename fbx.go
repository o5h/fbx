package fbx

import "io"

type FBX struct {
	Header *Header
	Top    *Node
	Nodes  []*Node
}

func (f *FBX) ReadFrom(r io.Reader) (int64, error) {
	reader := &Reader{f, 0, nil}
	return reader.ReadFrom(r)
}

func (f *FBX) GetNode(name string) *Node {
	for _, node := range f.Nodes {
		if node.Name == name {
			return node
		}
	}
	return nil
}

func (f *FBX) Filter(filter NodeFilter) (nodes []*Node) {
	for _, node := range f.Nodes {
		subNodes := node.Filter(filter)
		nodes = append(nodes, subNodes...)
	}
	return
}
