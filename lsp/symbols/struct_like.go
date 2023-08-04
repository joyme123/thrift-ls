package symbols

import (
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/parser"
	"go.lsp.dev/protocol"
)

func StructSymbol(st *parser.Struct) *protocol.DocumentSymbol {
	if st.IsBadNode() || st.ChildrenBadNode() {
		return nil
	}

	res := &protocol.DocumentSymbol{
		Name:           st.Identifier.Name.Text,
		Detail:         "Struct",
		Kind:           protocol.SymbolKindStruct,
		Range:          lsputils.ASTNodeToRange(st),
		SelectionRange: lsputils.ASTNodeToRange(st),
	}

	for i := range st.Fields {
		child := FieldSymbol(st.Fields[i])
		if child == nil {
			continue
		}
		res.Children = append(res.Children, *child)
	}

	return res

}

func UnionSymbol(un *parser.Union) *protocol.DocumentSymbol {
	if un.IsBadNode() || un.ChildrenBadNode() {
		return nil
	}

	res := &protocol.DocumentSymbol{
		Name:           un.Name.Name.Text,
		Detail:         "Union",
		Kind:           protocol.SymbolKindStruct,
		Range:          lsputils.ASTNodeToRange(un),
		SelectionRange: lsputils.ASTNodeToRange(un),
	}

	for i := range un.Fields {
		child := FieldSymbol(un.Fields[i])
		if child == nil {
			continue
		}
		res.Children = append(res.Children, *child)
	}

	return res

}

func ExceptionSymbol(ex *parser.Exception) *protocol.DocumentSymbol {
	if ex.IsBadNode() || ex.ChildrenBadNode() {
		return nil
	}

	res := &protocol.DocumentSymbol{
		Name:           ex.Name.Name.Text,
		Detail:         "Exception",
		Kind:           protocol.SymbolKindStruct,
		Range:          lsputils.ASTNodeToRange(ex),
		SelectionRange: lsputils.ASTNodeToRange(ex),
	}

	for i := range ex.Fields {
		child := FieldSymbol(ex.Fields[i])
		if child == nil {
			continue
		}
		res.Children = append(res.Children, *child)
	}

	return res
}
