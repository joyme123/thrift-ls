package symbols

import (
	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/parser"
)

func ConstSymbol(cst *parser.Const) *protocol.DocumentSymbol {
	if cst.IsBadNode() || cst.ChildrenBadNode() {
		return nil
	}

	res := &protocol.DocumentSymbol{
		Name:           cst.Name.Name.Text,
		Kind:           protocol.SymbolKindConstant,
		Range:          lsputils.ASTNodeToRange(cst.Name.Name),
		SelectionRange: lsputils.ASTNodeToRange(cst.Name.Name),
	}

	return res
}
