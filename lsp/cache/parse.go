package cache

import (
	"github.com/cloudwego/thriftgo/parser"
)

type ParsedFile struct {
	fh  FileHandle
	ast *parser.Thrift
}

func Parse(fh FileHandle) (*ParsedFile, error) {
	content, err := fh.Content()
	if err != nil {
		return nil, err
	}
	ast, err := parser.ParseString(fh.URI().Filename(), string(content))
	return &ParsedFile{
		fh:  fh,
		ast: ast,
	}, nil
}
