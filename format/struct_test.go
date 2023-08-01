package format

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func TestMustFormatStruct(t *testing.T) {
	doc := `
// comments


/*
 * comments 
 */
struct test {
  /*
   * field 1
   */
  1: required string test,
}      (a.b = "c")          // endline comments
`

	ast, err := parser.Parse("test.thrift", []byte(doc))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	type args struct {
		st *parser.Struct
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				st: ast.(*parser.Document).Structs[0],
			},
			want: `// comments

/*
 * comments
 */
struct test {
    /*
     * field 1
     */
    1: required string test,
} (a.b = "c") // endline comments
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MustFormatStruct(tt.args.st))
		})
	}
}
