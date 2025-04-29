package diagnostic

import (
	"context"
	"sort"
	"testing"

	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/uri"
)

func Test_FieldIDCheck_Diagnostic(t *testing.T) {
	file1 := `struct Test {
  1: required string name,
  1: required string email,
  0: required string test1,
  32768: required int test2,
}

union Test2 {
  1: required string name,
  1: required string email,
  0: required string test1,
  32768: required int test2,
}

exception Test3 {
  1: required string name,
  1: required string email,
  0: required string test1,
  32768: required int test2,
} // line 20

service Demo {
  Test Api1(0:Test arg, 1: Test2 arg1, 1: Test2 arg2, 32768: int arg4),
  Test Api2(1: Test2 arg1) throws (0:Test3 err, 1:Test3 err1, 1:Test3 err2, 32768:Test3 err4)
}
`

	ss := buildSnapshotForTest([]*cache.FileChange{
		{
			URI:     "file:///tmp/user.thrift",
			Version: 0,
			Content: []byte(file1),
			From:    cache.FileChangeTypeDidOpen,
		},
	})
	type args struct {
		ctx         context.Context
		ss          *cache.Snapshot
		changeFiles []uri.URI
	}
	tests := []struct {
		name      string
		c         *FieldIDCheck
		args      args
		want      DiagnosticResult
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "case1",
			c:    &FieldIDCheck{},
			args: args{
				ctx: context.TODO(),
				ss:  ss,
				changeFiles: []uri.URI{
					"file:///tmp/user.thrift",
				},
			},
			want: DiagnosticResult{
				"file:///tmp/user.thrift": {
					// struct
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      1,
								Character: 2,
							},
							End: protocol.Position{
								Line:      1,
								Character: 3,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id conflict",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      2,
								Character: 2,
							},
							End: protocol.Position{
								Line:      2,
								Character: 3,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id conflict",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      3,
								Character: 2,
							},
							End: protocol.Position{
								Line:      3,
								Character: 3,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id should be a positive integer in [1, 32767]",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      4,
								Character: 2,
							},
							End: protocol.Position{
								Line:      4,
								Character: 7,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id should be a positive integer in [1, 32767]",
					},

					// union
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      8,
								Character: 2,
							},
							End: protocol.Position{
								Line:      8,
								Character: 3,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id conflict",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      9,
								Character: 2,
							},
							End: protocol.Position{
								Line:      9,
								Character: 3,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id conflict",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      10,
								Character: 2,
							},
							End: protocol.Position{
								Line:      10,
								Character: 3,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id should be a positive integer in [1, 32767]",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      11,
								Character: 2,
							},
							End: protocol.Position{
								Line:      11,
								Character: 7,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id should be a positive integer in [1, 32767]",
					},

					// exception
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      15,
								Character: 2,
							},
							End: protocol.Position{
								Line:      15,
								Character: 3,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id conflict",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      16,
								Character: 2,
							},
							End: protocol.Position{
								Line:      16,
								Character: 3,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id conflict",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      17,
								Character: 2,
							},
							End: protocol.Position{
								Line:      17,
								Character: 3,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id should be a positive integer in [1, 32767]",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      18,
								Character: 2,
							},
							End: protocol.Position{
								Line:      18,
								Character: 7,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id should be a positive integer in [1, 32767]",
					},

					// function params
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      22,
								Character: 12,
							},
							End: protocol.Position{
								Line:      22,
								Character: 13,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id should be a positive integer in [1, 32767]",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      22,
								Character: 24,
							},
							End: protocol.Position{
								Line:      22,
								Character: 25,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id conflict",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      22,
								Character: 39,
							},
							End: protocol.Position{
								Line:      22,
								Character: 40,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id conflict",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      22,
								Character: 54,
							},
							End: protocol.Position{
								Line:      22,
								Character: 59,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id should be a positive integer in [1, 32767]",
					},

					// function throws
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      23,
								Character: 35,
							},
							End: protocol.Position{
								Line:      23,
								Character: 36,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id should be a positive integer in [1, 32767]",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      23,
								Character: 48,
							},
							End: protocol.Position{
								Line:      23,
								Character: 49,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id conflict",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      23,
								Character: 62,
							},
							End: protocol.Position{
								Line:      23,
								Character: 63,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id conflict",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      23,
								Character: 76,
							},
							End: protocol.Position{
								Line:      23,
								Character: 81,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field id should be a positive integer in [1, 32767]",
					},
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &FieldIDCheck{}
			got, err := c.Diagnostic(tt.args.ctx, tt.args.ss, tt.args.changeFiles)

			for key := range got {
				sort.SliceStable(got[key], func(i, j int) bool {
					if got[key][i].Range.Start.Line == got[key][j].Range.Start.Line {
						return got[key][i].Range.Start.Character < got[key][j].Range.Start.Character
					}

					return got[key][i].Range.Start.Line < got[key][j].Range.Start.Line
				})
			}

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
