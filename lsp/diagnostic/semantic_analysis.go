package diagnostic

import (
	"context"
	"errors"
	"fmt"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/parser"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

type SemanticAnalysis struct {
}

func (s *SemanticAnalysis) Diagnostic(ctx context.Context, ss *cache.Snapshot, changeFiles []uri.URI) (DiagnosticResult, error) {
	res := make(DiagnosticResult)
	for _, file := range changeFiles {
		items, err := s.diagnostic(ctx, ss, file)
		if err != nil {
			return nil, err
		}
		res[file] = items
	}

	return res, nil
}

func (s *SemanticAnalysis) Name() string {
	return "SemanticAnalysis"
}

func (s *SemanticAnalysis) diagnostic(ctx context.Context, ss *cache.Snapshot, changeFile uri.URI) ([]protocol.Diagnostic, error) {
	pf, err := ss.Parse(ctx, changeFile)
	if err != nil {
		return nil, err
	}
	if pf.AST() == nil {
		return nil, errors.New("parse ast failed")
	}

	for _, err := range pf.Errors() {
		log.Debugln("parse err", err)
	}

	res := s.checkDefineConflict(ctx, pf)

	return res, nil
}

func (s *SemanticAnalysis) checkDefineConflict(ctx context.Context, pf *cache.ParsedFile) []protocol.Diagnostic {
	var ret []protocol.Diagnostic

	processStructLike := func(fields []*parser.Field) {
		fieldMap := make(map[string]struct{})
		for i := range fields {
			field := fields[i]
			if field.IsBadNode() || field.ChildrenBadNode() {
				continue
			}
			if _, exist := fieldMap[field.Identifier.Name.Text]; exist {
				// struct conflict
				ret = append(ret, protocol.Diagnostic{
					Range:    lsputils.ASTNodeToRange(field.Identifier.Name),
					Severity: protocol.DiagnosticSeverityError,
					Source:   "thrift-ls",
					Message:  fmt.Sprintf("field name conflict with other field"),
				})
			}
			fieldMap[field.Identifier.Name.Text] = struct{}{}
		}
	}

	definitionNameMap := make(map[string]string)

	structMap := make(map[string]struct{})
	for _, st := range pf.AST().Structs {
		if st.IsBadNode() || st.ChildrenBadNode() {
			continue
		}

		if _, exist := structMap[st.Identifier.Name.Text]; exist {
			// struct conflict
			ret = append(ret, protocol.Diagnostic{
				Range:    lsputils.ASTNodeToRange(st.Identifier.Name),
				Severity: protocol.DiagnosticSeverityError,
				Source:   "thrift-ls",
				Message:  fmt.Sprintf("struct name conflict with other struct"),
			})
		}

		if t, exist := definitionNameMap[st.Identifier.Name.Text]; exist && t != st.Type() {
			// struct conflict
			ret = append(ret, protocol.Diagnostic{
				Range:    lsputils.ASTNodeToRange(st.Identifier.Name),
				Severity: protocol.DiagnosticSeverityHint,
				Source:   "thrift-ls",
				Message:  fmt.Sprintf("struct name conflict with other type"),
			})
		}

		structMap[st.Identifier.Name.Text] = struct{}{}
		definitionNameMap[st.Identifier.Name.Text] = st.Type()

		processStructLike(st.Fields)
	}

	unionMap := make(map[string]struct{})
	for _, union := range pf.AST().Unions {
		if union.IsBadNode() || union.ChildrenBadNode() {
			continue
		}

		if _, exist := unionMap[union.Name.Name.Text]; exist {
			// union conflict
			ret = append(ret, protocol.Diagnostic{
				Range:    lsputils.ASTNodeToRange(union.Name.Name),
				Severity: protocol.DiagnosticSeverityError,
				Source:   "thrift-ls",
				Message:  fmt.Sprintf("union name conflict with other union"),
			})
		}

		if t, exist := definitionNameMap[union.Name.Name.Text]; exist && t != union.Type() {
			// union conflict with others
			ret = append(ret, protocol.Diagnostic{
				Range:    lsputils.ASTNodeToRange(union.Name.Name),
				Severity: protocol.DiagnosticSeverityHint,
				Source:   "thrift-ls",
				Message:  fmt.Sprintf("union name conflict with other type"),
			})
		}

		unionMap[union.Name.Name.Text] = struct{}{}
		definitionNameMap[union.Name.Name.Text] = union.Type()

		processStructLike(union.Fields)
	}

	excepMap := make(map[string]struct{})
	for _, excep := range pf.AST().Exceptions {
		if excep.IsBadNode() || excep.ChildrenBadNode() {
			continue
		}

		if _, exist := excepMap[excep.Name.Name.Text]; exist {
			// exception conflict
			ret = append(ret, protocol.Diagnostic{
				Range:    lsputils.ASTNodeToRange(excep.Name.Name),
				Severity: protocol.DiagnosticSeverityError,
				Source:   "thrift-ls",
				Message:  fmt.Sprintf("exception name conflict with other exception"),
			})
		}

		if t, exist := definitionNameMap[excep.Name.Name.Text]; exist && t != excep.Type() {
			// union conflict with others
			ret = append(ret, protocol.Diagnostic{
				Range:    lsputils.ASTNodeToRange(excep.Name.Name),
				Severity: protocol.DiagnosticSeverityHint,
				Source:   "thrift-ls",
				Message:  fmt.Sprintf("exception name conflict with other type"),
			})
		}

		excepMap[excep.Name.Name.Text] = struct{}{}
		definitionNameMap[excep.Name.Name.Text] = excep.Type()

		processStructLike(excep.Fields)

	}

	svcMap := make(map[string]struct{})
	for _, svc := range pf.AST().Services {
		if svc.IsBadNode() || svc.ChildrenBadNode() {
			continue
		}
		if _, exist := svcMap[svc.Name.Name.Text]; exist {
			// service conflict
			ret = append(ret, protocol.Diagnostic{
				Range:    lsputils.ASTNodeToRange(svc.Name.Name),
				Severity: protocol.DiagnosticSeverityError,
				Source:   "thrift-ls",
				Message:  fmt.Sprintf("service name conflict with other service"),
			})
		}

		if t, exist := definitionNameMap[svc.Name.Name.Text]; exist && t != svc.Type() {
			// service conflict with others
			ret = append(ret, protocol.Diagnostic{
				Range:    lsputils.ASTNodeToRange(svc.Name.Name),
				Severity: protocol.DiagnosticSeverityHint,
				Source:   "thrift-ls",
				Message:  fmt.Sprintf("service name conflict with other type"),
			})
		}

		svcMap[svc.Name.Name.Text] = struct{}{}
		definitionNameMap[svc.Name.Name.Text] = svc.Type()

		fnMap := make(map[string]struct{})
		for _, fn := range svc.Functions {
			if fn.IsBadNode() || svc.ChildrenBadNode() {
				continue
			}
			if _, exist := fnMap[fn.Name.Name.Text]; exist {
				// function conflict
				ret = append(ret, protocol.Diagnostic{
					Range:    lsputils.ASTNodeToRange(fn.Name.Name),
					Severity: protocol.DiagnosticSeverityWarning,
					Source:   "thrift-ls",
					Message:  fmt.Sprintf("function name conflict with other function"),
				})
			}
			fnMap[fn.Name.Name.Text] = struct{}{}
			processStructLike(fn.Arguments)
			if fn.Throws != nil {
				processStructLike(fn.Throws.Fields)
			}
		}
	}

	return ret
}

// same sa goto definition
func (s *SemanticAnalysis) checkTypeExist(ctx context.Context, pf *cache.ParsedFile) []protocol.Diagnostic {
	ret := make([]protocol.Diagnostic, 0)

	return ret
}
