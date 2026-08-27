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
									Value: &parser.LiteralValue{
										Text: "value1",
									},
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

func Test_FormatFieldReference(t *testing.T) {
	tests := []struct {
		name string
		doc  string
		want string
	}{
		{
			name: "aligned struct fields and function argument",
			doc: `struct Node {
1: optional Node & next
2: required i64 value
3: optional list<Node> & children (cpp.ref = "true")
}

service Tree {
void insert(1: Node & node) throws (1: TreeError & err)
}`,
			want: `struct Node {
    1: optional Node       & next
    2: required i64        value
    3: optional list<Node> & children (cpp.ref = "true")
}

service Tree {
    void insert(1: Node & node) throws (1: TreeError & err)
}`,
		},
		{
			name: "marker is normalized to a single surrounding space",
			doc:  "struct A {\n1: Node&next\n}",
			want: "struct A {\n    1: Node & next\n}",
		},
		{
			name: "newline between marker and name is joined",
			doc:  "struct A {\n1: Node &\nnext\n}",
			want: "struct A {\n    1: Node & next\n}",
		},
		{
			name: "comment before the marker is kept",
			doc:  "struct A {\n1: Node /* a */ & next\n}",
			want: "struct A {\n    1: Node /* a */ & next\n}",
		},
		{
			name: "comment between marker and name is kept",
			doc:  "struct A {\n1: Node & /* retain */ next\n}",
			want: "struct A {\n    1: Node & /* retain */ next\n}",
		},
		{
			name: "comments on both sides of the marker are kept",
			doc:  "struct A {\n1: Node /* a */ & /* b */ next\n}",
			want: "struct A {\n    1: Node /* a */ & /* b */ next\n}",
		},
		{
			name: "comment before the name without a marker is kept",
			doc:  "struct A {\n1: Node /* c */ next\n}",
			want: "struct A {\n    1: Node /* c */ next\n}",
		},
		{
			name: "comment in a function argument is kept",
			doc:  "service S {\nvoid f(1: Node & /* c */ n)\n}",
			want: "service S {\n    void f(1: Node & /* c */ n)\n}",
		},
		{
			name: "marker with default value",
			doc:  "struct A {\n1: optional Node & next = 1\n}",
			want: "struct A {\n    1: optional Node & next = 1\n}",
		},
		{
			name: "marker in union and exception fields",
			doc:  "union U {\n1: U & other\n}\n\nexception E {\n1: E & other\n}",
			want: "union U {\n    1: U & other\n}\n\nexception E {\n    1: E & other\n}",
		},
	}

	FieldLineComma = FieldLineCommaDisable

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := parser.Parse("test.thrift", []byte(tt.doc))
			assert.NoError(t, err)
			assert.NotNil(t, ast)

			formated, err := FormatDocument(ast.(*parser.Document))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, formated)

			// self validation reparses the formatted result and compares it against
			// the original ast, so it fails if the reference marker or a comment
			// attached to the field identifier is dropped
			_, err = FormatDocumentWithValidation(ast.(*parser.Document), true)
			assert.NoError(t, err)

			// formatting is a fixed point
			reparsed, err := parser.Parse("test.thrift", []byte(formated))
			assert.NoError(t, err)
			again, err := FormatDocumentWithValidation(reparsed.(*parser.Document), true)
			assert.NoError(t, err)
			assert.Equal(t, formated, again)
		})
	}
}
