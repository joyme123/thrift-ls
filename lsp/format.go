package lsp

import (
	"context"

	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/format"
	"github.com/joyme123/thrift-ls/lsp/mapper"
)

func (s *Server) formatting(ctx context.Context, params *protocol.DocumentFormattingParams) (result []protocol.TextEdit, err error) {

	// TODO: 支持 format options

	document := params.TextDocument
	fileURI := document.URI
	view, err := s.session.ViewOf(fileURI)
	if err != nil {
		return nil, err
	}

	ss, release := view.Snapshot()
	defer release()

	fh, err := ss.ReadFile(ctx, fileURI)
	if err != nil {
		return nil, err
	}

	bytes, err := fh.Content()
	if err != nil {
		return nil, err
	}

	pf, err := ss.Parse(ctx, fileURI)
	if err != nil {
		return nil, err
	}
	if len(pf.Errors()) > 0 || pf.AST() == nil {
		return nil, pf.AggregatedError()
	}

	formatted, err := format.FormatDocument(pf.AST())
	if err != nil {
		return nil, err
	}

	mp := mapper.NewMapper(fileURI, bytes)
	endPos := mp.GetLSPEndPosition()
	textEdit := protocol.TextEdit{
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      0,
				Character: 0,
			},
			End: protocol.Position{
				Line:      endPos.Line,
				Character: endPos.Character,
			},
		},
		NewText: formatted,
	}

	result = append(result, textEdit)

	return

}
