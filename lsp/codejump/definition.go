package codejump

import (
	"context"
	"errors"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/lsp/types"
	"github.com/joyme123/thrift-ls/parser"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

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
	nodePath := parser.SearchNodePathByPosition(pf.AST(), astPos)
	targetNode := nodePath[len(nodePath)-1]

	switch targetNode.Type() {
	case "TypeName":
		return typeNameDefinition(ctx, ss, file, pf.AST(), targetNode)
	case "ConstValue":
		return constValueTypeDefinition(ctx, ss, file, pf.AST(), targetNode)
	case "IdentifierName": // service extends
		return serviceDefinition(ctx, ss, file, pf.AST(), targetNode)
	}

	return
}

func serviceDefinition(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, targetNode parser.Node) ([]protocol.Location, error) {
	res := make([]protocol.Location, 0)
	astFile, id, _, err := ServiceDefinitionIdentifier(ctx, ss, file, ast, targetNode)
	if err != nil {
		return res, err
	}
	if id != nil {
		res = append(res, jump(astFile, id.Name))
	}

	return res, nil
}

func ServiceDefinitionIdentifier(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, targetNode parser.Node) (uri.URI, *parser.Identifier, string, error) {
	identifierName := targetNode.(*parser.IdentifierName)

	include, identifier := lsputils.ParseIdent(file, ast.Includes, identifierName.Text)
	var astFile uri.URI
	if include == "" {
		astFile = file
	} else {
		path := lsputils.GetIncludePath(ast, include)
		if path == "" { // doesn't match any include path
			return "", nil, "", nil
		}
		astFile = lsputils.IncludeURI(file, path)
	}

	// now we can find destinate definition in `dstAst` by `identifier`
	dstAst, err := ss.Parse(ctx, astFile)
	if err != nil {
		return astFile, nil, "", err
	}

	if len(dstAst.Errors()) > 0 {
		log.Errorf("parse error: %v", dstAst.Errors())
	}

	dstService := GetServiceNode(dstAst.AST(), identifier)
	if dstService != nil {
		return astFile, dstService.Name, "Service", nil
	}

	return astFile, nil, "", nil
}

func typeNameDefinition(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, targetNode parser.Node) ([]protocol.Location, error) {
	res := make([]protocol.Location, 0)
	astFile, id, _, err := TypeNameDefinitionIdentifier(ctx, ss, file, ast, targetNode)
	if err != nil {
		return res, err
	}
	if id != nil {
		res = append(res, jump(astFile, id.Name))
	}

	return res, nil
}

func TypeNameDefinitionIdentifier(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, targetNode parser.Node) (uri.URI, *parser.Identifier, string, error) {
	typeName := targetNode.(*parser.TypeName)
	typeV := typeName.Name
	if IsBasicType(typeV) {
		return "", nil, "", nil
	}

	include, identifier := lsputils.ParseIdent(file, ast.Includes, typeV)
	var astFile uri.URI
	if include == "" {
		astFile = file
	} else {
		path := lsputils.GetIncludePath(ast, include)
		if path == "" { // doesn't match any include path
			return "", nil, "", nil
		}
		astFile = lsputils.IncludeURI(file, path)
	}

	// now we can find destinate definition in `dstAst` by `identifier`
	dstAst, err := ss.Parse(ctx, astFile)
	if err != nil {
		return astFile, nil, "", err
	}

	if len(dstAst.Errors()) > 0 {
		log.Errorf("parse error: %v", dstAst.Errors())
	}

	// struct, exception, enum or union
	dstException := GetExceptionNode(dstAst.AST(), identifier)
	if dstException != nil {
		return astFile, dstException.Name, "Exception", nil
	}
	dstStruct := GetStructNode(dstAst.AST(), identifier)
	if dstStruct != nil {
		return astFile, dstStruct.Identifier, "Struct", nil
	}
	dstEnum := GetEnumNode(dstAst.AST(), identifier)
	if dstEnum != nil {
		return astFile, dstEnum.Name, "Enum", nil
	}
	dstUnion := GetUnionNode(dstAst.AST(), identifier)
	if dstUnion != nil {
		return astFile, dstUnion.Name, "Union", nil
	}
	dstTypedef := GetTypedefNode(dstAst.AST(), identifier)
	if dstTypedef != nil {
		return astFile, dstTypedef.Alias, "Typedef", nil
	}

	return astFile, nil, "", nil
}

// search enum
func constValueTypeDefinition(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, targetNode parser.Node) ([]protocol.Location, error) {
	res := make([]protocol.Location, 0)
	astFile, id, err := ConstValueTypeDefinitionIdentifier(ctx, ss, file, ast, targetNode)
	if err != nil {
		return res, err
	}

	if id != nil {
		res = append(res, jump(astFile, id))
	}

	return res, nil
}

func ConstValueTypeDefinitionIdentifier(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, targetNode parser.Node) (uri.URI, *parser.Identifier, error) {
	constValue := targetNode.(*parser.ConstValue)
	if constValue.TypeName != "identifier" {
		return "", nil, nil
	}

	include, identifier := lsputils.ParseIdent(file, ast.Includes, constValue.Value.(string))
	var astFile uri.URI
	if include == "" {
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

	// now we can find destinate definition in `dstAst` by `identifier`
	dstAst, err := ss.Parse(ctx, astFile)
	if err != nil {
		return astFile, nil, err
	}

	dstEnumValueIdentifier := GetEnumValueIdentifierNode(dstAst.AST(), identifier)
	if dstEnumValueIdentifier != nil {
		return astFile, dstEnumValueIdentifier, nil
	}

	constIdentifier := GetConstIdentifierNode(dstAst.AST(), identifier)
	if constIdentifier != nil {
		return astFile, constIdentifier, nil
	}

	return astFile, nil, nil
}
