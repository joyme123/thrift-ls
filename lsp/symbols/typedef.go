package symbols

import (
	"github.com/joyme123/thrift-ls/format"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/parser"
	"go.lsp.dev/protocol"
)

func TypedefSymbol(td *parser.Typedef) *protocol.DocumentSymbol {
	res := &protocol.DocumentSymbol{
		Name:           td.Alias.Name.Text,
		Detail:         format.MustFormatFieldType(td.T),
		Kind:           protocol.SymbolKindTypeParameter,
		Range:          lsputils.ASTNodeToRange(td),
		SelectionRange: lsputils.ASTNodeToRange(td),
	}

	return res
}
