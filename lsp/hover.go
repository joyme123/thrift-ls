package lsp

import (
	"context"
	"strings"

	"github.com/joyme123/thrift-ls/lsp/codejump"
	"go.lsp.dev/protocol"
)

func (s *Server) hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	file := params.TextDocument.URI
	view, err := s.session.ViewOf(file)
	if err != nil {
		return nil, err
	}
	ss, release := view.Snapshot()
	defer release()

	content, err := codejump.Hover(ctx, ss, params.TextDocument.URI, params.Position)
	if err != nil {
		return nil, err
	}

	if content == "" {
		return nil, nil
	}

	markdown_prefix := "```thrift\n"
	if strings.HasPrefix(content, "\n") {
		markdown_prefix = "```thrift"
	}
	markdown_suffix := "\n```"
	if strings.HasSuffix(content, "\n") {
		markdown_suffix = "```"
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.Markdown,
			Value: markdown_prefix + content + markdown_suffix,
		},
	}, nil
}
