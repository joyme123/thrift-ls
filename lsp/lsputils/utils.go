package lsputils

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/parser"
	"go.lsp.dev/uri"
)

var (
	includeRootsMu sync.RWMutex
	includeRoots   []string
)

func SetIncludeRoots(roots []string) {
	includeRootsMu.Lock()
	defer includeRootsMu.Unlock()
	includeRoots = append([]string(nil), roots...)
}

func getIncludeRoots() []string {
	includeRootsMu.RLock()
	defer includeRootsMu.RUnlock()
	return append([]string(nil), includeRoots...)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

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
		if includeNameMatchesPath(include.Path.Value.Text, includeName) {
			return include.Path.Value.Text
		}
	}

	return ""
}

func includeNameMatchesPath(includePath, includeName string) bool {
	items := strings.Split(filepath.ToSlash(includePath), "/")
	if len(items) == 0 {
		return false
	}

	fileName := items[len(items)-1]
	if !strings.HasSuffix(fileName, ".thrift") {
		return false
	}

	fileName = strings.TrimSuffix(fileName, ".thrift")
	return includeName == fileName || (len(items) > 1 && includeName == items[0])
}

// cur is current file uri. for example file:///tmp/user.thrift
// includePath is include name used in code. for example: base.thrift
func IncludeURI(cur uri.URI, includePath string) uri.URI {
	filePath := cur.Filename()
	basePath := filepath.Dir(filePath)

	relativePath := filepath.Join(basePath, includePath)
	if fileExists(relativePath) {
		return uri.File(relativePath)
	}

	for _, root := range getIncludeRoots() {
		candidate := filepath.Join(root, includePath)
		if fileExists(candidate) {
			return uri.File(candidate)
		}
	}

	return uri.File(relativePath)
}

func IncludeExists(cur uri.URI, includePath string) bool {
	return fileExists(IncludeURI(cur, includePath).Filename())
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
			includeNames = append(includeNames, GetIncludeName(u))

			items := strings.Split(filepath.ToSlash(path), "/")
			if len(items) > 1 && items[0] != "." && items[0] != ".." && items[0] != includeNames[len(includeNames)-1] {
				includeNames = append(includeNames, items[0])
			}
		}
	}

	return includeNames
}
