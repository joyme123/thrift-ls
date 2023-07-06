package cache

import (
	"github.com/cloudwego/thriftgo/parser"
)

type ParsedFile struct {
	fh FileHandle
	// ast is latest available ast. current fh content may not to be parsed.
	// so it may be nil when fh content is invalid
	ast *parser.Thrift

	// errs hold all ast parsing errors
	errs []error
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

	ast, err := parser.ParseString(fh.URI().Filename(), string(content))
	if err != nil {
		pf.errs = append(pf.errs, err)
	}

	pf.ast = ast
	return pf, nil
}
