package format

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func TestMustFormatStruct(t *testing.T) {
	doc := `
// comments


/*
 * comments 
 */
struct test {
  /*
   * field 1
   */
  1: required string test,
}      (a.b = "c")          // endline comments
`

	type args struct {
		st *parser.Struct
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				st: func() *parser.Struct {
					ast, err := parser.Parse("test.thrift", []byte(doc))
					assert.NoError(t, err)
					return ast.(*parser.Document).Structs[0]
				}(),
			},
			want: `// comments

/*
 * comments
 */
struct test {
    /*
     * field 1
     */
    1: required string test,
} (a.b = "c") // endline comments
`,
		},
		{
			name: "test with CRLF",
			args: args{
				st: func() *parser.Struct {
					ast, err := parser.Parse("test.thrift", []byte("struct test {\r\n 1: required string test,\r\n 2: required string test2\r\n}"))
					assert.NoError(t, err)
					return ast.(*parser.Document).Structs[0]
				}(),
			},
			want: "struct test {\n    1: required string test,\n    2: required string test2\n}\n",
		},
		{
			name: "test with LF",
			args: args{
				st: func() *parser.Struct {
					ast, err := parser.Parse("test.thrift", []byte("struct test {\n 1: required string test,\n 2: required string test2\n}"))
					assert.NoError(t, err)
					return ast.(*parser.Document).Structs[0]
				}(),
			},
			want: "struct test {\n    1: required string test,\n    2: required string test2\n}\n",
		},
		{
			name: "test with additional LF",
			args: args{
				st: func() *parser.Struct {
					ast, err := parser.Parse("test.thrift", []byte("struct test {\n 1: required string test,\n\n 2: required string test2\n}"))
					assert.NoError(t, err)
					return ast.(*parser.Document).Structs[0]
				}(),
			},
			want: "struct test {\n    1: required string test,\n\n    2: required string test2\n}\n",
		},
		{
			name: "test struct annotation case 1",
			args: args{
				st: func() *parser.Struct {
					ast, err := parser.Parse("test.thrift", []byte(`
struct Foo {
  // 用户名列表
  2: list<string> strings (
    custom_tag_1 = "1"
    custom_tag_2 = "2"
  )
}
`))
					assert.NoError(t, err)
					return ast.(*parser.Document).Structs[0]
				}(),
			},
			want: `struct Foo {
    // 用户名列表
    2: list<string> strings (
        custom_tag_1 = "1"
        custom_tag_2 = "2"
    )
}
`,
		},
		{
			name: "test struct annotation case 2",
			args: args{
				st: func() *parser.Struct {
					ast, err := parser.Parse("test.thrift", []byte(`
struct Foo {
  // 用户名列表
  2: list<string> strings (
    custom_tag_1 = "1",
    // 自定义标签2
    custom_tag_2 = "2"
  )
}
`))
					assert.NoError(t, err)
					return ast.(*parser.Document).Structs[0]
				}(),
			},
			want: `struct Foo {
    // 用户名列表
    2: list<string> strings (
        custom_tag_1 = "1",
        // 自定义标签2
        custom_tag_2 = "2"
    )
}
`,
		},
		{
			name: "test struct annotation case 3",
			args: args{
				st: func() *parser.Struct {
					ast, err := parser.Parse("test.thrift", []byte(`
struct Foo {
  // 用户名列表
  2: list<string> strings (custom_tag_1 = "1", /* 自定义标签2 */    custom_tag_2 = "2")
}
`))
					assert.NoError(t, err)
					return ast.(*parser.Document).Structs[0]
				}(),
			},
			want: `struct Foo {
    // 用户名列表
    2: list<string> strings (custom_tag_1 = "1", /* 自定义标签2 */ custom_tag_2 = "2")
}
`,
		},
		{
			name: "test struct annotation case 4",
			args: args{
				st: func() *parser.Struct {
					ast, err := parser.Parse("test.thrift", []byte(`
struct Foo {
  // 用户名列表
  2: list<string> strings (
  	custom_tag_1 = "1", /* 自定义标签2 */ custom_tag_2 = "2"
  )
}
`))
					assert.NoError(t, err)
					return ast.(*parser.Document).Structs[0]
				}(),
			},
			want: `struct Foo {
    // 用户名列表
    2: list<string> strings (
        custom_tag_1 = "1", /* 自定义标签2 */ custom_tag_2 = "2"
    )
}
`,
		},
		{
			name: "test struct annotation case 5",
			args: args{
				st: func() *parser.Struct {
					ast, err := parser.Parse("test.thrift", []byte(`
struct Foo {
  // 用户名列表
  2: list<string> strings (
  	custom_tag_1 = "1", custom_tag_2 = "2"
  )
}
`))
					assert.NoError(t, err)
					return ast.(*parser.Document).Structs[0]
				}(),
			},
			want: `struct Foo {
    // 用户名列表
    2: list<string> strings (
        custom_tag_1 = "1", custom_tag_2 = "2"
    )
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MustFormatStruct(tt.args.st))
		})
	}
}
