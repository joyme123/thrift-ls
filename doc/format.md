# format

format 用来处理 document formatting request，对 thrift 文件进行格式化。

1. 对于没有语法错误的 thrift 代码。可以使用解析过的 AST，通过遍历 AST 的方式输出格式化后的代码。
2. 对于有语法错误的 thrfit 代码。理论上可以在解析时生成 BadNode 时，保存出错的文本。在 AST 生成时原样输出即可。(暂不实现)

实现的难点在于如何从 AST 到 source code。在第一版实现时发现，最难的地方在于如果正确的处理空格，换行，注释和缩进。


```
> from gofmt slide
Lessons learned: Application
Basic source code formatting is great initial goal.
True power lies in source code transformation tools.
Avoid formatting options.
Keep it simple.
Want:

Go parser: source code => syntax tree
Make it easy to manipulate syntax tree in any way possible.
Go printer: syntax tree => source code

Lessons learned: Implementation
Lots of trial and error in initial version.
Single biggest mistake: comments not attached to AST nodes.
=> Current design makes it extremely hard to manipulate AST
and maintain comments in right places.

Cludge: ast.CommentMap
Want:

Easy to manipulate syntax tree with comments attached.
```



## 参考资料

- [the design of a pretty-printing library](https://belle.sourceforge.net/doc/hughes95design.pdf)
- [a prettier printer](https://homepages.inf.ed.ac.uk/wadler/papers/prettier/prettier.pdf)
- [elastic tabstops](https://nick-gravgaard.com/elastic-tabstops/)
- [gofmt slides](https://go.dev/talks/2015/gofmt-en.slide#10)
