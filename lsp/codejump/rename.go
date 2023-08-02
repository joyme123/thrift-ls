package codejump

import (
	"context"
	"errors"
	"fmt"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/lsp/types"
	"github.com/joyme123/thrift-ls/parser"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func PrepareRename(ctx context.Context, ss *cache.Snapshot, file uri.URI, pos protocol.Position) (res protocol.Range, err error) {
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
	case "Identifier":
		return lsputils.ASTNodeToRange(targetNode), nil
	case "ConstValue":
		return
	default:
		err = fmt.Errorf("%s doesn't support rename", targetNode.Type())
		return
	}
}
