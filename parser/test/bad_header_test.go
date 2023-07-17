package test

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func Test_ParseBadHeader(t *testing.T) {
	demoContent := `include "aaa"
aaaaa`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	assert.NotNil(t, ast)

	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)

		errPos := []string{"2:1"}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, containsError(errList.Errors(), parser.InvalidHeaderError))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	demoContent2 := `inc`
	ast, err = parser.Parse("test.thrift", []byte(demoContent2))
	assert.Error(t, err)
	assert.NotNil(t, ast)

	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)

		errPos := []string{"1:1"}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, containsError(errList.Errors(), parser.InvalidHeaderError))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}
}
