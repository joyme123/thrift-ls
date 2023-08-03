package lsp

import (
	"context"

	"github.com/joyme123/thrift-ls/lsp/codejump"
	"go.lsp.dev/protocol"
)

func (s *Server) prepareRename(ctx context.Context, params *protocol.PrepareRenameParams) (*protocol.Range, error) {
	file := params.TextDocument.URI
	view, err := s.session.ViewOf(file)
	if err != nil {
		return nil, err
	}
	ss, release := view.Snapshot()
	defer release()

	return codejump.PrepareRename(ctx, ss, params.TextDocument.URI, params.Position)
}

func (s *Server) rename(ctx context.Context, params *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	file := params.TextDocument.URI
	view, err := s.session.ViewOf(file)
	if err != nil {
		return nil, err
	}
	ss, release := view.Snapshot()
	defer release()

	return codejump.Rename(ctx, ss, params.TextDocument.URI, params.Position, params.NewName)
}
