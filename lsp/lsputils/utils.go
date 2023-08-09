package lsputils

import (
	"path/filepath"
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
	index := strings.LastIndexByte(fileName, '/')
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
		if include.BadNode || include.Path == nil || include.Path.BadNode {
			continue
		}
		items := strings.Split(include.Path.Value.Text, "/")
		path := items[len(items)-1]
		name, _, found := strings.Cut(path, ".")
		if !found {
			continue
		}
		if name == includeName {
			return include.Path.Value.Text
		}
	}

	return ""
}

// cur is current file uri. for example file:///tmp/user.thrift
// includePath is include name used in code. for example: base
func IncludeURI(cur uri.URI, includePath string) uri.URI {
	filePath := cur.Filename()
	items := strings.Split(filePath, "/")
	basePath := strings.TrimSuffix(filePath, items[len(items)-1])

	path := filepath.Join(basePath, includePath)

	return uri.File(path)
}
