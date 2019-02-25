package fbx

type FBX struct {
	Header *Header
	Top    *Node
	Nodes  []*Node
}

func (f *FBX) Filter(filter NodeFilter) (nodes []*Node) {
	for _, node := range f.Nodes {
		subNodes := node.Filter(filter)
		nodes = append(nodes, subNodes...)
	}
	return
}
