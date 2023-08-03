package codejump

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/lsp/types"
	"github.com/joyme123/thrift-ls/parser"
	utilerrors "github.com/joyme123/thrift-ls/utils/errors"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

var validReferenceDefinitionType = map[string]struct{}{
	"Struct":    {},
	"Union":     {},
	"Enum":      {},
	"Exception": {},
	"Typedef":   {},
}

func Reference(ctx context.Context, ss *cache.Snapshot, file uri.URI, pos protocol.Position) (res []protocol.Location, err error) {
	res = make([]protocol.Location, 0)
	pf, err := ss.Parse(ctx, file)
	if err != nil {
		return
	}

	if pf.AST() == nil {
		err = errors.New("parse ast failed")
		return
	}

	astPos, err := pf.Mapper().LSPPosToParserPosition(types.Position{Line: pos.Line, Character: pos.Character})
	if err != nil {
		return
	}
	nodePath := parser.SearchNodePath(pf.AST(), astPos)
	targetNode := nodePath[len(nodePath)-1]

	switch targetNode.Type() {
	case "TypeName":
		return searchTypeNameReferences(ctx, ss, file, pf.AST(), nodePath, targetNode)
	case "IdentifierName":
		if len(nodePath) <= 2 {
			return
		}
		// identifierName -> identifier -> definition
		parentDefinitionNode := nodePath[len(nodePath)-3]
		definitionType := parentDefinitionNode.Type()
		if definitionType == "EnumValue" || definitionType == "Const" {
			var typeName string
			if definitionType == "Const" {
				typeName = fmt.Sprintf("%s.%s", lsputils.GetIncludeName(file), targetNode.(*parser.IdentifierName).Text)
			} else {
				enumNode := nodePath[len(nodePath)-4]
				typeName = fmt.Sprintf("%s.%s.%s", lsputils.GetIncludeName(file), enumNode.(*parser.Enum).Name.Name.Text, targetNode.(*parser.IdentifierName).Text)
			}
			// search in const value
			return searchConstValueIdentifierReferences(ctx, ss, file, typeName)
		}

		if _, ok := validReferenceDefinitionType[definitionType]; !ok {
			return
		}

		// typeName is base.User
		typeName := fmt.Sprintf("%s.%s", lsputils.GetIncludeName(file), targetNode.(*parser.IdentifierName).Text)
		return searchIdentifierReferences(ctx, ss, file, typeName, definitionType)
	case "ConstValue":
		return searchConstValueReferences(ctx, ss, file, pf.AST(), nodePath, targetNode)
	default:
		log.Warningln("unsupport type for reference:", targetNode.Type())
	}

	return
}

func searchTypeNameReferences(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, nodePath []parser.Node, targetNode parser.Node) (res []protocol.Location, err error) {
	res = make([]protocol.Location, 0)
	typeNameNode := targetNode.(*parser.TypeName)
	typeName := typeNameNode.Name
	if _, ok := baseType[typeName]; ok {
		return
	}

	var errs []error
	// search type definition
	definitionFile, identifierNode, definitionType, err := typeNameDefinitionIdentifier(ctx, ss, file, ast, nodePath, targetNode)
	if err != nil {
		errs = append(errs, err)
	}

	if identifierNode == nil {
		return
	}
	res = append(res, jump(definitionFile, identifierNode.Name))

	locations, err := searchIdentifierReferences(ctx, ss, definitionFile, typeName, definitionType)
	if err != nil {
		errs = append(errs, err)
	}
	res = append(res, locations...)

	if len(errs) > 0 {
		err = utilerrors.NewAggregate(errs)
	}

	return
}

func searchIdentifierReferences(ctx context.Context, ss *cache.Snapshot, file uri.URI, typeName string, definitionType string) (res []protocol.Location, err error) {
	log.Debugln("searchIdentifierReferences for file:", file, "typeName:", typeName)
	var errs []error

	// search in it self
	locations, err := searchDefinitionIdentifierReferences(ctx, ss, file,
		strings.TrimLeft(typeName, fmt.Sprintf("%s.", lsputils.GetIncludeName(file))), definitionType)
	if err != nil {
		errs = append(errs, err)
	}
	res = append(res, locations...)

	// search type references in other file
	includeNode := ss.Graph().Get(file)
	log.Debugln("includeNode: ", includeNode)
	if includeNode != nil {
		if len(includeNode.InDegree()) == 0 && len(includeNode.OutDegree()) == 0 {
			ss.Graph().Debug()
		}
		referenceFiles := includeNode.InDegree()
		for _, referenceFile := range referenceFiles {
			log.Debugln("reference file: ", referenceFile)
			locations, err := searchDefinitionIdentifierReferences(ctx, ss, referenceFile, typeName, definitionType)
			if err != nil {
				errs = append(errs, err)
			}
			res = append(res, locations...)
		}
	}

	if len(errs) > 0 {
		err = utilerrors.NewAggregate(errs)
	}
	return

}

