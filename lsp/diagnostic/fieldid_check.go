package diagnostic

import (
	"context"
	"errors"
	"fmt"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/parser"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

type FieldIDCheck struct {
}

// FieldIDCheck checks struct, union, exception, function pramas, function throws field id
// field id format: have a unique, positive integer identifier.
// ref doc: http://diwakergupta.github.io/thrift-missing-guide/#_defining_structs
func (c *FieldIDCheck) Diagnostic(ctx context.Context, ss *cache.Snapshot, changeFiles []uri.URI) (DiagnosticResult, error) {
	res := make(DiagnosticResult)
	for _, file := range changeFiles {
		items, err := c.diagnostic(ctx, ss, file)
		if err != nil {
			return nil, err
		}
		res[file] = items
	}

	return res, nil
}

func (c *FieldIDCheck) diagnostic(ctx context.Context, ss *cache.Snapshot, file uri.URI) ([]protocol.Diagnostic, error) {
	pf, err := ss.Parse(ctx, file)
	if err != nil {
		return nil, err
	}

	if pf.AST() == nil {
		return nil, errors.New("parse ast failed")
	}

	for _, err := range pf.Errors() {
		fmt.Println("parse err", err)
	}

	var ret []protocol.Diagnostic

	processStructLike := func(fields []*parser.Field) {
		fieldIDSet := make(map[int][]*parser.Field)
		for i := range fields {
			field := fields[i]
			if field.Index.BadNode {
				continue
			}
			fieldIDSet[field.Index.Value] = append(fieldIDSet[field.Index.Value], field)
		}

		for fieldID, set := range fieldIDSet {
			if fieldID < 1 || fieldID > 32767 {
				for _, field := range set {
					// field ID exceeded
					ret = append(ret, protocol.Diagnostic{
						Range:    lsputils.ASTNodeToRange(field.Index),
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  fmt.Sprintf("field id should be a positive integer in [1, 32767]"),
					})
				}
			}

			if len(set) == 1 {
				continue
			}
			for _, field := range set {
				// field id conflict
				ret = append(ret, protocol.Diagnostic{
					Range:    lsputils.ASTNodeToRange(field.Index),
					Severity: protocol.DiagnosticSeverityError,
					Source:   "thrift-ls",
					Message:  fmt.Sprintf("field id conflict"),
				})
			}
		}
	}

	for _, st := range pf.AST().Structs {
		processStructLike(st.Fields)
	}

	for _, union := range pf.AST().Unions {
		processStructLike(union.Fields)
	}

	for _, excep := range pf.AST().Exceptions {

		processStructLike(excep.Fields)

	}

	for _, svc := range pf.AST().Services {
		for _, fn := range svc.Functions {
			processStructLike(fn.Arguments)
			processStructLike(fn.Throws)
		}
	}

	return ret, nil
}
