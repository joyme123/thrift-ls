package completion

import (
	"context"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"go.lsp.dev/protocol"
)

func Completion(ctx context.Context, ss *cache.Snapshot, cmp *CompletionRequest) ([]*CompletionItem, error) {

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
