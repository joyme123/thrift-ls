package test

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func Test_ParseConstValue(t *testing.T) {

	demoContent := `const i16 a = 1
const i32 b = 0o1
const i64 c = 0xa1 
const double d = 1.33333333 
const double e = 1.3333e11
const double f = 1e11
const double g = 1.3333E11
const double h = 1E11`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	consts := ast.(*parser.Document).Consts
	assert.Len(t, consts, 8)
	assert.Equal(t, "1", consts[0].Value.ValueInText)
	assert.Equal(t, "0o1", consts[1].Value.ValueInText)
	assert.Equal(t, "0xa1", consts[2].Value.ValueInText)
	assert.Equal(t, "1.33333333", consts[3].Value.ValueInText)
	assert.Equal(t, "1.3333e11", consts[4].Value.ValueInText)
	assert.Equal(t, "1e11", consts[5].Value.ValueInText)
	assert.Equal(t, "1.3333E11", consts[6].Value.ValueInText)
	assert.Equal(t, "1E11", consts[7].Value.ValueInText)
}
