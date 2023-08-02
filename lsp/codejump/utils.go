package codejump

import (
	"strings"

	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/parser"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func GetExceptionNode(ast *parser.Document, name string) *parser.Exception {
	if ast == nil {
		return nil
	}
	for _, excep := range ast.Exceptions {
		if excep.BadNode || excep.Name == nil {
			continue
		}

		if excep.Name.Name == name {
			return excep
		}
	}

	return nil
}

func GetStructNode(ast *parser.Document, name string) *parser.Struct {
	if ast == nil {
		return nil
	}
	for _, st := range ast.Structs {
		if st.BadNode || st.Identifier == nil {
			continue
		}

		if st.Identifier.Name == name {
			return st
		}
	}

	return nil
}

func GetUnionNode(ast *parser.Document, name string) *parser.Union {
	if ast == nil {
		return nil
	}
	for _, st := range ast.Unions {
		if st.BadNode || st.Name == nil {
			continue
		}

		if st.Name.Name == name {
			return st
		}
	}

	return nil
}

func GetEnumNode(ast *parser.Document, name string) *parser.Enum {
	if ast == nil {
		return nil
	}
	for _, st := range ast.Enums {
		if st.BadNode || st.Name == nil {
			continue
		}

		if st.Name.Name == name {
			return st
		}
	}

	return nil
}

func GetEnumNodeByEnumValue(ast *parser.Document, enumValueName string) *parser.Enum {
	if ast == nil {
		return nil
	}

	enumName, _, found := strings.Cut(enumValueName, ".")
	if !found {
		return nil
	}

	return GetEnumNode(ast, enumName)
}

// GetEnumValueIdentifierNode enum A { ONE }, ONE is the target node
func GetEnumValueIdentifierNode(ast *parser.Document, name string) *parser.Identifier {
	if ast == nil {
		return nil
	}

	enumName, identifier, found := strings.Cut(name, ".")
	if !found {
		return nil
	}

	for _, enum := range ast.Enums {
		if enum.BadNode || enum.Name == nil || enum.Name.Name != enumName {
			continue
		}
		for _, enumValue := range enum.Values {
			if enumValue.Name == nil || enumValue.Name.BadNode || enumValue.Name.Name != identifier {
				continue
			}
			return enumValue.Name
		}
	}

	return nil
}

func GetConstNode(ast *parser.Document, name string) *parser.Const {
	if ast == nil {
		return nil
	}

	for _, cst := range ast.Consts {
		if cst.BadNode || cst.Name == nil || cst.Name.Name != name {
			continue
		}
		return cst
	}

	return nil
}

func GetConstIdentifierNode(ast *parser.Document, name string) *parser.Identifier {
	if ast == nil {
		return nil
	}

	for _, cst := range ast.Consts {
		if cst.BadNode || cst.Name == nil || cst.Name.Name != name {
			continue
		}
		return cst.Name
	}

	return nil
}

func GetTypedefNode(ast *parser.Document, name string) *parser.Typedef {
	if ast == nil {
		return nil
	}
	for _, td := range ast.Typedefs {
		if td.BadNode || td.Alias == nil {
			continue
		}
		if td.Alias.Name == name {
			return td
		}
	}

	return nil
}

func jump(file uri.URI, node parser.Node) protocol.Location {
	rng := lsputils.ASTNodeToRange(node)
	return protocol.Location{
		Range: rng,
		URI:   file,
	}
}
