package symbols

import (
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/parser"
	"go.lsp.dev/protocol"
)

func ServiceSymbol(svc *parser.Service) *protocol.DocumentSymbol {
	if svc.IsBadNode() || svc.ChildrenBadNode() {
		return nil
	}

	res := &protocol.DocumentSymbol{
		Name:           svc.Name.Name.Text,
		Kind:           protocol.SymbolKindInterface,
		Range:          lsputils.ASTNodeToRange(svc.Name.Name),
		SelectionRange: lsputils.ASTNodeToRange(svc.Name.Name),
	}

	for i := range svc.Functions {
		child := FunctionSymbol(svc.Functions[i])
		if child != nil {
			res.Children = append(res.Children, *child)
		}

	}

	return res
}

func FunctionSymbol(fn *parser.Function) *protocol.DocumentSymbol {
	if fn.IsBadNode() || fn.ChildrenBadNode() {
		return nil
	}

	res := &protocol.DocumentSymbol{
		Name:           fn.Name.Name.Text,
		Kind:           protocol.SymbolKindFunction,
		Range:          lsputils.ASTNodeToRange(fn.Name.Name),
		SelectionRange: lsputils.ASTNodeToRange(fn.Name.Name),
	}

	return res
}
