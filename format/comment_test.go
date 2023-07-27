package format

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func TestMustFormatComments(t *testing.T) {
	doc := `/* aaaaa
aaaaa
aaaaa
aaaaa
*/
// aaaaa
// aaaaa

include "a.thrift" // aaaaa`
	ast, err := parser.Parse("test.thrift", []byte(doc))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	type args struct {
		comments []*parser.Comment
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
			},
			want: `/* aaaaa
aaaaa
aaaaa
aaaaa
*/
// aaaaa
// aaaaa`,
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
			assert.Equal(t, tt.want, MustFormatComments(tt.args.comments))
		})
	}
}
