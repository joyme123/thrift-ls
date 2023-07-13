package parser

func SearchNodePath(root Node, pos Position) []Node {
	path := make([]Node, 0)
	searchNodePath(root, pos, &path)

	return path
}

func searchNodePath(root Node, pos Position, path *[]Node) {
	if !root.Contains(pos) {
		return
	}

	*path = append(*path, root)

	for _, child := range root.Children() {
		searchNodePath(child, pos, path)
	}
}
