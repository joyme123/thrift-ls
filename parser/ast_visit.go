package parser

import (
	"github.com/joyme123/thrift-ls/utils"
)

func SearchNodePathByPosition(root Node, pos Position) []Node {
	path := make([]Node, 0)
	searchNodePath(root, pos, &path)

	return path
}

func searchNodePath(root Node, pos Position, path *[]Node) {
	if utils.IsNil(root) {
		return
	}
	if !root.Contains(pos) {
		return
	}

	*path = append(*path, root)

	for _, child := range root.Children() {
		searchNodePath(child, pos, path)
	}
}
