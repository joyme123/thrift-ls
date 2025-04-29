package completion

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/parser"
	"github.com/joyme123/thrift-ls/utils"
	log "github.com/sirupsen/logrus"
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
	"uuid":                  protocol.InsertTextFormatPlainText,
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
	showText   string
	insertText string
	format     protocol.InsertTextFormat
}

func (c *TokenCompletion) Completion(ctx context.Context, ss *cache.Snapshot, cmp *CompletionRequest) ([]*CompletionItem, protocol.Range, error) {
	tokens := ss.Tokens()

	rng := protocol.Range{
		Start: protocol.Position{
			Line:      cmp.Pos.Line,
			Character: cmp.Pos.Character,
		},
		End: protocol.Position{
			Line:      cmp.Pos.Line,
			Character: cmp.Pos.Character,
		},
	}

	log.Debugln("all tokens:", tokens)

	parsedFile, err := ss.Parse(ctx, cmp.Fh.URI())
	if err != nil {
		return nil, rng, err
	}

	if parsedFile.AST() == nil {
		return nil, rng, fmt.Errorf("parser ast failed")
	}

	pos, err := parsedFile.Mapper().LSPPosToParserPosition(cmp.Pos)
	if err != nil {
		return nil, rng, err
	}

	candidates := make([]Candidate, 0)

	log.Debugln("parser pos: ", pos)
	includeLiteralPos := pos
	includeLiteralPos.Col = includeLiteralPos.Col - 1 // remove quote
	nodePath := parser.SearchNodePathByPosition(parsedFile.AST(), includeLiteralPos)
	if items, includeRng, err := c.includeCompletion(ss, cmp.Fh.URI(), nodePath); err == nil {
		for i := range items {
			candidates = append(candidates, items[i])
		}
		if len(items) > 0 {
			rng = includeRng
			log.Debugln("include completion candidates: ", candidates)
		}
	}

	if len(candidates) == 0 {
		nodePath = parser.SearchNodePathByPosition(parsedFile.AST(), pos)
		content, err := cmp.Fh.Content()
		if err != nil {
			return nil, rng, err
		}
		var prefix []byte
		// get prefix by pos
		for i := pos.Offset - 1; i >= 0; i-- {
			if utils.Space(content[i]) || content[i] == '.' || content[i] == '\'' || content[i] == '"' {
				prefix = content[i+1 : pos.Offset]
				rng.Start.Character = rng.Start.Character - uint32(len(prefix))
				break
			}
		}

		if len(prefix) == 0 {
			// prefix is empty, set prefix to content
			prefix = content
			rng.Start.Character = rng.Start.Character - uint32(len(prefix))
		}

		searchCandidate := func(token string, format protocol.InsertTextFormat) {
			if len(token) > len(prefix) && strings.HasPrefix(token, string(prefix)) {
				candidates = append(candidates, Candidate{
					showText:   token,
					insertText: token,
					format:     format,
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

	return res, rng, nil
}

func (c *TokenCompletion) includeCompletion(ss *cache.Snapshot, file uri.URI, nodePath []parser.Node) (res []Candidate, rng protocol.Range, err error) {
	if len(nodePath) < 3 {
		return
	}

	if nodePath[len(nodePath)-1].Type() != "LiteralValue" || nodePath[len(nodePath)-2].Type() != "Literal" || nodePath[len(nodePath)-3].Type() != "Include" {
		return
	}

	targetNode := nodePath[len(nodePath)-1].(*parser.LiteralValue)
	pathPrefix := targetNode.Text
	rng = lsputils.ASTNodeToRange(targetNode)

	currentDir := filepath.Dir(file.Filename())

	log.Debugf("search prefix %s in path %s", pathPrefix, currentDir)

	res, err = ListDirAndFiles(currentDir, pathPrefix)

	log.Debugln("include completion: ", res, "err", err)
	return
}
