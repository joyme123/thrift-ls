package symbols

import (
	"strconv"

	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/parser"
)

func EnumSymbol(enum *parser.Enum) *protocol.DocumentSymbol {
	if enum.IsBadNode() || enum.ChildrenBadNode() {
		return nil
	}

	res := &protocol.DocumentSymbol{
		Name:           enum.Name.Name.Text,
		Detail:         "Enum",
		Kind:           protocol.SymbolKindEnum,
		Range:          lsputils.ASTNodeToRange(enum.Name.Name),
		SelectionRange: lsputils.ASTNodeToRange(enum.Name.Name),
	}

	for i := range enum.Values {
		child := EnumValueSymbol(enum.Values[i])
		if child == nil {
			continue
		}
		res.Children = append(res.Children, *child)
	}

	return res

}

func EnumValueSymbol(v *parser.EnumValue) *protocol.DocumentSymbol {
	if v.IsBadNode() || v.ChildrenBadNode() {
		return nil
	}

	res := &protocol.DocumentSymbol{
		Name:           v.Name.Name.Text,
		Detail:         strconv.FormatInt(v.Value, 10),
		Kind:           protocol.SymbolKindNumber,
		Range:          lsputils.ASTNodeToRange(v.Name.Name),
		SelectionRange: lsputils.ASTNodeToRange(v.Name.Name),
	}

	return res
}
