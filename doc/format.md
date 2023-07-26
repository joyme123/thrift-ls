# format

format 用来在 document/didSave 时，对 thrift 文件进行格式化。

1. 对于没有语法错误的 thrift 代码。可以使用解析过的 AST，通过遍历 AST 的方式输出格式化后的代码。
2. 对于有语法错误的 thrfit 代码。理论上可以在解析时生成 BadNode 时，保存出错的文本。在 AST 生成时原样输出即可。(暂不实现)
