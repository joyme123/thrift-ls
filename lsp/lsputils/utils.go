package lsputils

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/joyme123/thrift-ls/parser"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func ASTNodeToRange(node parser.Node) protocol.Range {
	return protocol.Range{
		Start: protocol.Position{
			Line:      uint32(node.Pos().Line - 1),
			Character: uint32(node.Pos().Col - 1),
		},
		End: protocol.Position{
			Line:      uint32(node.End().Line - 1),
			Character: uint32(node.End().Col - 1),
		},
	}
}

// GetIncludeName return include name by file uri
// for example: file uri is file:///base.thrift, then `base` is include name
func GetIncludeName(file uri.URI) string {
	fileName := file.Filename()
	index := strings.LastIndexByte(fileName, filepath.Separator)
	if index == -1 {
		return fileName
	}
	fileName = string(fileName[index+1:])

	index = strings.LastIndexByte(fileName, '.')
	if index == -1 {
		return fileName
	}
	return string(fileName[0:index])
}

// includeName: base.User. `base` is the includeName. returns ../../base.thrift
// if doesn't match, return empty string
func GetIncludePath(ast *parser.Document, includeName string) string {
	for _, include := range ast.Includes {
		if include.BadNode || include.Path == nil || include.Path.BadNode || include.Path.Value == nil {
			continue
		}
		items := strings.Split(include.Path.Value.Text, "/")
		path := items[len(items)-1]
		if !strings.HasSuffix(path, ".thrift") {
			continue
		}
		name := strings.TrimSuffix(path, ".thrift")
		if name == includeName {
			return include.Path.Value.Text
		}
	}

	return ""
}

// cur is current file uri. for example file:///tmp/user.thrift
// includePath is include name used in code. for example: base.thrift
func IncludeURI(cur uri.URI, includePath string) uri.URI {
	filePath := cur.Filename()
	items := strings.Split(filePath, string(filepath.Separator))
	basePath := strings.TrimSuffix(filePath, items[len(items)-1])

	path := filepath.Join(basePath, includePath)

	return uri.File(path)
}

// ParseIdent parse an identifier. identifier format:
//  1. identifier
//  2. include.identifier
//
// it returns include, ident
func ParseIdent(cur uri.URI, includes []*parser.Include, identifier string) (include, ident string) {
	includeNames := IncludeNames(cur, includes)
	// parse include from includeNames

	sort.SliceStable(includeNames, func(i, j int) bool {
		// sort by string length, make sure longest include match early
		// examples:
		// user.extra
		// user
		return len(includeNames[i]) > len(includeNames[j])
	})

	for _, incName := range includeNames {
		prefix := incName + "."
		if strings.HasPrefix(identifier, prefix) {
			return incName, strings.TrimPrefix(identifier, prefix)
		}
	}

	return "", identifier
}

// IncludeNames returns include names from include ast nodes
func IncludeNames(cur uri.URI, includes []*parser.Include) (includeNames []string) {
	for _, inc := range includes {
		if inc.Path != nil && inc.Path.Value != nil {
			path := inc.Path.Value.Text
			u := IncludeURI(cur, path)
			includeName := GetIncludeName(u)
			includeNames = append(includeNames, includeName)
		}
	}

	return includeNames
}
