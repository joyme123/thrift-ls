package completion

import (
	"github.com/joyme123/thrift-ls/lsp/cache"
	"go.lsp.dev/protocol"
)

type Position struct {
	Line      uint32
	Character uint32
}

type CompletionRequest struct {
	TriggerKind int
	Pos         Position
	Fh          cache.FileHandle
}

type CompletionItem struct {
	// Label holds the primary text user sees
	Label string

	Detail string

	// InsertText holds the text to insert when user selects this completion.
	// It may be same with Label
	InsertText string

	Kind       protocol.CompletionItemKind
	Deprecated bool

	Score int

	// Documentation holds document text for this completion
	Documentation string
}
