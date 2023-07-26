package format

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func TestMustFormatFieldType(t *testing.T) {
	type args struct {
		ft *parser.FieldType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "identifier type",
			args: args{
				ft: &parser.FieldType{
					TypeName: &parser.TypeName{
						Name: "User",
					},
				},
			},
			want: "User",
		},
		{
			name: "map type",
			args: args{
				ft: &parser.FieldType{
					TypeName: &parser.TypeName{
						Name: "map",
					},
					KeyType: &parser.FieldType{
						TypeName: &parser.TypeName{
							Name: "string",
						},
					},
					ValueType: &parser.FieldType{
						TypeName: &parser.TypeName{
							Name: "i32",
						},
					},
				},
			},
			want: "map<string,i32>",
		},
		{
			name: "set type",
			args: args{
				ft: &parser.FieldType{
					TypeName: &parser.TypeName{
						Name: "set",
					},
					KeyType: &parser.FieldType{
						TypeName: &parser.TypeName{
							Name: "string",
						},
					},
				},
			},
			want: "set<string>",
		},
		{
			name: "list type",
			args: args{
				ft: &parser.FieldType{
					TypeName: &parser.TypeName{
						Name: "list",
					},
					KeyType: &parser.FieldType{
						TypeName: &parser.TypeName{
							Name: "string",
						},
					},
				},
			},
			want: "list<string>",
		},
		{
			name: "embedding type",
			args: args{
				ft: &parser.FieldType{
					TypeName: &parser.TypeName{
						Name: "list",
					},
					KeyType: &parser.FieldType{
						TypeName: &parser.TypeName{
							Name: "map",
						},
						KeyType: &parser.FieldType{
							TypeName: &parser.TypeName{
								Name: "string",
							},
						},
						ValueType: &parser.FieldType{
							TypeName: &parser.TypeName{
								Name: "string",
							},
						},
					},
				},
			},
			want: "list<map<string,string>>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MustFormatFieldType(tt.args.ft))
		})
	}
}
