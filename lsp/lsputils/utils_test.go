package lsputils

import (
	"context"
	"testing"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/memoize"
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

func buildSnapshotForTest(files []*cache.FileChange) *cache.Snapshot {
	store := &memoize.Store{}
	c := cache.New(store)
	fs := cache.NewOverlayFS(c)
	fs.Update(context.TODO(), files)

	view := cache.NewView("test", "file:///tmp", fs, store)
	ss := cache.NewSnapshot(view, store)

	return ss
}
