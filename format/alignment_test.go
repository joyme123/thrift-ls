package format

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/require"
)

func TestFormatAlignsEndLineCommentsWithFields(t *testing.T) {
	previous := Align
	t.Cleanup(func() { Align = previous })
	Align = AlignTypeField

	ast, err := parser.Parse("test.thrift", []byte(`struct A {
1: string a // first
2: i64 longer // second
}`))
	require.NoError(t, err)

	formatted, err := FormatDocument(ast.(*parser.Document))
	require.NoError(t, err)
	require.Equal(t, `struct A {
    1: string a      // first
    2: i64    longer // second
}`, formatted)
}

func TestFormatAlignsEnumAssignments(t *testing.T) {
	previous := Align
	t.Cleanup(func() { Align = previous })
	Align = AlignTypeAssign

	ast, err := parser.Parse("test.thrift", []byte(`enum A {
ONE = 1
TWO_LONG = 2
THREE = 333
}`))
	require.NoError(t, err)

	formatted, err := FormatDocument(ast.(*parser.Document))
	require.NoError(t, err)
	require.Equal(t, `enum A {
    ONE      = 1
    TWO_LONG = 2
    THREE    = 333
}`, formatted)
}
