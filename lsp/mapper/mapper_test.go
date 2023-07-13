package mapper

import (
	"testing"

	"github.com/joyme123/thrift-ls/lsp/types"
	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/uri"
)

func TestMapper_LSPPosToParserPosition(t *testing.T) {
	type fields struct {
		fileURI uri.URI
		content []byte
	}
	type args struct {
		pos types.Position
	}

	content := `struct demo {
  1: required string name,
}`

	runeContent := `struct 😀😂 {
  1: required string name,
}`

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      parser.Position
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "ascii",
			fields: fields{
				fileURI: "test/test.thrift",
				content: []byte(content),
			},
			args: args{
				pos: types.Position{
					Line:      1,
					Character: 5,
				},
			},
			want: parser.Position{
				Line:   2,
				Col:    6,
				Offset: 19,
			},
			assertion: assert.NoError,
		},
		{
			name: "ascii line exceeded",
			fields: fields{
				fileURI: "test/test.thrift",
				content: []byte(content),
			},
			args: args{
				pos: types.Position{
					Line:      3,
					Character: 5,
				},
			},
			want:      parser.InvalidPosition,
			assertion: assert.Error,
		},
		{
			name: "ascii character exceeded",
			fields: fields{
				fileURI: "test/test.thrift",
				content: []byte(content),
			},
			args: args{
				pos: types.Position{
					Line:      1,
					Character: 27,
				},
			},
			want:      parser.InvalidPosition,
			assertion: assert.Error,
		},
		{
			name: "rune",
			fields: fields{
				fileURI: "test/test.thrift",
				content: []byte(runeContent),
			},
			args: args{
				pos: types.Position{
					Line:      0,
					Character: 12,
				},
			},
			want: parser.Position{
				Line:   1,
				Col:    11,
				Offset: 16,
			},
			assertion: assert.NoError,
		},
		{
			name: "rune line exceeded",
			fields: fields{
				fileURI: "test/test.thrift",
				content: []byte(content),
			},
			args: args{
				pos: types.Position{
					Line:      2,
					Character: 12,
				},
			},
			want:      parser.InvalidPosition,
			assertion: assert.Error,
		},
		{
			name: "rune character exceeded",
			fields: fields{
				fileURI: "test/test.thrift",
				content: []byte(content),
			},
			args: args{
				pos: types.Position{
					Line:      0,
					Character: 14,
				},
			},
			want:      parser.InvalidPosition,
			assertion: assert.Error,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			m := &Mapper{
				fileURI: tt.fields.fileURI,
				content: tt.fields.content,
			}
			got, err := m.LSPPosToParserPosition(tt.args.pos)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
