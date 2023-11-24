package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions_GetIndent(t *testing.T) {
	type fields struct {
		Write  bool
		Indent string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "4spaces",
			fields: fields{
				Indent: "4spaces",
			},
			want: "    ",
		},
		{
			name: "1space",
			fields: fields{
				Indent: "1spaces",
			},
			want: " ",
		},
		{
			name: "space",
			fields: fields{
				Indent: "space",
			},
			want: " ",
		},
		{
			name: "2tabs",
			fields: fields{
				Indent: "2tabs",
			},
			want: "\t\t",
		},
		{
			name: "1tab",
			fields: fields{
				Indent: "1tab",
			},
			want: "\t",
		},
		{
			name: "tab",
			fields: fields{
				Indent: "tab",
			},
			want: "\t",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				Write:  tt.fields.Write,
				Indent: tt.fields.Indent,
			}
			assert.Equal(t, tt.want, o.GetIndent())
		})
	}
}
