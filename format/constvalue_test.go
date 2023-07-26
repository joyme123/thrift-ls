package format

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func TestMustFormatConstValue(t *testing.T) {
	type args struct {
		cv *parser.ConstValue
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "map",
			args: args{
				cv: &parser.ConstValue{
					TypeName: "map",
					Value: []*parser.ConstValue{
						{
							TypeName: "pair",
							Key: &parser.ConstValue{
								TypeName: "string",
								Value:    "key",
							},
							Value: &parser.ConstValue{
								TypeName: "string",
								Value:    "value",
							},
						},
						{
							TypeName: "pair",
							Key: &parser.ConstValue{
								TypeName: "string",
								Value:    "key2",
							},
							Value: &parser.ConstValue{
								TypeName: "string",
								Value:    "value2",
							},
						},
					},
				},
			},
			want: `{"key": "value", "key2": "value2"}`,
		},
		{
			name: "list",
			args: args{
				cv: &parser.ConstValue{
					TypeName: "list",
					Value: []*parser.ConstValue{
						{
							TypeName: "string",
							Value:    "value1",
						},
						{
							TypeName: "string",
							Value:    "value2",
						},
					},
				},
			},
			want: `["value1", "value2"]`,
		},
		{
			name: "i64",
			args: args{
				cv: &parser.ConstValue{
					TypeName:    "i64",
					Value:       1,
					ValueInText: "1",
				},
			},
			want: "1",
		},
		{
			name: "i64 in hex",
			args: args{
				cv: &parser.ConstValue{
					TypeName:    "i64",
					Value:       26,
					ValueInText: "0x1a",
				},
			},
			want: "0x1a",
		},
		{
			name: "i64 in oct",
			args: args{
				cv: &parser.ConstValue{
					TypeName:    "i64",
					Value:       1,
					ValueInText: "0o1",
				},
			},
			want: "0o1",
		},
		{
			name: "double",
			args: args{
				cv: &parser.ConstValue{
					TypeName:    "i64",
					Value:       1e11,
					ValueInText: "1E11",
				},
			},
			want: "1E11",
		},
		{
			name: "identifier",
			args: args{
				cv: &parser.ConstValue{
					TypeName: "identifier",
					Value:    "User.Name",
				},
			},
			want: "User.Name",
		},
		{
			name: "literal",
			args: args{
				cv: &parser.ConstValue{
					TypeName: "string",
					Value:    "value",
				},
			},
			want: `"value"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MustFormatConstValue(tt.args.cv))
		})
	}
}
