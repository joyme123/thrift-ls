package completion

import (
	"context"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"go.lsp.dev/protocol"
)

type Interface interface {
	Completion(ctx context.Context, ss *cache.Snapshot, cmp *CompletionRequest) ([]*CompletionItem, error)
}

// SemanticBasedCompletion generates completion list based on semantic. It is more precisely than token based completion
type SemanticBasedCompletion struct {
}

func (c *SemanticBasedCompletion) Completion(ctx context.Context, ss *cache.Snapshot, cmp *CompletionRequest) ([]*CompletionItem, error) {
	return []*CompletionItem{
		{
			Label:         "test",
			Detail:        "test",
			InsertText:    "test",
			Kind:          protocol.CompletionItemKindText,
			Deprecated:    false,
			Score:         90,
			Documentation: "test doc",
		},
		{
			Label:         "test2",
			Detail:        "test2",
			InsertText:    "test2",
			Kind:          protocol.CompletionItemKindText,
			Deprecated:    false,
			Score:         80,
			Documentation: "test2 doc",
		},
	}, nil
}

func BuildCompletionItem(candidate Candidate) *CompletionItem {
	return &CompletionItem{
		Label:            candidate.text,
		Detail:           candidate.text,
		InsertText:       candidate.text,
		InsertTextFormat: candidate.format,
		Kind:             protocol.CompletionItemKindText,
		Deprecated:       false,
		Score:            90,
		Documentation:    "",
	}
}
