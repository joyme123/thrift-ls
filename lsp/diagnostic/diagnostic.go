package diagnostic

import (
	"context"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/utils/errors"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

var registry []Interface

func init() {
	registry = []Interface{
		&CycleCheck{},
		&Parse{},
		&FieldIDCheck{},
		&SemanticAnalysis{},
	}
}

type Interface interface {
	Diagnostic(ctx context.Context, ss *cache.Snapshot, changeFiles []uri.URI) (DiagnosticResult, error)
	Name() string
}

type Diagnostic struct {
}

func NewDiagnostic() Interface {
	return &Diagnostic{}
}

func (d *Diagnostic) Diagnostic(ctx context.Context, ss *cache.Snapshot, changeFiles []uri.URI) (DiagnosticResult, error) {
	res := make(DiagnosticResult)
	var errs []error
	for _, impl := range registry {
		log.Debugln("diagnostic called: ", impl.Name())
		diagRes, err := impl.Diagnostic(ctx, ss, changeFiles)
		if err != nil {
			errs = append(errs, err)
		}
		for key, items := range diagRes {
			res[key] = append(res[key], items...)
		}
	}
	if len(errs) > 0 {
		return res, errors.NewAggregate(errs)

	}
	return res, nil
}

func (d *Diagnostic) Name() string {
	return "Diagnostic"
}

type DiagnosticResult map[uri.URI][]protocol.Diagnostic
