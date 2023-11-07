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
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func PrepareRename(ctx context.Context, ss *cache.Snapshot, file uri.URI, pos protocol.Position) (res *protocol.Range, err error) {
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
	nodePath := parser.SearchNodePathByPosition(pf.AST(), astPos)
	targetNode := nodePath[len(nodePath)-1]

	switch targetNode.Type() {
	case "IdentifierName", "ConstValue":
		rg := lsputils.ASTNodeToRange(targetNode)
		return &rg, nil
	default:
		err = fmt.Errorf("%s doesn't support rename", targetNode.Type())
		return
	}
}

func Rename(ctx context.Context, ss *cache.Snapshot, file uri.URI, pos protocol.Position, newName string) (res *protocol.WorkspaceEdit, err error) {
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
	nodePath := parser.SearchNodePathByPosition(pf.AST(), astPos)
	targetNode := nodePath[len(nodePath)-1]

	self := lsputils.ASTNodeToRange(targetNode)

	switch targetNode.Type() {
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
			locations, err := searchConstValueIdentifierReferences(ctx, ss, file, typeName)
			if err != nil {
				return nil, err
			}

			locations = append(locations, protocol.Location{
				URI:   file,
				Range: self,
			})

			return convertLocationToWorkspaceEdit(locations, file, newName), nil
		} else if definitionType == "Service" {
			svcName := targetNode.(*parser.IdentifierName).Text
			if !strings.Contains(svcName, ".") {
				svcName = fmt.Sprintf("%s.%s", lsputils.GetIncludeName(file), svcName)
			} else {
				include, _, _ := strings.Cut(svcName, ".")
				path := lsputils.GetIncludePath(pf.AST(), include)
				if path != "" { // doesn't match any include path
					file = lsputils.IncludeURI(file, path)
				}
			}
			locations, err := searchServiceReferences(ctx, ss, file, svcName)
			if err != nil {
				return nil, err
			}

			locations = append(locations, protocol.Location{
				URI:   file,
				Range: self,
			})

			return convertLocationToWorkspaceEdit(locations, file, newName), nil
		}

		if _, ok := validReferenceDefinitionType[definitionType]; !ok {
			return
		}

		// typeName is base.User
		typeName := fmt.Sprintf("%s.%s", lsputils.GetIncludeName(file), targetNode.(*parser.IdentifierName).Text)
		locations, err := searchIdentifierReferences(ctx, ss, file, typeName, definitionType)
		if err != nil {
			return nil, err
		}

		locations = append(locations, protocol.Location{
			URI:   file,
			Range: self,
		})
		return convertLocationToWorkspaceEdit(locations, file, newName), nil
	case "ConstValue":
		locations, err := searchConstValueReferences(ctx, ss, file, pf.AST(), nodePath, targetNode)
		if err != nil {
			return nil, err
		}

		locations = append(locations, protocol.Location{
			URI:   file,
			Range: self,
		})
		return convertLocationToWorkspaceEdit(locations, file, newName), nil
	default:
		err = fmt.Errorf("%s doesn't support rename", targetNode.Type())
		return
	}
}

func convertLocationToWorkspaceEdit(locations []protocol.Location, fileURI uri.URI, newName string) *protocol.WorkspaceEdit {
	res := &protocol.WorkspaceEdit{
		Changes: make(map[protocol.DocumentURI][]protocol.TextEdit),
	}

	for _, loc := range locations {
		newText := newName
		if loc.URI != fileURI {
			newText = lsputils.GetIncludeName(fileURI) + "." + newName
		}
		res.Changes[loc.URI] = append(res.Changes[loc.URI], protocol.TextEdit{
			Range:   loc.Range,
			NewText: newText,
		})
	}

	return res
}
