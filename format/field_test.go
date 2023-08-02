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
						Comments: []*parser.Comment{
							{
								Text: "/* aaa */",
							},
						},
					},
					KeyType: &parser.FieldType{
						TypeName: &parser.TypeName{
							Name: "string",
							Comments: []*parser.Comment{
								{
									Text: "/* aaa */",
								},
							},
						},
					},
					ValueType: &parser.FieldType{
						TypeName: &parser.TypeName{
							Name: "i32",
							Comments: []*parser.Comment{
								{
									Text: "/* aaa */",
								},
							},
						},
					},
					Annotations: &parser.Annotations{
						LParKeyword: &parser.LParKeyword{
							Keyword: parser.Keyword{
								Comments: []*parser.Comment{
									{
										Text: "/* aaa */",
									},
								},
								Literal: &parser.KeywordLiteral{
									Text: "(",
								},
							},
						},
						RParKeyword: &parser.RParKeyword{
							Keyword: parser.Keyword{
								Comments: []*parser.Comment{
									{
										Text: "/* aaa */",
									},
								},
								Literal: &parser.KeywordLiteral{
									Text: ")",
								},
							},
						},
						Annotations: []*parser.Annotation{
							{
								Identifier: &parser.Identifier{
									Name: &parser.IdentifierName{
										Text: "key1",
									},
									Comments: []*parser.Comment{
										{
											Text: "/* aaa */",
										},
									},
								},
								Value: &parser.Literal{
									Value: "value1",
									Quote: "'",
									Comments: []*parser.Comment{
										{
											Text: "/* aaa */",
										},
									},
								},
								EqualKeyword: &parser.EqualKeyword{
									Keyword: parser.Keyword{
										Literal: &parser.KeywordLiteral{
											Text: "=",
										},
										Comments: []*parser.Comment{
											{
												Text: "/* aaa */",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: "/* aaa */ map</* aaa */ string,/* aaa */ i32> /* aaa */ (/* aaa */ key1 /* aaa */ = /* aaa */ 'value1'/* aaa */ )",
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
