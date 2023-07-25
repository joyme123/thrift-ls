package test

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func Test_ParseServiceIdentifierError(t *testing.T) {

	demoContent := `service "aa" 
	{}  // err1, line 1, col 8

`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)

		errPos := []string{"1:9"}
		errs := []error{parser.InvalidServiceIdentifierError}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, equalErrors(errList.Errors(), errs))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}

func Test_ParseServiceBlockRCURError(t *testing.T) {

	demoContent := `service aa {`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)
		errPos := []string{"1:13"}
		errs := []error{parser.InvalidServiceBlockRCURError}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, equalErrors(errList.Errors(), errs))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}

func Test_ParseServiceFunctionError(t *testing.T) {

	demoContent := `service aa {
  string Tes
}`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)
		errPos := []string{"2:3"}
		errs := []error{parser.InvalidServiceFunctionError}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, equalErrors(errList.Errors(), errs))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}

func Test_ParseServiceOK(t *testing.T) {

	demoContent := `service aa extends bb 
{
  string Test(1:required string test) throws (1:required error err)

}`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.NoError(t, err)
	assert.NotNil(t, ast)
}
