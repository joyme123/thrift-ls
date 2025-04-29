package diagnostic

import (
	"context"

	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/parser"
	"github.com/joyme123/thrift-ls/utils/errors"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/uri"
)

type Parse struct {
}

func (p *Parse) Diagnostic(ctx context.Context, ss *cache.Snapshot, changeFiles []uri.URI) (DiagnosticResult, error) {
	var errs []error

	res := make(DiagnosticResult)
	for _, uri := range changeFiles {
		parseRes, err := ss.Parse(ctx, uri)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// TODO(jpf): 递归解析
		for _, err := range parseRes.Errors() {
			parseErr, ok := err.(parser.ParserError)
			if !ok {
				continue
			}
			log.Debugf("diagnostic parse err: %v", parseErr)
			diag := parseErrToDiagnostic(parseErr)
			res[uri] = append(res[uri], diag)
		}
	}

	if len(errs) > 0 {
		return res, errors.NewAggregate(errs)
	}

	return res, nil
}

func (p *Parse) Name() string {
	return "Parse"
}

func parseErrToDiagnostic(err parser.ParserError) protocol.Diagnostic {
	line, col, _ := err.Pos()
	diag := protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      uint32(line - 1),
				Character: uint32(col - 1),
			},
			End: protocol.Position{
				Line:      uint32(line - 1),
				Character: uint32(col - 1),
			},
		},
		Severity: protocol.DiagnosticSeverityError,
		Source:   "thrift-ls",
		Message:  err.InnerError().Error(),
	}

	return diag
}
