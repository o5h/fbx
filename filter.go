package fbx

type NodeFilter func(*Node) bool

func FilterName(name string) NodeFilter {
	return func(n *Node) bool {
		return n.Name == name
	}
}
