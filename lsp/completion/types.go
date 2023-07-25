package completion

import (
	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/types"
	"go.lsp.dev/protocol"
)

type CompletionRequest struct {
	TriggerKind int
	Pos         types.Position
	Fh          cache.FileHandle
}

type CompletionItem struct {
	// Label holds the primary text user sees
	Label string

	// Detail a human-readable string with additional information
	// about this item, like type or symbol information.
	Detail string

	// InsertText holds the text to insert when user selects this completion.
	// It may be same with Label
	InsertText       string
	InsertTextFormat protocol.InsertTextFormat

	Kind       protocol.CompletionItemKind
	Deprecated bool

	Score int

	// Documentation holds document text for this completion
	Documentation string
}
