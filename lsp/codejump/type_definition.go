package codejump

import (
	"context"
	"errors"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/types"
	"github.com/joyme123/thrift-ls/parser"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func TypeDefinition(ctx context.Context, ss *cache.Snapshot, file uri.URI, pos protocol.Position) (res []protocol.Location, err error) {
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
	case "IdentifierName":
		// no parent
		if len(nodePath) <= 2 {
			return res, nil
		}
		parent := nodePath[len(nodePath)-3]
		var fieldType *parser.FieldType
		switch parent.Type() {
		case "Field":
			field := parent.(*parser.Field)
			fieldType = field.FieldType
		case "Typedef":
			typedef := parent.(*parser.Typedef)
			fieldType = typedef.T
		case "Function":
			fn := parent.(*parser.Function)
			fieldType = fn.FunctionType
		case "Const":
			cst := parent.(*parser.Const)
			fieldType = cst.ConstType
		}
		if fieldType != nil && !fieldType.BadNode && fieldType.TypeName != nil {
			return typeNameDefinition(ctx, ss, file, pf.AST(), fieldType.TypeName)
		}
	case "ConstValue":
		return constValueTypeDefinition(ctx, ss, file, pf.AST(), targetNode)
	default:
		log.Warningln("unsupport type for type definition:", targetNode.Type())
	}

	return
}
