package diagnostic

import (
	"context"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"go.lsp.dev/uri"
)

type SemanticAnalysis struct {
}

func (s *SemanticAnalysis) Diagnostic(ctx context.Context, ss *cache.Snapshot, changeFiles []uri.URI) (DiagnosticResult, error) {

	return nil, nil
}

func (s *SemanticAnalysis) Name() string {
	return "SemanticAnalysis"
}
