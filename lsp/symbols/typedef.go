package symbols

import (
	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/format"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/parser"
)

func TypedefSymbol(td *parser.Typedef) *protocol.DocumentSymbol {
	res := &protocol.DocumentSymbol{
		Name:           td.Alias.Name.Text,
		Detail:         format.MustFormatFieldType(td.T),
		Kind:           protocol.SymbolKindTypeParameter,
		Range:          lsputils.ASTNodeToRange(td.Alias.Name),
		SelectionRange: lsputils.ASTNodeToRange(td.Alias.Name),
	}

	return res
}
