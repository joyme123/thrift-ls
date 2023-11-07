package codejump

import (
	"context"
	"errors"
	"strings"

	"github.com/joyme123/thrift-ls/format"
	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/lsp/types"
	"github.com/joyme123/thrift-ls/parser"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func Hover(ctx context.Context, ss *cache.Snapshot, file uri.URI, pos protocol.Position) (res string, err error) {
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

	log.Info("node type:", targetNode.Type())

	switch targetNode.Type() {
	case "TypeName":
		return hoverDefinition(ctx, ss, file, pf.AST(), targetNode)
	case "ConstValue":
		return hoverConstValue(ctx, ss, file, pf.AST(), targetNode)
	case "IdentifierName": // service extends
		return hoverService(ctx, ss, file, pf.AST(), targetNode)
	}

	return
}

func hoverService(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, targetNode parser.Node) (string, error) {
	identifierName := targetNode.(*parser.IdentifierName)
	name := identifierName.Text
	include, identifier, found := strings.Cut(name, ".")
	var astFile uri.URI
	if !found {
		identifier = include
		include = ""
		astFile = file
	} else {
		path := lsputils.GetIncludePath(ast, include)
		if path == "" { // doesn't match any include path
			return "", nil
		}
		astFile = lsputils.IncludeURI(file, path)
	}

	// now we can find destinate definition in `dstAst` by `identifier`
	dstAst, err := ss.Parse(ctx, astFile)
	if err != nil {
		return "", err
	}

	if len(dstAst.Errors()) > 0 {
		log.Errorf("parse error: %v", dstAst.Errors())
	}

	dstService := GetServiceNode(dstAst.AST(), identifier)
	if dstService != nil {
		return format.MustFormatService(dstService), nil
	}

	return "", nil
}

func hoverDefinition(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, targetNode parser.Node) (string, error) {
	typeName := targetNode.(*parser.TypeName)
	typeV := typeName.Name
	if IsBasicType(typeV) {
		return "", nil
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
			return "", nil
		}
		astFile = lsputils.IncludeURI(file, path)
	}

	// now we can find destinate definition in `dstAst` by `identifier`
	dstAst, err := ss.Parse(ctx, astFile)
	if err != nil {
		return "", err
	}

	if len(dstAst.Errors()) > 0 {
		log.Errorf("parse error: %v", dstAst.Errors())
	}

	// struct, exception, enum or union
	dstException := GetExceptionNode(dstAst.AST(), identifier)
	if dstException != nil {
		return format.MustFormatException(dstException), nil
	}
	dstStruct := GetStructNode(dstAst.AST(), identifier)
	if dstStruct != nil {
		return format.MustFormatStruct(dstStruct), nil
	}
	dstEnum := GetEnumNode(dstAst.AST(), identifier)
	if dstEnum != nil {
		return format.MustFormatEnum(dstEnum), nil
	}
	dstUnion := GetUnionNode(dstAst.AST(), identifier)
	if dstUnion != nil {
		return format.MustFormatUnion(dstUnion), nil
	}
	dstTypedef := GetTypedefNode(dstAst.AST(), identifier)
	if dstTypedef != nil {
		return format.MustFormatTypedef(dstTypedef), nil
	}

	return "", nil
}

func hoverConstValue(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, targetNode parser.Node) (string, error) {
	constValue := targetNode.(*parser.ConstValue)
	if constValue.TypeName != "identifier" {
		return "", nil
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

	// now we can find destinate definition in `dstAst` by `identifier`
	dstAst, err := ss.Parse(ctx, astFile)
	if err != nil {
		return "", err
	}

	dstEnum := GetEnumNodeByEnumValue(dstAst.AST(), identifier)
	if dstEnum != nil {
		return format.MustFormatEnum(dstEnum), nil
	}

	dstConst := GetConstNode(dstAst.AST(), identifier)
	if dstConst != nil {
		return format.MustFormatConst(dstConst), nil
	}

	return "", nil
}
