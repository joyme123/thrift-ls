package parser

import (
	"fmt"
	"testing"

	"github.com/joyme123/thrift-ls/utils"
	"github.com/stretchr/testify/assert"
)

func TestSearchNodePath(t *testing.T) {
	demoContent := `struct Demo {
1: optional set<string> with_default = [ "😀", "aaa" ]
}

service Demo {
  user.Test Api(1:user.Test2 arg1, 2:user.Test3 arg2) throws (1:user.Error1 err)
}

enum Test {
  ONE
  TWO = 2
}
`

	ast, err := Parse("test.thrift", []byte(demoContent))
	assert.NoError(t, err)
	assert.NotNil(t, ast)
	doc := ast.(*Document)

	utils.MustDumpJsonToFile(doc, "/tmp/ast")

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
			want: []string{"Document", "Struct", "Identifier", "IdentifierName"},
		},
		{
			name: "field type",
			args: args{
				root: doc,
				pos: Position{
					Line: 2,
					Col:  13,
				},
			},
			want: []string{"Document", "Struct", "Field", "FieldType", "TypeName"},
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
		{
			name: "function type",
			args: args{
				root: doc,
				pos: Position{
					Line: 6,
					Col:  3,
				},
			},
			want: []string{"Document", "Service", "Function", "FieldType", "TypeName"},
		},
		{
			name: "function type",
			args: args{
				root: doc,
				pos: Position{
					Line: 6,
					Col:  8,
				},
			},
			want: []string{"Document", "Service", "Function", "FieldType", "TypeName"},
		},
		{
			name: "function argument",
			args: args{
				root: doc,
				pos: Position{
					Line: 6,
					Col:  24,
				},
			},
			want: []string{"Document", "Service", "Function", "Field", "FieldType", "TypeName"},
		},
		{
			name: "function throws",
			args: args{
				root: doc,
				pos: Position{
					Line: 6,
					Col:  70,
				},
			},
			want: []string{"Document", "Service", "Function", "Throws", "Field", "FieldType", "TypeName"},
		},
		{
			name: "enum value",
			args: args{
				root: doc,
				pos: Position{
					Line: 11,
					Col:  3,
				},
			},
			want: []string{"Document", "Enum", "EnumValue", "Identifier", "IdentifierName"},
		},
		{
			name: "enum value const value",
			args: args{
				root: doc,
				pos: Position{
					Line: 11,
					Col:  9,
				},
			},
			want: []string{"Document", "Enum", "EnumValue", "ConstValue"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "enum value" {
				fmt.Println()
			}
			path := SearchNodePathByPosition(tt.args.root, tt.args.pos)
			pathStr := make([]string, 0, len(path))
			for i := range path {
				nodeType := path[i].Type()
				pathStr = append(pathStr, nodeType)
			}
			assert.Equal(t, tt.want, pathStr)
		})
	}
}
