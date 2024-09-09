package diagnostic

import (
	"context"
	"sort"
	"testing"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func Test_SemanticAnalysis_Diagnostic(t *testing.T) {

	file1 := `struct Student {
	1: required string name,
	2: required User user1,
	3: required Student user2,
}
// line 6
	struct Student {}
// line 8
union Test {
	1: required string name,
	2: required User user1,
	3: required Student user2 = TestEnum.User, // enum doesn't exist
}
// line 14
exception TestError {
	1: required string name,
	2: required User user1,
	3: required Student user2,
}
// line 20
service TestService {
	Student Get(1: User user1) throws(1: TestError err1, 2: DoesNotExistError err2)
}
// line 24
struct TestContainer {
	1: required list<Student> Students
	2: required i32 failed1 = true
	3: required i32 failed2 = ""
	4: required string failed3 = true
	5: required string failed4 = 71

	100: required i32 user2 = 1
	101: required i64 user3 = 2
	102: required bool isUser = true
}
// line 36
struct TestUUID {
	1: required uuid id
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
		args      args
		want      DiagnosticResult
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "case 1",
			args: args{
				ctx: context.TODO(),
				ss:  ss,
				changeFiles: []uri.URI{
					"file:///tmp/user.thrift",
				},
			},
			want: DiagnosticResult{
				"file:///tmp/user.thrift": {
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      2,
								Character: 13,
							},
							End: protocol.Position{
								Line:      2,
								Character: 18,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field type doesn't exist",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      6,
								Character: 8,
							},
							End: protocol.Position{
								Line:      6,
								Character: 15,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "struct name conflict with other struct",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      10,
								Character: 13,
							},
							End: protocol.Position{
								Line:      10,
								Character: 18,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field type doesn't exist",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      11,
								Character: 29,
							},
							End: protocol.Position{
								Line:      11,
								Character: 42,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "default value doesn't exist",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      16,
								Character: 13,
							},
							End: protocol.Position{
								Line:      16,
								Character: 18,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field type doesn't exist",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      21,
								Character: 16,
							},
							End: protocol.Position{
								Line:      21,
								Character: 21,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field type doesn't exist",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      21,
								Character: 57,
							},
							End: protocol.Position{
								Line:      21,
								Character: 75,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "field type doesn't exist",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      26,
								Character: 27,
							},
							End: protocol.Position{
								Line:      26,
								Character: 31,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "expect i32 but got bool",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      27,
								Character: 27,
							},
							End: protocol.Position{
								Line:      27,
								Character: 29,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "expect i32 but got string",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      28,
								Character: 30,
							},
							End: protocol.Position{
								Line:      28,
								Character: 34,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "expect string but got bool",
					},
					{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      29,
								Character: 30,
							},
							End: protocol.Position{
								Line:      29,
								Character: 32,
							},
						},
						Severity: protocol.DiagnosticSeverityError,
						Source:   "thrift-ls",
						Message:  "expect string but got i64",
					},
				},
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SemanticAnalysis{}
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
