package test

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func Test_ParseConstError(t *testing.T) {

	demoContent := `const string a = 1aa  // err1: InvalidConstValueError
const string b // err2: InvalidConstMissingValueError
const string c = "hello // err3: InvalidLiteral1MissingRightError
const string d = 'hello // err4: InvalidLiteral2MissingRightError
const i32 111e = 1 // err5: InvalidConstIdentifierError
const string f = "hello"`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)

		errPos := []string{"1:18", "2:16", "3:18", "4:18", "5:11"}
		errs := []error{
			parser.InvalidConstConstValueError,
			parser.InvalidConstMissingValueError,
			parser.InvalidLiteral1MissingRightError,
			parser.InvalidLiteral2MissingRightError,
			parser.InvalidConstIdentifierError}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, equalErrors(errList.Errors(), errs))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}
