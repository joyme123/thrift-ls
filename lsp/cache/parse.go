package cache

import (
	"fmt"
	"sync"

	"github.com/joyme123/thrift-ls/lsp/mapper"
	"github.com/joyme123/thrift-ls/parser"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/uri"
)

type ParseCaches struct {
	mu     sync.RWMutex
	caches map[uri.URI]*ParsedFile
	tokens map[string]struct{}
}

func NewParseCaches() *ParseCaches {
	return &ParseCaches{
		caches: make(map[uri.URI]*ParsedFile),
	}
}

func (c *ParseCaches) Set(filePath uri.URI, res *ParsedFile) {
	c.mu.Lock()
	c.caches[filePath] = res
	c.tokens = nil
	c.mu.Unlock()
}

func (c *ParseCaches) Get(filePath uri.URI) *ParsedFile {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.caches[filePath]
}

func (c *ParseCaches) Forget(filePath uri.URI) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.caches, filePath)
	c.tokens = nil
}

func (c *ParseCaches) Clone() *ParseCaches {
	c.mu.RLock()
	defer c.mu.RUnlock()

	clone := make(map[uri.URI]*ParsedFile)
	for i := range c.caches {
		clone[i] = c.caches[i]
	}
	newCaches := &ParseCaches{
		caches: clone,
	}
	return newCaches
}

func (c *ParseCaches) Tokens() map[string]struct{} {
	if len(c.tokens) > 0 {
		return c.tokens
	}

	tokens := make(map[string]struct{})
	for _, parsed := range c.caches {
		if parsed.ast == nil {
			continue
		}
		for _, item := range parsed.ast.Includes {
			if item.Path == nil || item.Path.BadNode {
				continue
			}
			tokens[item.Name()] = struct{}{}
		}
		for _, item := range parsed.ast.Enums {
			if item.Name == nil || item.Name.BadNode {
				continue
			}
			tokens[item.Name.Name] = struct{}{}
		}
		for _, item := range parsed.ast.Consts {
			if item.Name == nil || item.Name.BadNode {
				continue
			}
			tokens[item.Name.Name] = struct{}{}
		}
		for _, item := range parsed.ast.Typedefs {
			if item.Alias == nil || item.Alias.BadNode {
				continue
			}
			tokens[item.Alias.Name] = struct{}{}
		}
		for _, item := range parsed.ast.Services {
			if item.Name == nil || item.Name.BadNode {
				continue
			}
			tokens[item.Name.Name] = struct{}{}
		}
		for _, item := range parsed.ast.Unions {
			if item.Name == nil || item.Name.BadNode {
				continue
			}
			tokens[item.Name.Name] = struct{}{}
		}
		for _, item := range parsed.ast.Structs {
			if item.Identifier == nil || item.Identifier.BadNode {
				continue
			}
			tokens[item.Identifier.Name] = struct{}{}

			for _, field := range item.Fields {
				if field.BadNode || field.Identifier == nil || field.Identifier.BadNode {
					continue
				}
				tokens[field.Identifier.Name] = struct{}{}
			}
		}
		for _, item := range parsed.ast.Exceptions {
			if item.Name == nil || item.Name.BadNode {
				continue
			}
			tokens[item.Name.Name] = struct{}{}
		}
	}
	c.tokens = tokens

	return tokens
}

type ParsedFile struct {
	fh FileHandle
	// ast is latest available ast. current fh content may not to be parsed.
	// so it may be nil when fh content is invalid
	ast *parser.Document

	mapper *mapper.Mapper

	// errs hold all ast parsing errors
	errs []parser.ParserError
}

func (p *ParsedFile) Mapper() *mapper.Mapper {
	return p.mapper
}

func (p *ParsedFile) AST() *parser.Document {
	return p.ast
}

func (p *ParsedFile) Errors() []parser.ParserError {
	return p.errs
}

func (p *ParsedFile) AggregatedError() error {
	if len(p.errs) == 0 {
		return nil
	}
	return fmt.Errorf("aggregated error: %v", p.errs)
}

// TODO(jpf): use promise
func Parse(fh FileHandle) (*ParsedFile, error) {
	content, err := fh.Content()
	if err != nil {
		return nil, err
	}

	pf := &ParsedFile{
		fh: fh,
	}

	psr := &parser.PEGParser{}

	ast, errs := psr.Parse(fh.URI().Filename(), content)
	for i := range errs {
		parserErr, ok := errs[i].(parser.ParserError)
		if ok {
			pf.errs = append(pf.errs, parserErr)
		}
	}
	pf.ast = ast
	log.Debugf("peg parsed err: %v", errs)

	mp := mapper.NewMapper(fh.URI(), content)
	pf.mapper = mp

	return pf, nil
}

// type ParseError struct {
// 	Pos Position
// 	Msg string
// }
//
// func (e *ParseError) Error() string {
// 	if e.Pos.Filename != "" || e.Pos.IsValid() {
// 		// don't print "<unknown position>"
// 		// TODO(gri) reconsider the semantics of Position.IsValid
// 		return e.Pos.String() + ": " + e.Msg
// 	}
// 	return e.Msg
// }

type Position struct {
	Filename string // filename, if any
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (byte count)
}

func (p Position) IsValid() bool {
	return p.Line > 0
}

func (p Position) String() string {
	s := p.Filename
	if p.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d", p.Line)
		if p.Column != 0 {
			s += fmt.Sprintf(":%d", p.Column)
		}
	}
	if s == "" {
		s = "-"
	}
	return s
}
