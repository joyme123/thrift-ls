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

func Test_ParseStructFieldReference(t *testing.T) {
	demoContent := `struct Node {
1: optional Node & next
2: required i64 value
3: optional list<Node> & children (cpp.ref = "true")
}
`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	doc := ast.(*parser.Document)
	assert.Len(t, doc.Structs, 1)
	assert.False(t, doc.Structs[0].ChildrenBadNode())

	fields := doc.Structs[0].Fields
	assert.Len(t, fields, 3)

	assert.NotNil(t, fields[0].ReferenceKeyword)
	assert.Equal(t, "&", fields[0].ReferenceKeyword.Literal.Text)
	assert.Equal(t, "next", fields[0].Identifier.Name.Text)
	assert.Equal(t, 2, fields[0].ReferenceKeyword.Pos().Line)
	assert.Equal(t, 18, fields[0].ReferenceKeyword.Pos().Col)

	assert.Nil(t, fields[1].ReferenceKeyword)
	assert.Equal(t, "value", fields[1].Identifier.Name.Text)

	assert.NotNil(t, fields[2].ReferenceKeyword)
	assert.Equal(t, "&", fields[2].ReferenceKeyword.Literal.Text)
	assert.Equal(t, "children", fields[2].Identifier.Name.Text)
}

func Test_ParseFieldReferenceCases(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
		// reference is the field name the reference marker must land on, keyed by
		// nothing else because every valid case below declares a single field
		wantIdentifier string
		// wantComments are the comments attached to the field identifier
		wantComments []string
	}{
		{
			name:           "no space around marker",
			content:        "struct A {\n1: Node&next\n}\n",
			wantIdentifier: "next",
		},
		{
			name:           "comment before marker",
			content:        "struct A {\n1: Node /* a */ & next\n}\n",
			wantIdentifier: "next",
		},
		{
			name:           "comment between marker and name",
			content:        "struct A {\n1: Node & /* retain */ next\n}\n",
			wantIdentifier: "next",
			wantComments:   []string{"/* retain */"},
		},
		{
			name:           "comment on both sides of marker",
			content:        "struct A {\n1: Node /* a */ & /* b */ next\n}\n",
			wantIdentifier: "next",
			wantComments:   []string{"/* b */"},
		},
		{
			name:           "newline between marker and name",
			content:        "struct A {\n1: Node &\nnext\n}\n",
			wantIdentifier: "next",
		},
		{
			name:           "line comment between marker and name",
			content:        "struct A {\n1: Node & // c\nnext\n}\n",
			wantIdentifier: "next",
			wantComments:   []string{"// c"},
		},
		{
			name:           "marker with default value",
			content:        "struct A {\n1: optional Node & next = 1\n}\n",
			wantIdentifier: "next",
		},
		{
			name:           "marker with annotations and separator",
			content:        "struct A {\n1: Node & next (cpp.ref = \"true\"),\n}\n",
			wantIdentifier: "next",
		},
		{
			name:           "marker on container type",
			content:        "struct A {\n1: map<string,Node> & m\n}\n",
			wantIdentifier: "m",
		},
		{
			name:           "marker in union field",
			content:        "union A {\n1: Node & next\n}\n",
			wantIdentifier: "next",
		},
		{
			name:           "marker in exception field",
			content:        "exception A {\n1: Node & next\n}\n",
			wantIdentifier: "next",
		},
		{
			name:    "double marker is rejected",
			content: "struct A {\n1: Node && next\n}\n",
			wantErr: true,
		},
		{
			name:    "marker without identifier is rejected",
			content: "struct A {\n1: Node &\n}\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := parser.Parse("test.thrift", []byte(tt.content))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, ast)

			field := onlyField(t, ast.(*parser.Document))
			assert.NotNil(t, field.ReferenceKeyword)
			assert.Equal(t, "&", field.ReferenceKeyword.Literal.Text)
			assert.Equal(t, tt.wantIdentifier, field.Identifier.Name.Text)

			comments := make([]string, 0, len(field.Identifier.Comments))
			for _, c := range field.Identifier.Comments {
				comments = append(comments, c.Text)
			}
			assert.Equal(t, tt.wantComments, nilIfEmpty(comments))
		})
	}
}

// onlyField returns the single field declared by doc, whichever struct-like
// definition holds it, and asserts the definition has no bad child nodes.
func onlyField(t *testing.T, doc *parser.Document) *parser.Field {
	t.Helper()

	var fields []*parser.Field
	switch {
	case len(doc.Structs) == 1:
		assert.False(t, doc.Structs[0].ChildrenBadNode())
		fields = doc.Structs[0].Fields
	case len(doc.Unions) == 1:
		assert.False(t, doc.Unions[0].ChildrenBadNode())
		fields = doc.Unions[0].Fields
	case len(doc.Exceptions) == 1:
		assert.False(t, doc.Exceptions[0].ChildrenBadNode())
		fields = doc.Exceptions[0].Fields
	default:
		t.Fatal("document declares no struct-like definition")
	}

	assert.Len(t, fields, 1)
	return fields[0]
}

func nilIfEmpty(v []string) []string {
	if len(v) == 0 {
		return nil
	}
	return v
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
