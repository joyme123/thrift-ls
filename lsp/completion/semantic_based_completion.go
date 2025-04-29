package completion

import (
	"context"

	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/lsp/cache"
)

type Interface interface {
	Completion(ctx context.Context, ss *cache.Snapshot, cmp *CompletionRequest) ([]*CompletionItem, protocol.Range, error)
}

// SemanticBasedCompletion generates completion list based on semantic. It is more precisely than token based completion
// TODO(jpf)
type SemanticBasedCompletion struct {
}

func BuildCompletionItem(candidate Candidate) *CompletionItem {
	return &CompletionItem{
		Label:            candidate.showText,
		Detail:           candidate.showText,
		InsertText:       candidate.insertText,
		InsertTextFormat: candidate.format,
		Kind:             protocol.CompletionItemKindText,
		Deprecated:       false,
		Score:            90,
		Documentation:    "",
	}
}
