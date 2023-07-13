package completion

import (
	"context"
	"strings"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/utils"
	log "github.com/sirupsen/logrus"
)

var DefaultTokenCompletion Interface = &TokenCompletion{}

type TokenCompletion struct {
}

func (c *TokenCompletion) Completion(ctx context.Context, ss *cache.Snapshot, cmp *CompletionRequest) ([]*CompletionItem, error) {
	tokens := ss.Tokens()

	log.Debugln("all tokens:", tokens)

	parsedFile := ss.GetParsedFile(cmp.Fh.URI())
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

	candidates := make([]string, 0)
	for i := range tokens {
		if len(i) > len(prefix) && strings.HasPrefix(i, string(prefix)) {
			candidates = append(candidates, i)
		}
		if len(candidates) >= 10 {
			break
		}
	}

	log.Debugln("prefix:", string(prefix), "candidates: ", candidates)

	res := make([]*CompletionItem, 0, len(candidates))
	for i := range candidates {
		res = append(res, BuildCompletionItem(candidates[i]))
	}

	return res, nil
}
