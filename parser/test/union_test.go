package test

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func Test_ParseUnionIdentifierError(t *testing.T) {

	demoContent := `union {  // err1, line 1, col 7
  // Name is demo name
  1: required string Name;
  2: optional boo Required = true;
}

union 123Demos {  // err3, line 7, col 7
}

union Demos{
}
`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)

		errPos := []string{"1:7", "7:7"}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, containsError(errList.Errors(), parser.InvalidUnionIdentifierError))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}

func Test_ParseUnionBlockRCURError(t *testing.T) {

	demoContent := `union Demo {
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
		assert.True(t, containsError(errList.Errors(), parser.InvalidUnionBlockRCURError))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}

func Test_ParseUnionFieldError(t *testing.T) {
	demoContent := `union Demo {
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
		assert.True(t, containsError(errList.Errors(), parser.InvalidUnionFieldError))
		assert.True(t, containsError(errList.Errors(), parser.InvalidFieldIndexError))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}

func Test_ParseUnionFieldDefault(t *testing.T) {
	demoContent := `union Demo {
1: optional set<string> with_default = [ "test", "aaa" ]
2: optional set<binary> bin_set = {}
3: optional map<binary,i32> bin_map = {}
}
`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.NoError(t, err)

	assert.NotNil(t, ast)
}

func Test_ParseUnionLocation(t *testing.T) {
	demoContent := `union Demo {
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
	assert.Equal(t, 71, docEnd.Offset)

	assert.Len(t, doc.Unions, 1)
	unionNamePos := doc.Unions[0].Name.Pos()
	unionNameEnd := doc.Unions[0].Name.End()

	// union pos
	assert.Equal(t, 1, unionNamePos.Line)
	assert.Equal(t, 7, unionNamePos.Col)
	assert.Equal(t, 6, unionNamePos.Offset)
	// union end
	assert.Equal(t, 1, unionNameEnd.Line)
	assert.Equal(t, 11, unionNameEnd.Col)
	assert.Equal(t, 10, unionNameEnd.Offset)
}
