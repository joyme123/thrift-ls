package format

import (
	"strings"
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

/*
 * multi line comments
 */
func TestMustFormatComments(t *testing.T) {
	doc := strings.TrimSpace(`
/*
 * aaaaa
 * aaaaa
 * aaaaa
 * aaaaa
 */


// aaaaa



// aaaaa

include "a.thrift" // aaaaa

// endline comments`)
	ast, err := parser.Parse("test.thrift", []byte(doc))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	type args struct {
		comments []*parser.Comment
		indent   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "comments",
			args: args{
				comments: ast.(*parser.Document).Includes[0].Comments,
				indent:   Indent,
			},
			want: strings.TrimSpace(`
/*
 * aaaaa
 * aaaaa
 * aaaaa
 * aaaaa
 */

// aaaaa

// aaaaa`),
		},
		{
			name: "endline comments",
			args: args{
				comments: ast.(*parser.Document).Includes[0].EndLineComments,
			},
			want: `// aaaaa`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MustFormatComments(tt.args.comments, ""))
		})
	}
}
