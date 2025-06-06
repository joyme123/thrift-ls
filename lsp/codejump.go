package lsp

import (
	"context"

	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/lsp/codejump"
)

func (s *Server) definition(ctx context.Context, params *protocol.DefinitionParams) (result []protocol.Location, err error) {
	file := params.TextDocument.URI
	view, err := s.session.ViewOf(file)
	if err != nil {
		return nil, err
	}
	ss, release := view.Snapshot()
	defer release()

	return codejump.Definition(ctx, ss, params.TextDocument.URI, params.Position)
}

func (s *Server) references(ctx context.Context, params *protocol.ReferenceParams) (result []protocol.Location, err error) {
	file := params.TextDocument.URI
	view, err := s.session.ViewOf(file)
	if err != nil {
		return nil, err
	}
	ss, release := view.Snapshot()
	defer release()

	return codejump.Reference(ctx, ss, params.TextDocument.URI, params.Position)
}

func (s *Server) typeDefinition(ctx context.Context, params *protocol.TypeDefinitionParams) (result []protocol.Location, err error) {
	file := params.TextDocument.URI
	view, err := s.session.ViewOf(file)
	if err != nil {
		return nil, err
	}
	ss, release := view.Snapshot()
	defer release()

	return codejump.TypeDefinition(ctx, ss, params.TextDocument.URI, params.Position)
}
