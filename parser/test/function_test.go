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
