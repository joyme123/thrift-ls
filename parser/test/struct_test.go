package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func Test_ParseStructIdentifierError(t *testing.T) {

	demoContent := `struct {  // err1, line 1, col 8
  // Name is demo name
  1: required string Name;
  2: optional boo Required = true;
}

struct 123Demos {  // err3, line 7, col 8
}

struct Demos{
}
`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)

		errPos := []string{"1:8", "7:8"}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, containsError(errList.Errors(), parser.InvalidStructIdentifierError))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)

	data, err := json.MarshalIndent(ast, "", "  ")
	assert.NoError(t, err)

	os.WriteFile("/tmp/ast", data, os.ModePerm)
}

func Test_ParseStructBlockLCURError(t *testing.T) {

	demoContent := `struct Demo
  // Name is demo name, err1, line 3, col 12
  2: optional boo Required = true; 
}
struct Demos!!{ } // err2, line5, col 13
struct Demos} // err2, line 6, col 13
struct Demos{}
`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)
		errPos := []string{"3:3", "5:13", "6:13"}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, containsError(errList.Errors(), parser.InvalidStructBlockLCURError))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)

	data, err := json.MarshalIndent(ast, "", "  ")
	assert.NoError(t, err)

	os.WriteFile("/tmp/ast", data, os.ModePerm)
}

func Test_ParseStructBlockRCURError(t *testing.T) {

	demoContent := `struct Demo {
  // Name is demo name

  2: optional boo Required = true; // err1, line 6, col 1

struct Demos{}
struct Demos{ // err2, line 8, col 1
struct Demos{}
`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)
		errPos := []string{"6:1", "8:1"}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, containsError(errList.Errors(), parser.InvalidStructBlockRCURError))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)

	data, err := json.MarshalIndent(ast, "", "  ")
	assert.NoError(t, err)

	os.WriteFile("/tmp/ast", data, os.ModePerm)
}

func Test_ParseStructFieldError(t *testing.T) {
	demoContent := `struct Demo {
  1: optional i64 count
  a: optional boo Required = true; // err1, line 3, col 0
  2: required i32 test4;
  required string test; // err2, line 5, col 0
  4: required i32 test;
  5 required string test; // err3, line 7, col 0
  6: required test test;
  no comment // err4, line 9, col 0
}
`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)
		errPos := []string{"3:0", "5:0", "7:0", "9:0"}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, containsError(errList.Errors(), parser.InvalidStructFieldError))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)

	data, err := json.MarshalIndent(ast, "", "  ")
	assert.NoError(t, err)

	os.WriteFile("/tmp/ast", data, os.ModePerm)
}