func searchDefinitionIdentifierReferences(ctx context.Context, ss *cache.Snapshot, file uri.URI, typeName string, definitionType string) (res []protocol.Location, err error) {
	ast, err := ss.Parse(ctx, file)
	if err != nil {
		return
	}

	if ast.AST() == nil {
		return
	}

	var searchFieldType func(fieldType *parser.FieldType)
	searchFieldType = func(fieldType *parser.FieldType) {
		if fieldType.KeyType != nil && !fieldType.KeyType.BadNode {
			searchFieldType(fieldType.KeyType)
		}
		if fieldType.ValueType != nil && !fieldType.ValueType.BadNode {
			searchFieldType(fieldType.ValueType)
		}

		if fieldType.TypeName != nil && fieldType.TypeName.Name == typeName {
			res = append(res, jump(file, fieldType.TypeName))
		}
	}

	jumpField := func(field *parser.Field) {
		if field.BadNode || field.FieldType == nil || field.FieldType.BadNode || field.FieldType.TypeName == nil {
			return
		}
		searchFieldType(field.FieldType)
	}

	for _, svc := range ast.AST().Services {
		for _, fn := range svc.Functions {
			if fn.BadNode {
				continue
			}

			if fn.FunctionType != nil && !fn.FunctionType.BadNode {
				searchFieldType(fn.FunctionType)
			}

			for i := range fn.Arguments {
				jumpField(fn.Arguments[i])
			}

			if fn.Throws != nil {
				for i := range fn.Throws.Fields {
					jumpField(fn.Throws.Fields[i])
				}
			}
		}
	}
	if definitionType == "Exception" {
		return
	}

	for _, st := range ast.AST().Structs {
		if st.BadNode {
			continue
		}

		for _, field := range st.Fields {
			jumpField(field)
		}
	}

	for _, st := range ast.AST().Unions {
		if st.BadNode {
			continue
		}

		for _, field := range st.Fields {
			jumpField(field)
		}
	}

	for _, st := range ast.AST().Exceptions {
		if st.BadNode {
			continue
		}

		for _, field := range st.Fields {
			jumpField(field)
		}
	}

	for _, typedef := range ast.AST().Typedefs {
		if typedef.BadNode || typedef.T == nil || typedef.T.BadNode {
			continue
		}
		searchFieldType(typedef.T)
	}

	for _, cst := range ast.AST().Consts {
		if cst.BadNode || cst.ConstType == nil || cst.ConstType.BadNode {
			continue
		}
		searchFieldType(cst.ConstType)
	}

	return res, nil
}

func searchConstValueReferences(ctx context.Context, ss *cache.Snapshot, file uri.URI, ast *parser.Document, nodePath []parser.Node, targetNode parser.Node) (res []protocol.Location, err error) {
	res = make([]protocol.Location, 0)

	var errs []error
	// search type definition
	definitionFile, identifierNode, err := constValueTypeDefinitionIdentifier(ctx, ss, file, ast, targetNode)
	if err != nil {
		errs = append(errs, err)
	}

	if identifierNode == nil {
		return
	}
	res = append(res, jump(definitionFile, identifierNode.Name))

	valueName := targetNode.(*parser.ConstValue).Value.(string)
	locations, err := searchConstValueIdentifierReferences(ctx, ss, definitionFile, valueName)
	if err != nil {
		errs = append(errs, err)
	}
	res = append(res, locations...)

	if len(errs) > 0 {
		err = utilerrors.NewAggregate(errs)
	}

	return
}

// a const value maybe a const defintion or enum value definition
func searchConstValueIdentifierReferences(ctx context.Context, ss *cache.Snapshot, file uri.URI, valueName string) (res []protocol.Location, err error) {
	var errs []error
	// search in it self
	locations, err := searchConstValueIdentifierReference(ctx, ss, file, strings.TrimLeft(valueName, fmt.Sprintf("%s.", lsputils.GetIncludeName(file))))
	if err != nil {
		errs = append(errs, err)
	}
	res = append(res, locations...)

	// search type references in other file
	includeNode := ss.Graph().Get(file)
	if includeNode != nil {
		referenceFiles := includeNode.InDegree()
		for _, referenceFile := range referenceFiles {
			locations, err := searchConstValueIdentifierReference(ctx, ss, referenceFile, valueName)
			if err != nil {
				errs = append(errs, err)
			}
			res = append(res, locations...)
		}
	}

	if len(errs) > 0 {
		err = utilerrors.NewAggregate(errs)
	}

	return
}

func searchConstValueIdentifierReference(ctx context.Context, ss *cache.Snapshot, file uri.URI, valueName string) (res []protocol.Location, err error) {
	ast, err := ss.Parse(ctx, file)
	if err != nil {
		return
	}

	if ast.AST() == nil {
		return
	}

	jumpField := func(field *parser.Field) {
		if field.BadNode || field.ConstValue == nil || field.ConstValue.TypeName != "identifier" {
			return
		}
		if field.ConstValue.Value == valueName {
			res = append(res, jump(file, field.ConstValue))
		}
	}

	for _, st := range ast.AST().Structs {
		if st.BadNode {
			continue
		}
		for _, field := range st.Fields {
			jumpField(field)
		}
	}

	for _, union := range ast.AST().Unions {
		if union.BadNode {
			continue
		}
		for _, field := range union.Fields {
			jumpField(field)
		}
	}

	for _, excep := range ast.AST().Exceptions {
		if excep.BadNode {
			continue
		}
		for _, field := range excep.Fields {
			jumpField(field)
		}
	}

	for _, cst := range ast.AST().Consts {
		if cst.BadNode || cst.Value == nil || cst.Value.TypeName != "identifier" {
			continue
		}
		if cst.Value.Value == valueName {
			res = append(res, jump(file, cst.Value))
		}
	}

	for _, enum := range ast.AST().Enums {
		if enum.BadNode {
			continue
		}
		for _, enumValue := range enum.Values {
			if enumValue.ValueNode == nil || enumValue.ValueNode.TypeName != "identifier" {
				continue
			}

			if enumValue.ValueNode.Value == valueName {
				res = append(res, jump(file, enumValue.ValueNode))
			}
		}
	}

	for _, svc := range ast.AST().Services {
		for _, fn := range svc.Functions {
			for _, field := range fn.Arguments {
				jumpField(field)
			}
			if fn.Throws != nil {
				for _, field := range fn.Throws.Fields {
					jumpField(field)
				}
			}
		}
	}

	return
}
