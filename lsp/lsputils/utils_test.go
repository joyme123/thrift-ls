package lsputils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/uri"
)

func TestIncludeExists(t *testing.T) {
	root := t.TempDir()
	include := filepath.Join(root, "shared", "v1", "identity.thrift")
	assert.NoError(t, os.MkdirAll(filepath.Dir(include), 0o755))
	assert.NoError(t, os.WriteFile(include, []byte(""), 0o644))
	SetIncludeRoots([]string{root})
	t.Cleanup(func() { SetIncludeRoots(nil) })

	current := uri.File(filepath.Join(root, "auth", "auth.thrift"))
	assert.True(t, IncludeExists(current, "shared/v1/identity.thrift"))
	assert.False(t, IncludeExists(current, "shared/v1/missing.thrift"))
}

func Test_IncludeURI(t *testing.T) {
	t.Cleanup(func() { SetIncludeRoots(nil) })
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	packageDir := filepath.Dir(file)
	repositoryRoot := filepath.Clean(filepath.Join(packageDir, "..", ".."))

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
		{
			name: "case4",
			args: args{
				cur:         uri.File("/tmp/workspace/app.subpath.thrift"),
				includePath: "user.subpath.thrift",
			},
			want: uri.File("/tmp/workspace/user.subpath.thrift"),
		},
		{
			name: "relative include takes precedence",
			args: args{
				cur:         uri.File(filepath.Join(packageDir, "utils.go")),
				includePath: "utils_test.go",
			},
			want: uri.File(filepath.Join(packageDir, "utils_test.go")),
		},
		{
			name: "workspace root include",
			args: args{
				cur:         uri.File(filepath.Join(packageDir, "missing.thrift")),
				includePath: "go.mod",
			},
			want: uri.File(filepath.Join(repositoryRoot, "go.mod")),
		},
		{
			name: "missing include falls back to relative path",
			args: args{
				cur:         uri.File(filepath.Join(packageDir, "missing.thrift")),
				includePath: "does-not-exist.thrift",
			},
			want: uri.File(filepath.Join(packageDir, "does-not-exist.thrift")),
		},
	}
	SetIncludeRoots([]string{repositoryRoot})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IncludeURI(tt.args.cur, tt.args.includePath))
		})
	}
}

func TestGetIncludePath(t *testing.T) {
	file := `include "../../user.thrift"
include "../../user.extra.thrift"
include "shared/v1/identity.thrift"
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
		{
			name: "case",
			args: args{
				ast:         ast.(*parser.Document),
				includeName: "user.extra",
			},
			want: "../../user.extra.thrift",
		},
		{
			name: "directory alias",
			args: args{
				ast:         ast.(*parser.Document),
				includeName: "shared",
			},
			want: "shared/v1/identity.thrift",
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
		{
			name: "file name with .",
			args: args{
				file: uri.New("/tmp/base.subpath.thrift"),
			},
			want: "base.subpath",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetIncludeName(tt.args.file))
		})
	}
}

func TestIncludeNames(t *testing.T) {
	type args struct {
		cur      uri.URI
		includes []*parser.Include
	}
	tests := []struct {
		name             string
		args             args
		wantIncludeNames []string
	}{
		{
			name: "case 1",
			args: args{
				cur: uri.New("/tmp/app.thrift"),
				includes: []*parser.Include{
					{
						Path: &parser.Literal{
							Value: &parser.LiteralValue{
								Text: "../../base.sub.thrift",
							},
						},
					},
					{
						Path: &parser.Literal{
							Value: &parser.LiteralValue{
								Text: "user.sub.thrift",
							},
						},
					},
					{
						Path: &parser.Literal{
							Value: &parser.LiteralValue{
								Text: "app.thrift",
							},
						},
					},
				},
			},
			wantIncludeNames: []string{
				"base.sub",
				"user.sub",
				"app",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantIncludeNames, IncludeNames(tt.args.cur, tt.args.includes))
		})
	}
}

func TestParseIdent(t *testing.T) {
	type args struct {
		cur        uri.URI
		includes   []*parser.Include
		identifier string
	}
	tests := []struct {
		name        string
		args        args
		wantInclude string
		wantIdent   string
	}{
		{
			name: "directory alias",
			args: args{
				cur: uri.New("/tmp/auth.thrift"),
				includes: []*parser.Include{
					{
						Path: &parser.Literal{
							Value: &parser.LiteralValue{
								Text: "shared/v1/identity.thrift",
							},
						},
					},
				},
				identifier: "shared.AccountIdentity",
			},
			wantInclude: "shared",
			wantIdent:   "AccountIdentity",
		},
		{
			name: "case 1",
			args: args{
				cur: uri.New("/tmp/app.thrift"),
				includes: []*parser.Include{
					{
						Path: &parser.Literal{
							Value: &parser.LiteralValue{
								Text: "user.sub.thrift",
							},
						},
					},
					{
						Path: &parser.Literal{
							Value: &parser.LiteralValue{
								Text: "user.thrift",
							},
						},
					},
				},
				identifier: "user.Name",
			},
			wantInclude: "user",
			wantIdent:   "Name",
		},
		{
			name: "case 2",
			args: args{
				cur: uri.New("/tmp/app.thrift"),
				includes: []*parser.Include{
					{
						Path: &parser.Literal{
							Value: &parser.LiteralValue{
								Text: "user.sub.thrift",
							},
						},
					},
					{
						Path: &parser.Literal{
							Value: &parser.LiteralValue{
								Text: "user.thrift",
							},
						},
					},
				},
				identifier: "user.sub.Name",
			},
			wantInclude: "user.sub",
			wantIdent:   "Name",
		},
		{
			name: "case 3",
			args: args{
				cur: uri.New("/tmp/app.thrift"),
				includes: []*parser.Include{
					{
						Path: &parser.Literal{
							Value: &parser.LiteralValue{
								Text: "user.thrift",
							},
						},
					},
				},
				identifier: "user.sub.Name",
			},
			wantInclude: "user",
			wantIdent:   "sub.Name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInclude, gotIdent := ParseIdent(tt.args.cur, tt.args.includes, tt.args.identifier)
			assert.Equal(t, tt.wantInclude, gotInclude)
			assert.Equal(t, tt.wantIdent, gotIdent)
		})
	}
}
