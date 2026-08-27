package test

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func Test_ParseFunctionIdentifierError(t *testing.T) {
	demoContent := `service Demo {
  string 11GetName(1:required string name) throws(1:required error err1)
}
`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)
		errPos := []string{"2:10"}
		errs := []error{parser.InvalidFunctionIdentifierError}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, equalErrors(errList.Errors(), errs))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}

func Test_ParseFunctionArgumentError(t *testing.T) {
	demoContent := `service Demo 
{
  string GetName(1:required string 11name) throws(1:required error err1)
  string GetName(1:required string name) throws(1:required error err1)
}
`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)
		errPos := []string{"3:18"}
		errs := []error{parser.InvalidFunctionArgumentError}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, equalErrors(errList.Errors(), errs))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}

func Test_ParseFunctionArgumentReference(t *testing.T) {
	demoContent := `service Tree {
void insert(1: Node & node) throws (1: TreeError & err)
}
`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	doc := ast.(*parser.Document)
	assert.Len(t, doc.Services, 1)
	assert.False(t, doc.Services[0].ChildrenBadNode())

	fn := doc.Services[0].Functions[0]

	assert.Len(t, fn.Arguments, 1)
	assert.NotNil(t, fn.Arguments[0].ReferenceKeyword)
	assert.Equal(t, "&", fn.Arguments[0].ReferenceKeyword.Literal.Text)
	assert.Equal(t, "node", fn.Arguments[0].Identifier.Name.Text)

	assert.Len(t, fn.Throws.Fields, 1)
	assert.NotNil(t, fn.Throws.Fields[0].ReferenceKeyword)
	assert.Equal(t, "&", fn.Throws.Fields[0].ReferenceKeyword.Literal.Text)
	assert.Equal(t, "err", fn.Throws.Fields[0].Identifier.Name.Text)
}
