package test

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func Test_ParseEnum(t *testing.T) {
	demoContent := `enum TweetType {
    TWEET,         // 0
    RETWEET = 2,   // 2
    DM = 0xa,      // 10
    REPLY          // 11
    POST = 0o21    // 17
}`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	enums := ast.(*parser.Document).Enums
	assert.Len(t, enums, 1)
	assert.Equal(t, "TweetType", enums[0].Name.Name.Text)
	assert.Len(t, enums[0].Values, 5)
	assert.Equal(t, "TWEET", enums[0].Values[0].Name.Name.Text)
	assert.Equal(t, int64(0), enums[0].Values[0].Value)
	assert.Equal(t, "RETWEET", enums[0].Values[1].Name.Name.Text)
	assert.Equal(t, int64(2), enums[0].Values[1].Value)
	assert.Equal(t, "DM", enums[0].Values[2].Name.Name.Text)
	assert.Equal(t, int64(10), enums[0].Values[2].Value)
	assert.Equal(t, "REPLY", enums[0].Values[3].Name.Name.Text)
	assert.Equal(t, int64(11), enums[0].Values[3].Value)
	assert.Equal(t, "POST", enums[0].Values[4].Name.Name.Text)
	assert.Equal(t, int64(17), enums[0].Values[4].Value)
}

func Test_ParseEnumIdentifierError(t *testing.T) {

	demoContent := `enum {  // err1, line 1, col 6
  TWEET,         // 0
  RETWEET = 2,   // 2
  DM = 0xa,      // 10
  REPLY          // 11
  POST = 0o21    // 17
}

enum 123Demos {  // err2, line 9, col 6
}

enum Demos{
}
`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)

		errPos := []string{"1:6", "9:6"}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, containsError(errList.Errors(), parser.InvalidEnumIdentifierError))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}

func Test_ParseEnumBlockRCURError(t *testing.T) {

	demoContent := `enum Demo {  
  TWEET,         // 0
  RETWEET = 2,   // 2
  DM = 0xa,      // 10
  REPLY          // 11
  POST = 0o21    // 17 // err1, line 8, col 1

enum Demos1 {}  
enum Demos{ // err2, line 10, col 1
`

	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)

		errPos := []string{"8:1", "10:1"}
		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}

		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, containsError(errList.Errors(), parser.InvalidEnumBlockRCURError))

	}

	assert.NotNil(t, ast)
}

func Test_ParseEnumValueError(t *testing.T) {
	demoContent := `enum Demo {  
  TWEET,         // 0
  RETWEET = "aaa",   // 2, err1, line 3, col 13
  DM = 0xa,      // 10
  REPLY          // 11
}

enum Demos1 {}
`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.Error(t, err)
	if err != nil {
		errList, ok := err.(parser.ErrorLister)
		assert.True(t, ok)
		errPos := []string{"3:13"}
		errs := []error{parser.InvalidEnumValueIntConstantError}
		assert.Len(t, errList.Errors(), len(errPos))
		assert.True(t, equalErrors(errList.Errors(), errs))

		for i, err := range errList.Errors() {
			assert.Contains(t, err.Error(), errPos[i])
			t.Logf("error %d: %v\n", i, err)
		}
	}

	assert.NotNil(t, ast)
}

func Test_ParseEnumLocation(t *testing.T) {
	demoContent := `enum Demo {
  TWEET,         // 0
  RETWEET = 2,   // 2, err1, line 3, col 13
  DM = 0xa,      // 10
  REPLY          // 11
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
	assert.Equal(t, 6, docEnd.Line)
	assert.Equal(t, 3, docEnd.Col)
	assert.Equal(t, len(demoContent), docEnd.Offset)

	assert.Len(t, doc.Enums, 1)
	namePos := doc.Enums[0].Name.Name.Pos()
	nameEnd := doc.Enums[0].Name.Name.End()

	// name pos
	assert.Equal(t, 1, namePos.Line)
	assert.Equal(t, 6, namePos.Col)
	assert.Equal(t, 5, namePos.Offset)
	// name end
	assert.Equal(t, 1, nameEnd.Line)
	assert.Equal(t, 10, nameEnd.Col)
	assert.Equal(t, 9, nameEnd.Offset)

	// enum value
	assert.Len(t, doc.Enums[0].Values, 4)
	enumValue := doc.Enums[0].Values[0]
	enumValuePos := enumValue.Pos()
	enumValueEnd := enumValue.End()

	// enum value pos
	assert.Equal(t, 2, enumValuePos.Line)
	assert.Equal(t, 3, enumValuePos.Col)
	assert.Equal(t, 14, enumValuePos.Offset)

	// enum value end
	assert.Equal(t, 2, enumValueEnd.Line)
	assert.Equal(t, 18, enumValueEnd.Col)
	assert.Equal(t, 29, enumValueEnd.Offset)
}
