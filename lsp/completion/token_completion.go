package completion

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/parser"
	"github.com/joyme123/thrift-ls/utils"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

var DefaultTokenCompletion Interface = &TokenCompletion{}

// TokenCompletion is token based completion. It generates completion list based on identifier in ast
type TokenCompletion struct {
}

var keywords = map[string]protocol.InsertTextFormat{
	"bool":                  protocol.InsertTextFormatPlainText,
	"byte":                  protocol.InsertTextFormatPlainText,
	"i16":                   protocol.InsertTextFormatPlainText,
	"i32":                   protocol.InsertTextFormatPlainText,
	"i64":                   protocol.InsertTextFormatPlainText,
	"double":                protocol.InsertTextFormatPlainText,
	"binary":                protocol.InsertTextFormatPlainText,
	"string":                protocol.InsertTextFormatPlainText,
	"required":              protocol.InsertTextFormatPlainText,
	"optional":              protocol.InsertTextFormatPlainText,
	"include":               protocol.InsertTextFormatPlainText,
	"cpp_include":           protocol.InsertTextFormatPlainText,
	"list<$1>":              protocol.InsertTextFormatSnippet,
	"set<$1>":               protocol.InsertTextFormatSnippet,
	"map<$1, $2>":           protocol.InsertTextFormatSnippet,
	"struct $1 {\n$2\n}":    protocol.InsertTextFormatSnippet,
	"const $1 $2 = $3":      protocol.InsertTextFormatSnippet,
	"service $1 {\n$2\n}":   protocol.InsertTextFormatSnippet,
	"union $1 {\n$2\n}":     protocol.InsertTextFormatSnippet,
	"exception $1 {\n$2\n}": protocol.InsertTextFormatSnippet,
	"throws ($1)":           protocol.InsertTextFormatSnippet,
	"typedef $1 $2":         protocol.InsertTextFormatSnippet,
}

type Candidate struct {
	text   string
	format protocol.InsertTextFormat
}

func (c *TokenCompletion) Completion(ctx context.Context, ss *cache.Snapshot, cmp *CompletionRequest) ([]*CompletionItem, error) {
	tokens := ss.Tokens()

	log.Debugln("all tokens:", tokens)

	parsedFile, err := ss.Parse(ctx, cmp.Fh.URI())
	if err != nil {
		return nil, err
	}

	if parsedFile.AST() == nil {
		return nil, fmt.Errorf("parser ast failed")
	}

	pos, err := parsedFile.Mapper().LSPPosToParserPosition(cmp.Pos)
	if err != nil {
		return nil, err
	}

	content, err := cmp.Fh.Content()
	if err != nil {
		return nil, err
	}

	var prefix []byte
	// get prefix by pos
	for i := pos.Offset - 1; i >= 0; i-- {
		if utils.Space(content[i]) || content[i] == '.' {
			prefix = content[i+1 : pos.Offset]
			break
		}
	}

	candidates := make([]Candidate, 0)

	nodePath := parser.SearchNodePath(parsedFile.AST(), pos)
	if items, err := c.includeCompletion(ss, cmp.Fh.URI(), nodePath); err == nil {
		for i := range items {
			candidates = append(candidates, Candidate{
				text:   items[i],
				format: protocol.InsertTextFormatPlainText,
			})
		}
		log.Debugln("include completion candidates: ", candidates)
	}

	if len(candidates) == 0 {
		searchCandidate := func(token string, format protocol.InsertTextFormat) {
			if len(token) > len(prefix) && strings.HasPrefix(token, string(prefix)) {
				candidates = append(candidates, Candidate{
					text:   token,
					format: format,
				})
			}
		}
		for i := range keywords {
			searchCandidate(i, keywords[i])
			if len(candidates) >= 10 {
				break
			}
		}
		for i := range tokens {
			searchCandidate(i, protocol.InsertTextFormatPlainText)
			if len(candidates) >= 10 {
				break
			}
		}
		log.Debugln("token prefix:", string(prefix), "candidates: ", candidates)
	}

	res := make([]*CompletionItem, 0, len(candidates))
	for i := range candidates {
		res = append(res, BuildCompletionItem(candidates[i]))
	}

	return res, nil
}

func (c *TokenCompletion) includeCompletion(ss *cache.Snapshot, file uri.URI, nodePath []parser.Node) (res []string, err error) {
	if len(nodePath) < 2 {
		return
	}

	if nodePath[len(nodePath)-1].Type() != "Literal" || nodePath[len(nodePath)-2].Type() != "Include" {
		return
	}

	pathPrefix := nodePath[len(nodePath)-1].(*parser.Literal).Value

	currentDir := filepath.Dir(file.Filename())

	log.Debugf("search prefix %s in path %s", pathPrefix, currentDir)

	res, err = ListDirAndFiles(currentDir, pathPrefix)

	log.Debugln("include completion: ", res, "err", err)
	return
}
