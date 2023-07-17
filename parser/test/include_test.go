package test

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func Test_ParseInclude(t *testing.T) {
	demoContent := `include "../user.thrift"
include "../base.thrift"
`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	includes := ast.(*parser.Document).Includes
	assert.Len(t, includes, 2)
	assert.Equal(t, "../user.thrift", includes[0].Path.Value)
	assert.Equal(t, "../base.thrift", includes[1].Path.Value)
}

func Test_ParseIncludeWithError(t *testing.T) {
	demoContent := `include "../user.thrift"
include "../base.thrift
`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	assert.NotNil(t, ast)

	includes := ast.(*parser.Document).Includes
	assert.Len(t, includes, 2)
	assert.Equal(t, "../user.thrift", includes[0].Path.Value)
	assert.Equal(t, true, includes[1].Path.BadNode)
}
