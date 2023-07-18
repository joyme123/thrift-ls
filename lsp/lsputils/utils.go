package lsputils

import (
	"github.com/joyme123/thrift-ls/parser"
	"go.lsp.dev/protocol"
)

func ASTNodeToRange(node parser.Node) protocol.Range {
	return protocol.Range{
		Start: protocol.Position{
			Line:      uint32(node.Pos().Line - 1),
			Character: uint32(node.Pos().Col - 1),
		},
		End: protocol.Position{
			Line:      uint32(node.End().Line - 1),
			Character: uint32(node.End().Col - 1),
		},
	}
}
