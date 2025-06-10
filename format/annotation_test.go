package format

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func TestMustFormatAnnotations(t *testing.T) {

	doc := `struct a {}/* xxx */ ( /* xxx */ key /* xxx */ = /* xxx */ "value", key2 = "value2", key3 = "value3" ) /* xxx */`
	ast, err := parser.Parse("test.thrift", []byte(doc))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	type args struct {
		annotations *parser.Annotations
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				annotations: ast.(*parser.Document).Structs[0].Annotations,
			},
			want: `/* xxx */ (/* xxx */ key /* xxx */ = /* xxx */ "value", key2 = "value2", key3 = "value3")`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MustFormatAnnotations(tt.args.annotations))
		})
	}
}

func TestMustFormatAnnotationsInvalidCase(t *testing.T) {

	doc := `
struct Foo {
  // 用户名列表
  2: list<string> strings (
    custom_tag_1 = "1"
    custom_tag_2 = "2"
  )
}
`
	ast, err := parser.Parse("test.thrift", []byte(doc))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	type args struct {
		annotations *parser.Annotations
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				annotations: ast.(*parser.Document).Structs[0].Annotations,
			},
			want: `
struct Foo {
  // 用户名列表
  2: list<string> strings (
    custom_tag_1 = "1"
    custom_tag_2 = "2"
  )
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MustFormatAnnotations(tt.args.annotations))
		})
	}
}
