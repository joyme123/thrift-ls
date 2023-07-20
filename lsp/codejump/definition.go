package codejump

import (
	"context"
	"errors"
	"strings"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/lsp/types"
	"github.com/joyme123/thrift-ls/parser"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

var baseType = map[string]struct{}{
	"map":    {},
	"set":    {},
	"list":   {},
	"string": {},
	"i16":    {},
	"i32":    {},
	"i64":    {},
	"i8":     {},
	"double": {},
	"bool":   {},
	"byte":   {},
	"binary": {},
}

func Definition(ctx context.Context, ss *cache.Snapshot, file uri.URI, pos protocol.Position) (res []protocol.Location, err error) {
	res = make([]protocol.Location, 0)
	pf, err := ss.Parse(ctx, file)
	if err != nil {
		return
	}

	if pf.AST() == nil {
		err = errors.New("parse ast failed")
		return
	}

	astPos, err := pf.Mapper().LSPPosToParserPosition(types.Position{Line: pos.Line, Character: pos.Character})
	if err != nil {
		return
	}
	nodePath := parser.SearchNodePath(pf.AST(), astPos)
	targetNode := nodePath[len(nodePath)-1]

	switch targetNode.Type() {
	case "TypeName":
		return fieldTypeDefinition(ctx, ss, file, pf.AST(), nodePath, targetNode)
	case "ConstValue":
		return constValueTypeDefinition(ctx, ss, file, pf.AST(), targetNode)
	}

	return
}

func fieldTypeDefinition(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, nodePath []parser.Node, targetNode parser.Node) ([]protocol.Location, error) {
	res := make([]protocol.Location, 0)

	typeName := targetNode.(*parser.TypeName)
	typeV := typeName.Name
	if _, ok := baseType[typeV]; ok {
		return res, nil
	}

	include, identifier, found := strings.Cut(typeV, ".")
	var astFile uri.URI
	if !found {
		identifier = include
		include = ""
		astFile = file
	} else {
		path := lsputils.GetIncludePath(ast, include)
		if path == "" { // doesn't match any include path
			return nil, nil
		}
		astFile = lsputils.IncludeURI(file, path)
	}

	// now we can find destinate definition in `dstAst` by `identifier`
	dstAst, err := ss.Parse(ctx, astFile)
	if err != nil {
		return nil, err
	}

	// struct, exception, enum or union
	dstException := GetExceptionNode(dstAst.AST(), identifier)
	if dstException != nil {
		res = append(res, jump(astFile, dstException.Name))
	}
	dstStruct := GetStructNode(dstAst.AST(), identifier)
	if dstStruct != nil {
		res = append(res, jump(astFile, dstStruct.Identifier))
	}
	dstEnum := GetEnumNode(dstAst.AST(), identifier)
	if dstEnum != nil {
		res = append(res, jump(astFile, dstEnum.Name))
	}
	dstUnion := GetUnionNode(dstAst.AST(), identifier)
	if dstUnion != nil {
		res = append(res, jump(astFile, dstUnion.Name))
	}
	dstTypedef := GetTypedefNode(dstAst.AST(), identifier)
	if dstTypedef != nil {
		res = append(res, jump(astFile, dstTypedef.Alias))
	}

	return res, nil
}

// search enum
func constValueTypeDefinition(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, targetNode parser.Node) ([]protocol.Location, error) {
	res := make([]protocol.Location, 0)

	constValue := targetNode.(*parser.ConstValue)
	if constValue.TypeName != "identifier" {
		return res, nil
	}

	include, identifier, found := strings.Cut(constValue.Value.(string), ".")
	var astFile uri.URI
	if !found {
		identifier = include
		include = ""
		astFile = file
	} else {
		path := lsputils.GetIncludePath(ast, include)
		if path == "" { // doesn't match any include path, maybe enum value
			include = ""
			identifier = constValue.Value.(string)
			astFile = file
		} else {
			astFile = lsputils.IncludeURI(file, path)
		}
	}

	// check if identifier is enum value
	if !strings.Contains(identifier, ".") {
		return res, nil
	}

	// now we can find destinate definition in `dstAst` by `identifier`
	dstAst, err := ss.Parse(ctx, astFile)
	if err != nil {
		return nil, err
	}

	dstEnumValueIdentifier := GetEnumValueIdentifierNode(dstAst.AST(), identifier)
	if dstEnumValueIdentifier != nil {
		res = append(res, jump(astFile, dstEnumValueIdentifier))
	}

	return res, nil
}
