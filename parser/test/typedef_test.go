package test

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func Test_ParseTypedefError(t *testing.T) {

	demoContent := `typedef a "b"
typedef a b
typedef a 1a`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)

		errPos := []string{"1:11", "3:11"}
		errs := []error{parser.InvalidTypedefIdentifierError, parser.InvalidTypedefIdentifierError}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, equalErrors(errList.Errors(), errs))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}

func Test_ParseTypedefLocation(t *testing.T) {

	demoContent := `typedef a b`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	doc := ast.(*parser.Document)
	docPos := doc.Location.Pos()
	docEnd := doc.Location.End()

	// doc pos
	assert.Equal(t, 1, docPos.Line)
	assert.Equal(t, 1, docPos.Col)
	assert.Equal(t, 0, docPos.Offset)
	// doc end
	assert.Equal(t, 1, docEnd.Line)
	assert.Equal(t, 12, docEnd.Col)
	assert.Equal(t, len(demoContent), docEnd.Offset)

	assert.Len(t, doc.Typedefs, 1)
	aliasPos := doc.Typedefs[0].Alias.Name.Pos()
	aliasEnd := doc.Typedefs[0].Alias.Name.End()

	// union pos
	assert.Equal(t, 1, aliasPos.Line)
	assert.Equal(t, 11, aliasPos.Col)
	assert.Equal(t, 10, aliasPos.Offset)
	// union end
	assert.Equal(t, 1, aliasEnd.Line)
	assert.Equal(t, 12, aliasEnd.Col)
	assert.Equal(t, 11, aliasEnd.Offset)
}
