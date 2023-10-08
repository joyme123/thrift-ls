package diagnostic

import (
	"context"
	"errors"
	"fmt"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/codejump"
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

	items := s.checkDefinitionExist(ctx, ss, changeFile, pf)
	res = append(res, items...)

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

// same as goto definition
// struct/union/exception field type
func (s *SemanticAnalysis) checkDefinitionExist(ctx context.Context, ss *cache.Snapshot, file uri.URI, pf *cache.ParsedFile) []protocol.Diagnostic {
	ret := make([]protocol.Diagnostic, 0)

	// struct/union/exception/function arguments/throw fields field type
	processStructLike := func(fields []*parser.Field) {
		for i := range fields {
			field := fields[i]
			if field.IsBadNode() || field.ChildrenBadNode() {
				continue
			}
			items := s.checkTypeExist(ctx, ss, file, pf, field.FieldType)
			ret = append(ret, items...)

			// default value check
			if field.ConstValue != nil {
				items := s.checkConstValueExist(ctx, ss, file, pf, field.ConstValue)
				ret = append(ret, items...)

				dig := s.checkConstValueMatchType(ctx, field)
				if dig != nil {
					ret = append(ret, *dig)
				}
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

	for _, cst := range pf.AST().Consts {
		items := s.checkConstValueExist(ctx, ss, file, pf, cst.Value)
		ret = append(ret, items...)
	}

	for _, svc := range pf.AST().Services {
		for _, fn := range svc.Functions {
			if fn.FunctionType != nil {
				items := s.checkTypeExist(ctx, ss, file, pf, fn.FunctionType)
				ret = append(ret, items...)
			}

			processStructLike(fn.Arguments)
			if fn.Throws != nil {
				processStructLike(fn.Throws.Fields)
			}
		}
	}

	return ret
}

func (s *SemanticAnalysis) checkConstValueExist(ctx context.Context, ss *cache.Snapshot,
	file uri.URI, pf *cache.ParsedFile, cst *parser.ConstValue) (res []protocol.Diagnostic) {
	if cst.TypeName != "identifier" {
		return
	}

	if cst.Value == "true" || cst.Value == "false" {
		return
	}

	_, id, err := codejump.ConstValueTypeDefinitionIdentifier(ctx, ss, file, pf.AST(), cst)
	if err != nil || id == nil {
		res = append(res, protocol.Diagnostic{
			Range:    lsputils.ASTNodeToRange(cst),
			Severity: protocol.DiagnosticSeverityError,
			Source:   "thrift-ls",
			Message:  fmt.Sprintf("default value doesn't exist"),
		})
	}

	return
}

func (s *SemanticAnalysis) checkConstValueMatchType(ctx context.Context, field *parser.Field) (res *protocol.Diagnostic) {
	if field.BadNode || field.ChildrenBadNode() {
		return nil
	}

	expectTypeName := field.FieldType.TypeName

	if field.ConstValue != nil {
		valueType := field.ConstValue.TypeName
		// TypeName can be: list, map, pair, string, identifier, i64, double
		switch valueType {
		case "list", "map", "string", "double":
			if expectTypeName.Name != valueType {
				return &protocol.Diagnostic{
					Range:    lsputils.ASTNodeToRange(field.ConstValue),
					Severity: protocol.DiagnosticSeverityError,
					Source:   "thrift-ls",
					Message:  fmt.Sprintf("expect %s but got %s", expectTypeName.Name, valueType),
				}
			}
		case "identifier":
			if valueType == "identifier" &&
				(field.ConstValue.Value == "true" || field.ConstValue.Value == "false") {
				valueType = "bool"
			}
			if expectTypeName.Name == "bool" {
				if field.ConstValue.Value != "true" && field.ConstValue.Value != "false" {
					return &protocol.Diagnostic{
						Range:    lsputils.ASTNodeToRange(field.ConstValue),
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  fmt.Sprintf("expect %s but got %s", expectTypeName.Name, valueType),
					}
				}
			} else if codejump.IsBasicType(valueType) {
				return &protocol.Diagnostic{
					Range:    lsputils.ASTNodeToRange(field.ConstValue),
					Severity: protocol.DiagnosticSeverityError,
					Source:   "thrift-ls",
					Message:  fmt.Sprintf("expect %s but got %s", expectTypeName.Name, valueType),
				}
			}
		case "i64":
			if expectTypeName.Name != "i8" &&
				expectTypeName.Name != "i16" &&
				expectTypeName.Name != "i32" &&
				expectTypeName.Name != "i64" {
				return &protocol.Diagnostic{
					Range:    lsputils.ASTNodeToRange(field.ConstValue),
					Severity: protocol.DiagnosticSeverityError,
					Source:   "thrift-ls",
					Message:  fmt.Sprintf("expect %s but got %s", expectTypeName.Name, valueType),
				}
			}
		}
	}

	return nil
}

func (s *SemanticAnalysis) checkTypeExist(ctx context.Context, ss *cache.Snapshot,
	file uri.URI, pf *cache.ParsedFile, ft *parser.FieldType) (res []protocol.Diagnostic) {
	if codejump.IsContainerType(ft.TypeName.Name) {
		return s.checkContainerTypeExist(ctx, ss, file, pf, ft)
	} else if codejump.IsBasicType(ft.TypeName.Name) {
		return nil
	} else {
		_, id, _, err := codejump.TypeNameDefinitionIdentifier(ctx, ss, file, pf.AST(), ft.TypeName)
		if err != nil || id == nil {
			res = append(res, protocol.Diagnostic{
				Range:    lsputils.ASTNodeToRange(ft),
				Severity: protocol.DiagnosticSeverityError,
				Source:   "thrift-ls",
				Message:  fmt.Sprintf("field type doesn't exist"),
			})
		}
	}

	return res
}

func (s *SemanticAnalysis) checkContainerTypeExist(ctx context.Context,
	ss *cache.Snapshot, file uri.URI, pf *cache.ParsedFile, ft *parser.FieldType) (res []protocol.Diagnostic) {

	if ft.KeyType != nil {
		items := s.checkTypeExist(ctx, ss, file, pf, ft.KeyType)
		res = append(res, items...)
	}

	if ft.ValueType != nil {
		items := s.checkTypeExist(ctx, ss, file, pf, ft.ValueType)
		res = append(res, items...)
	}

	return res
}

// TODO(jpf): 类型和默认值的类型要一致
