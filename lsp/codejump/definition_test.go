package codejump

import (
	"context"
	"testing"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func TestDefinition(t *testing.T) {
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
}`

	// user.extra.thrift
	file3 := `struct Test {}`

	file4 := `include "user.extra.thrift"
include "user.thrift"

struct Person {
  1: required user.extra.Test field1,
  2: required user.Test field2,
}`

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
		{
			URI:     "file:///tmp/user.extra.thrift",
			Version: 0,
			Content: []byte(file3),
			From:    cache.FileChangeTypeDidOpen,
		},
		{
			URI:     "file:///tmp/app.thrift",
			Version: 0,
			Content: []byte(file4),
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
		want      []protocol.Location
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "case struct",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/api.thrift",
				pos: protocol.Position{
					Line:      2,
					Character: 7,
				},
			},
			want: []protocol.Location{
				{
					URI: "file:///tmp/user.thrift",
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
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case union",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/api.thrift",
				pos: protocol.Position{
					Line:      2,
					Character: 23,
				},
			},
			want: []protocol.Location{
				{
					URI: "file:///tmp/user.thrift",
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
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case enum",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/api.thrift",
				pos: protocol.Position{
					Line:      2,
					Character: 42,
				},
			},
			want: []protocol.Location{
				{
					URI: "file:///tmp/user.thrift",
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
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case exceptions",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/api.thrift",
				pos: protocol.Position{
					Line:      2,
					Character: 69,
				},
			},
			want: []protocol.Location{
				{
					URI: "file:///tmp/user.thrift",
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
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case typedef",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/api.thrift",
				pos: protocol.Position{
					Line:      3,
					Character: 7,
				},
			},
			want: []protocol.Location{
				{
					URI: "file:///tmp/user.thrift",
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
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case enumvalue",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/api.thrift",
				pos: protocol.Position{
					Line:      3,
					Character: 61,
				},
			},
			want: []protocol.Location{
				{
					URI: "file:///tmp/user.thrift",
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      16,
							Character: 2,
						},
						End: protocol.Position{
							Line:      16,
							Character: 5,
						},
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case const",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/api.thrift",
				pos: protocol.Position{
					Line:      3,
					Character: 85,
				},
			},
			want: []protocol.Location{
				{
					URI: "file:///tmp/user.thrift",
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
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case include 1",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/app.thrift",
				pos: protocol.Position{
					Line:      4,
					Character: 25,
				},
			},
			want: []protocol.Location{
				{
					URI: "file:///tmp/user.extra.thrift",
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
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "case include 2",
			args: args{
				ctx:  context.TODO(),
				ss:   ss,
				file: "file:///tmp/app.thrift",
				pos: protocol.Position{
					Line:      5,
					Character: 19,
				},
			},
			want: []protocol.Location{
				{
					URI: "file:///tmp/user.thrift",
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
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Definition(tt.args.ctx, tt.args.ss, tt.args.file, tt.args.pos)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
