package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchNodePath(t *testing.T) {
	demoContent := `struct Demo {
1: optional set<string> with_default = [ "😀", "aaa" ]
}`

	ast, err := Parse("test.thrift", []byte(demoContent))
	assert.NoError(t, err)
	assert.NotNil(t, ast)
	doc := ast.(*Document)

	type args struct {
		root Node
		pos  Position
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "struct name",
			args: args{
				root: doc,
				pos: Position{
					Line: 1,
					Col:  11,
				},
			},
			want: []string{"Document", "Struct", "Identifier"},
		},
		{
			name: "struct field value",
			args: args{
				root: doc,
				pos: Position{
					Line: 2,
					Col:  44,
				},
			},
			want: []string{"Document", "Struct", "Field", "ConstValue"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := SearchNodePath(tt.args.root, tt.args.pos)
			pathStr := make([]string, 0, len(path))
			for i := range path {
				nodeType := path[i].Type()
				pathStr = append(pathStr, nodeType)
			}
			assert.Equal(t, tt.want, pathStr)
		})
	}
}
