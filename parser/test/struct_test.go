package test

import (
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
}

func Test_ParseStructFieldError(t *testing.T) {
	demoContent := `struct Demo {
  1: optional i64 count
  a: optional boo Required = true; // err1, line 3, col 3
  2: required i32 test4;
  required string test; // err2, line 5, col 3
  4: required i32 test;
  5 required string test; // err3, line 7, col 3
  6: required test test;
  no comment // err4, line 9, col 3
}
`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)
		errPos := []string{"3:3", "5:3", "7:3", "9:3"}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, containsError(errList.Errors(), parser.InvalidStructFieldError))
		assert.True(t, containsError(errList.Errors(), parser.InvalidFieldIndexError))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}

func Test_ParseStructFieldDefault(t *testing.T) {
	demoContent := `struct Demo {
1: optional set<string> with_default = [ "test", "aaa" ]
2: optional set<binary> bin_set = {}
3: optional map<binary,i32> bin_map = {}
}
`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.NoError(t, err)

	assert.NotNil(t, ast)
}

func Test_ParseLocation(t *testing.T) {
	demoContent := `struct Demo {
1: optional set<string> with_default = [ "😀", "aaa" ]
}`
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
	assert.Equal(t, 3, docEnd.Line)
	assert.Equal(t, 3, docEnd.Col)
	assert.Equal(t, 72, docEnd.Offset)

	assert.Len(t, doc.Structs, 1)
	structNamePos := doc.Structs[0].Identifier.Name.Pos()
	structNameEnd := doc.Structs[0].Identifier.Name.End()

	// struct pos
	assert.Equal(t, 1, structNamePos.Line)
	assert.Equal(t, 8, structNamePos.Col)
	assert.Equal(t, 7, structNamePos.Offset)
	// struct end
	assert.Equal(t, 1, structNameEnd.Line)
	assert.Equal(t, 12, structNameEnd.Col)
	assert.Equal(t, 11, structNameEnd.Offset)
}

func Test_ParseStructErr(t *testing.T) {
	demoContent := `struct Demo {`
	_, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
}
