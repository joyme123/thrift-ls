package lsp

import (
	"context"

	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/lsp/symbols"
)

func (s *Server) documentSymbol(ctx context.Context, params *protocol.DocumentSymbolParams) (result []interface{}, err error) {
	file := params.TextDocument.URI
	view, err := s.session.ViewOf(file)
	if err != nil {
		return nil, err
	}
	ss, release := view.Snapshot()
	defer release()

	symbols := symbols.DocumentSymbols(ctx, ss, file)

	for i := range symbols {
		result = append(result, symbols[i])
	}

	return
}
