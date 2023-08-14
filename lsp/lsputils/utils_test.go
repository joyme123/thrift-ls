package lsputils

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/uri"
)

func Test_IncludeURI(t *testing.T) {
	type args struct {
		cur         uri.URI
		includePath string
	}
	tests := []struct {
		name string
		args args
		want uri.URI
	}{
		{
			name: "case1",
			args: args{
				cur:         uri.File("/tmp/workspace/app.thrift"),
				includePath: "../user.thrift",
			},
			want: uri.File("/tmp/user.thrift"),
		},
		{
			name: "case2",
			args: args{
				cur:         uri.File("/tmp/workspace/app.thrift"),
				includePath: "user.thrift",
			},
			want: uri.File("/tmp/workspace/user.thrift"),
		},
		{
			name: "case3",
			args: args{
				cur:         uri.URI("file:///c:/Users/Administrator/Downloads/galaxy-thrift-api-master/galaxy-thrift-api-master/sds/Common.thrift"),
				includePath: "Errors.thrift",
			},
			want: uri.URI("file:///c:/Users/Administrator/Downloads/galaxy-thrift-api-master/galaxy-thrift-api-master/sds/Errors.thrift"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IncludeURI(tt.args.cur, tt.args.includePath))
		})
	}
}

func TestGetIncludePath(t *testing.T) {
	file := `include "../../user.thrift"
service Demo {
  user.Test Api(1:user.Test2 arg1, 2:user.Test3 arg2) throws (1:user.Error1 err)
}`
	ast, err := parser.Parse("file:///test.thrift", []byte(file))
	assert.NoError(t, err)

	type args struct {
		ast         *parser.Document
		includeName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case",
			args: args{
				ast:         ast.(*parser.Document),
				includeName: "user",
			},
			want: "../../user.thrift",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetIncludePath(tt.args.ast, tt.args.includeName))
		})
	}
}

func TestGetIncludeName(t *testing.T) {
	type args struct {
		file uri.URI
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "file name",
			args: args{
				file: uri.New("base.thrift"),
			},
			want: "base",
		},
		{
			name: "file name with dir",
			args: args{
				file: uri.New("/tmp/base.thrift"),
			},
			want: "base",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetIncludeName(tt.args.file))
		})
	}
}
