package diagnostic

import (
	"context"
	"sort"
	"testing"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/memoize"
	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/uri"
)

func Test_cycleDetect(t *testing.T) {

	includesMap := map[uri.URI][]Include{
		"/user.thrift": {
			Include{file: "/goods.thrift"},
			Include{file: "/address.thrift"},
		},
		"/goods.thrift":   {Include{file: "/user.thrift"}},
		"/address.thrift": {Include{file: "/user.thrift"}},
	}

	type args struct {
		includesMap *map[uri.URI][]Include
	}
	tests := []struct {
		name string
		args args
		want []CyclePair
	}{
		{
			name: "cycle",
			args: args{
				includesMap: &includesMap,
			},
			want: []CyclePair{
				{
					file: "/user.thrift",
					include: Include{
						file: "/goods.thrift",
					},
				},
				{
					file:    "/goods.thrift",
					include: Include{file: "/user.thrift"},
				},
				{
					file:    "/user.thrift",
					include: Include{file: "/address.thrift"},
				},
				{
					file:    "/address.thrift",
					include: Include{file: "/user.thrift"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.SliceStable(tt.want, func(i, j int) bool {
				if tt.want[i].file == tt.want[j].file {
					return tt.want[i].include.file < tt.want[j].include.file
				}
				return tt.want[i].file < tt.want[j].file
			})

			got := cycleDetect(tt.args.includesMap)
			sort.SliceStable(got, func(i, j int) bool {
				if got[i].file == got[j].file {
					return got[i].include.file < got[j].include.file
				}
				return got[i].file < got[j].file
			})

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getIncludes(t *testing.T) {
	file1 := `include "./test/goods.thrift"
include "./test/address.thrift"`
	file2 := `include "../user.thrift"`
	file3 := `include "../user.thrift"`

	ss := buildSnapshotForTest([]*cache.FileChange{
		{
			URI:     "file:///tmp/user.thrift",
			Version: 0,
			Content: []byte(file1),
			From:    cache.FileChangeTypeDidOpen,
		},
		{
			URI:     "file:///tmp/test/goods.thrift",
			Version: 0,
			Content: []byte(file2),
			From:    cache.FileChangeTypeDidOpen,
		},
		{
			URI:     "file:///tmp/test/address.thrift",
			Version: 0,
			Content: []byte(file3),
			From:    cache.FileChangeTypeDidOpen,
		},
	})

	expectIncludeMap := map[uri.URI][]Include{
		"file:///tmp/user.thrift": {
			Include{
				file: "file:///tmp/test/goods.thrift",
				include: &parser.Include{
					IncludeKeyword: &parser.IncludeKeyword{
						Keyword: parser.Keyword{
							Literal: &parser.KeywordLiteral{
								Text:    "include",
								BadNode: false,
								Location: parser.Location{
									StartPos: parser.Position{
										Line:   1,
										Col:    1,
										Offset: 0,
									},
									EndPos: parser.Position{
										Line:   1,
										Col:    8,
										Offset: 7,
									},
								},
							},
							BadNode: false,
							Location: parser.Location{
								StartPos: parser.Position{
									Line:   1,
									Col:    1,
									Offset: 0,
								},
								EndPos: parser.Position{
									Line:   1,
									Col:    9,
									Offset: 8,
								},
							},
						},
					},
					Path: &parser.Literal{
						Value:   "./test/goods.thrift",
						BadNode: false,
						Quote:   "\"",
						Location: parser.NewLocationFromPos(
							parser.Position{
								Line:   1,
								Col:    9,
								Offset: 8,
							},
							parser.Position{
								Line:   1,
								Col:    30,
								Offset: 29,
							},
						),
					},
					Location: parser.NewLocationFromPos(
						parser.Position{
							Line:   1,
							Col:    1,
							Offset: 0,
						},
						parser.Position{
							Line:   1,
							Col:    30,
							Offset: 29,
						},
					),
				},
			},
			Include{
				file: "file:///tmp/test/address.thrift",
				include: &parser.Include{
					IncludeKeyword: &parser.IncludeKeyword{
						Keyword: parser.Keyword{
							Literal: &parser.KeywordLiteral{
								Text:    "include",
								BadNode: false,
								Location: parser.Location{
									StartPos: parser.Position{
										Line:   2,
										Col:    1,
										Offset: 30,
									},
									EndPos: parser.Position{
										Line:   2,
										Col:    8,
										Offset: 37,
									},
								},
							},
							BadNode: false,
							Location: parser.Location{
								StartPos: parser.Position{
									Line:   2,
									Col:    1,
									Offset: 30,
								},
								EndPos: parser.Position{
									Line:   2,
									Col:    9,
									Offset: 38,
								},
							},
						},
					},
					Path: &parser.Literal{
						Value:   "./test/address.thrift",
						Quote:   "\"",
						BadNode: false,
						Location: parser.NewLocationFromPos(
							parser.Position{
								Line:   2,
								Col:    9,
								Offset: 38,
							},
							parser.Position{
								Line:   2,
								Col:    32,
								Offset: 61,
							},
						),
					},
					Location: parser.NewLocationFromPos(
						parser.Position{
							Line:   2,
							Col:    1,
							Offset: 30,
						},
						parser.Position{
							Line:   2,
							Col:    32,
							Offset: 61,
						},
					),
				},
			},
		},
		"file:///tmp/test/goods.thrift": {
			Include{
				file: "file:///tmp/user.thrift",
				include: &parser.Include{
					IncludeKeyword: &parser.IncludeKeyword{
						Keyword: parser.Keyword{
							Literal: &parser.KeywordLiteral{
								Text:    "include",
								BadNode: false,
								Location: parser.Location{
									StartPos: parser.Position{
										Line:   1,
										Col:    1,
										Offset: 0,
									},
									EndPos: parser.Position{
										Line:   1,
										Col:    8,
										Offset: 7,
									},
								},
							},
							BadNode: false,
							Location: parser.Location{
								StartPos: parser.Position{
									Line:   1,
									Col:    1,
									Offset: 0,
								},
								EndPos: parser.Position{
									Line:   1,
									Col:    9,
									Offset: 8,
								},
							},
						},
					},
					Path: &parser.Literal{
						Value:   "../user.thrift",
						Quote:   "\"",
						BadNode: false,
						Location: parser.NewLocationFromPos(
							parser.Position{
								Line:   1,
								Col:    9,
								Offset: 8,
							},
							parser.Position{
								Line:   1,
								Col:    25,
								Offset: 24,
							},
						),
					},
					Location: parser.NewLocationFromPos(
						parser.Position{
							Line:   1,
							Col:    1,
							Offset: 0,
						},
						parser.Position{
							Line:   1,
							Col:    25,
							Offset: 24,
						},
					),
				},
			},
		},
		"file:///tmp/test/address.thrift": {
			Include{
				file: "file:///tmp/user.thrift",
				include: &parser.Include{
					IncludeKeyword: &parser.IncludeKeyword{
						Keyword: parser.Keyword{
							Literal: &parser.KeywordLiteral{
								Text:    "include",
								BadNode: false,
								Location: parser.Location{
									StartPos: parser.Position{
										Line:   1,
										Col:    1,
										Offset: 0,
									},
									EndPos: parser.Position{
										Line:   1,
										Col:    8,
										Offset: 7,
									},
								},
							},
							BadNode: false,
							Location: parser.Location{
								StartPos: parser.Position{
									Line:   1,
									Col:    1,
									Offset: 0,
								},
								EndPos: parser.Position{
									Line:   1,
									Col:    9,
									Offset: 8,
								},
							},
						},
					},
					Path: &parser.Literal{
						Value:   "../user.thrift",
						Quote:   "\"",
						BadNode: false,
						Location: parser.NewLocationFromPos(
							parser.Position{
								Line:   1,
								Col:    9,
								Offset: 8,
							},
							parser.Position{
								Line:   1,
								Col:    25,
								Offset: 24,
							},
						),
					},
					Location: parser.NewLocationFromPos(
						parser.Position{
							Line:   1,
							Col:    1,
							Offset: 0,
						},
						parser.Position{
							Line:   1,
							Col:    25,
							Offset: 24,
						},
					),
				},
			},
		},
	}
	includeMap := make(map[uri.URI][]Include)

	type args struct {
		ctx         context.Context
		ss          *cache.Snapshot
		file        uri.URI
		includesMap *map[uri.URI][]Include
	}
	tests := []struct {
		name      string
		args      args
		want      *map[uri.URI][]Include
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				ctx:         context.TODO(),
				ss:          ss,
				file:        "file:///tmp/user.thrift",
				includesMap: &includeMap,
			},
			want:      &expectIncludeMap,
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, getIncludes(tt.args.ctx, tt.args.ss, tt.args.file, tt.args.includesMap))

			assert.Equal(t, tt.want, tt.args.includesMap)
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
