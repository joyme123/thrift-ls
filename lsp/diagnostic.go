package lsp

import (
	"context"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/diagnostic"
	"github.com/joyme123/thrift-ls/utils/errors"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func (s *Server) diagnostic(ctx context.Context, ss *cache.Snapshot, changeFile *cache.FileChange) error {
	if s.client == nil {
		return nil
	}

	log.Debugln("-----------diagnostic called-----------")
	defer log.Debugln("-----------diagnostic finish-----------")

	diag := diagnostic.NewDiagnostic()
	diagRes, err := diag.Diagnostic(ctx, ss, []uri.URI{changeFile.URI})
	if err != nil {
		log.Errorf("diagnostic failed: %v", err)
	}

	log.Debugln("publish diagnostric result: ", len(diagRes))

	var errs []error
	for file, res := range diagRes {
		if res == nil {
			res = make([]protocol.Diagnostic, 0)
		}
		params := &protocol.PublishDiagnosticsParams{
			URI:         file,
			Diagnostics: res,
		}
		log.Debugln("file:", file, "diagnostics", res)
		err = s.client.PublishDiagnostics(ctx, params)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.NewAggregate(errs)
	}
	return nil
}
