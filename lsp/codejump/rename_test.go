package codejump

import (
	"context"
	"testing"

	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/uri"
)

func TestPrepareRename(t *testing.T) {

	file1 := `struct Test {
  1: required string name,
  2: required string email,
  3: required string test1,
  4: required int test2,
}

union Test2 {
  1: required string name,
  2: required string email,
  3: required string test1,
  4: required int test2,
}

enum Test3 {
  ONE = 1,
  TWO
}

exception Error1 {
  1: required string name,
  2: required string email,
  3: required string test1,
  4: required int test2,
}

typedef string UserType
const string DefaultName="nickname"`

	file2 := `include "user.thrift"
service Demo {
  user.Test Api(1:user.Test2 arg1, 2:user.Test3 arg2) throws (1:user.Error1 err)
  list<user.UserType> UserTypes(1:user.Test3 arg1=user.Test3.TWO, 2:string arg2=user.DefaultName)
}

typedef user.UserType UserKind
const user.UserType usermale = "male"
const UserKind kind = "1"
`

	ss := cache.BuildSnapshotForTest([]*cache.FileChange{
		{
			URI:     "file:///tmp/user.thrift",
			Version: 0,
			Content: []byte(file1),
			From:    cache.FileChangeTypeDidOpen,
		},
		{
			URI:     "file:///tmp/api.thrift",
			Version: 0,
			Content: []byte(file2),
			From:    cache.FileChangeTypeDidOpen,
		},
	})

	type args struct {
		ctx  context.Context
		ss   *cache.Snapshot
		file uri.URI
		pos  protocol.Position
	}
	tests := []struct {
		name      string
		args      args
		wantRes   *protocol.Range
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "case struct",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/user.thrift",
				pos: protocol.Position{
					Line:      0,
					Character: 7,
				},
			},
			wantRes: &protocol.Range{
				Start: protocol.Position{
					Line:      0,
					Character: 7,
				},
				End: protocol.Position{
					Line:      0,
					Character: 11,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case union",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/user.thrift",
				pos: protocol.Position{
					Line:      7,
					Character: 6,
				},
			},
			wantRes: &protocol.Range{
				Start: protocol.Position{
					Line:      7,
					Character: 6,
				},
				End: protocol.Position{
					Line:      7,
					Character: 11,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case enum",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/user.thrift",
				pos: protocol.Position{
					Line:      14,
					Character: 5,
				},
			},
			wantRes: &protocol.Range{
				Start: protocol.Position{
					Line:      14,
					Character: 5,
				},
				End: protocol.Position{
					Line:      14,
					Character: 10,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case exception",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/user.thrift",
				pos: protocol.Position{
					Line:      19,
					Character: 10,
				},
			},
			wantRes: &protocol.Range{
				Start: protocol.Position{
					Line:      19,
					Character: 10,
				},
				End: protocol.Position{
					Line:      19,
					Character: 16,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "typedef",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/user.thrift",
				pos: protocol.Position{
					Line:      26,
					Character: 15,
				},
			},
			wantRes: &protocol.Range{
				Start: protocol.Position{
					Line:      26,
					Character: 15,
				},
				End: protocol.Position{
					Line:      26,
					Character: 23,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "const",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/user.thrift",
				pos: protocol.Position{
					Line:      27,
					Character: 13,
				},
			},
			wantRes: &protocol.Range{
				Start: protocol.Position{
					Line:      27,
					Character: 13,
				},
				End: protocol.Position{
					Line:      27,
					Character: 24,
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := PrepareRename(tt.args.ctx, tt.args.ss, tt.args.file, tt.args.pos)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestRename(t *testing.T) {

	file1 := `struct Test {
  1: required string name,
  2: required string email,
  3: required string test1,
  4: required int test2,
}

union Test2 {
  1: required string name,
  2: required string email,
  3: required string test1,
  4: required int test2,
}

enum Test3 {
  ONE = 1,
  TWO
}

exception Error1 {
  1: required string name,
  2: required string email,
  3: required string test1,
  4: required int test2,
}

typedef string UserType
const string DefaultName="nickname"`

	file2 := `include "user.thrift"
service Demo {
  user.Test Api(1:user.Test2 arg1, 2:user.Test3 arg2) throws (1:user.Error1 err)
  list<user.UserType> UserTypes(1:user.Test3 arg1=user.Test3.TWO, 2:string arg2=user.DefaultName)
}

typedef user.UserType UserKind
const user.UserType usermale = "male"
const UserKind kind = "1"
`

	ss := cache.BuildSnapshotForTest([]*cache.FileChange{
		{
			URI:     "file:///tmp/user.thrift",
			Version: 0,
			Content: []byte(file1),
			From:    cache.FileChangeTypeDidOpen,
		},
		{
			URI:     "file:///tmp/api.thrift",
			Version: 0,
			Content: []byte(file2),
			From:    cache.FileChangeTypeDidOpen,
		},
	})

	type args struct {
		ctx     context.Context
		ss      *cache.Snapshot
		file    uri.URI
		pos     protocol.Position
		newText string
	}
	tests := []struct {
		name      string
		args      args
		wantRes   *protocol.WorkspaceEdit
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "case struct",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/user.thrift",
				pos: protocol.Position{
					Line:      0,
					Character: 7,
				},
				newText: "newtext",
			},
			wantRes: &protocol.WorkspaceEdit{
				Changes: map[uri.URI][]protocol.TextEdit{
					"file:///tmp/user.thrift": []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 7,
								},
								End: protocol.Position{
									Line:      0,
									Character: 11,
								},
							},
							NewText: "newtext",
						},
					},
					"file:///tmp/api.thrift": []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      2,
									Character: 2,
								},
								End: protocol.Position{
									Line:      2,
									Character: 11,
								},
							},
							NewText: "user.newtext",
						},
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case union",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/user.thrift",
				pos: protocol.Position{
					Line:      7,
					Character: 6,
				},
				newText: "newtext",
			},
			wantRes: &protocol.WorkspaceEdit{
				Changes: map[uri.URI][]protocol.TextEdit{
					"file:///tmp/user.thrift": []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      7,
									Character: 6,
								},
								End: protocol.Position{
									Line:      7,
									Character: 11,
								},
							},
							NewText: "newtext",
						},
					},
					"file:///tmp/api.thrift": []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      2,
									Character: 18,
								},
								End: protocol.Position{
									Line:      2,
									Character: 28,
								},
							},
							NewText: "user.newtext",
						},
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case enum",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/user.thrift",
				pos: protocol.Position{
					Line:      14,
					Character: 5,
				},
				newText: "newtext",
			},
			wantRes: &protocol.WorkspaceEdit{
				Changes: map[uri.URI][]protocol.TextEdit{
					"file:///tmp/user.thrift": []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      14,
									Character: 5,
								},
								End: protocol.Position{
									Line:      14,
									Character: 10,
								},
							},
							NewText: "newtext",
						},
					},
					"file:///tmp/api.thrift": []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      2,
									Character: 37,
								},
								End: protocol.Position{
									Line:      2,
									Character: 47,
								},
							},
							NewText: "user.newtext",
						},
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      3,
									Character: 34,
								},
								End: protocol.Position{
									Line:      3,
									Character: 44,
								},
							},
							NewText: "user.newtext",
						},
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case exception",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/user.thrift",
				pos: protocol.Position{
					Line:      19,
					Character: 10,
				},
				newText: "newtext",
			},
			wantRes: &protocol.WorkspaceEdit{
				Changes: map[uri.URI][]protocol.TextEdit{
					"file:///tmp/user.thrift": []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      19,
									Character: 10,
								},
								End: protocol.Position{
									Line:      19,
									Character: 16,
								},
							},
							NewText: "newtext",
						},
					},
					"file:///tmp/api.thrift": []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      2,
									Character: 64,
								},
								End: protocol.Position{
									Line:      2,
									Character: 75,
								},
							},
							NewText: "user.newtext",
						},
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "typedef",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/user.thrift",
				pos: protocol.Position{
					Line:      26,
					Character: 15,
				},
				newText: "newtext",
			},
			wantRes: &protocol.WorkspaceEdit{
				Changes: map[uri.URI][]protocol.TextEdit{
					"file:///tmp/user.thrift": []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      26,
									Character: 15,
								},
								End: protocol.Position{
									Line:      26,
									Character: 23,
								},
							},
							NewText: "newtext",
						},
					},
					"file:///tmp/api.thrift": []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      3,
									Character: 7,
								},
								End: protocol.Position{
									Line:      3,
									Character: 20,
								},
							},
							NewText: "user.newtext",
						},
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      6,
									Character: 8,
								},
								End: protocol.Position{
									Line:      6,
									Character: 21,
								},
							},
							NewText: "user.newtext",
						},
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      7,
									Character: 6,
								},
								End: protocol.Position{
									Line:      7,
									Character: 19,
								},
							},
							NewText: "user.newtext",
						},
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "const",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/user.thrift",
				pos: protocol.Position{
					Line:      27,
					Character: 13,
				},
				newText: "newtext",
			},
			wantRes: &protocol.WorkspaceEdit{
				Changes: map[uri.URI][]protocol.TextEdit{
					"file:///tmp/user.thrift": []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      27,
									Character: 13,
								},
								End: protocol.Position{
									Line:      27,
									Character: 24,
								},
							},
							NewText: "newtext",
						},
					},
					"file:///tmp/api.thrift": []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      3,
									Character: 80,
								},
								End: protocol.Position{
									Line:      3,
									Character: 96,
								},
							},
							NewText: "user.newtext",
						},
					},
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := Rename(tt.args.ctx, tt.args.ss, tt.args.file, tt.args.pos, tt.args.newText)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}
