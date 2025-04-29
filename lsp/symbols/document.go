package symbols

import (
	"context"
	"errors"

	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/lsp/cache"
	"go.lsp.dev/uri"
)

func DocumentSymbols(ctx context.Context, ss *cache.Snapshot, file uri.URI) []*protocol.DocumentSymbol {
	res := make([]*protocol.DocumentSymbol, 0)
	pf, err := ss.Parse(ctx, file)
	if err != nil {
		return res
	}

	if pf.AST() == nil {
		err = errors.New("parse ast failed")
		return res
	}

	doc := pf.AST()

	for i := range doc.Typedefs {
		child := TypedefSymbol(doc.Typedefs[i])
		if child != nil {
			res = append(res, child)
		}
	}

	for i := range doc.Consts {
		child := ConstSymbol(doc.Consts[i])
		if child != nil {
			res = append(res, child)
		}
	}

	for i := range doc.Structs {
		child := StructSymbol(doc.Structs[i])
		if child != nil {
			res = append(res, child)
		}
	}

	for i := range doc.Unions {
		child := UnionSymbol(doc.Unions[i])
		if child != nil {
			res = append(res, child)
		}
	}

	for i := range doc.Exceptions {
		child := ExceptionSymbol(doc.Exceptions[i])
		if child != nil {
			res = append(res, child)
		}
	}

	for i := range doc.Services {
		child := ServiceSymbol(doc.Services[i])
		if child != nil {
			res = append(res, child)
		}
	}

	return res
}
