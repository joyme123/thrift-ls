package symbols

import (
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/parser"
	"go.lsp.dev/protocol"
)

func ConstSymbol(cst *parser.Const) *protocol.DocumentSymbol {
	if cst.IsBadNode() || cst.ChildrenBadNode() {
		return nil
	}

	res := &protocol.DocumentSymbol{
		Name:           cst.Name.Name.Text,
		Kind:           protocol.SymbolKindConstant,
		Range:          lsputils.ASTNodeToRange(cst),
		SelectionRange: lsputils.ASTNodeToRange(cst),
	}

	return res
}
