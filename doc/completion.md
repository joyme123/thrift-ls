# completion 实现

completion 要求在用户输入时实时的提供补全的候选项。lsp client 会发送请求到 server，提供用户当前打开的文件以及输入的光标位置。server 在接收到请求后，根据用户输入的位置进行补全项的计算。

补全项的计算主要考虑以下几点：

1. 用户当前输入位置的期望类型。比如是 Header 还是 Definition。如果是 Struct 类型的 Definition，那么是 Identifier 还是 Field。
2. 根据类型进行输入的补全推荐。
- 是 Identifier 的话，可以根据已存在的 Identifier 已经补全推荐。
- 是 Field 中的 FieldType 的话，则提供 thrift 中常见类型的补全推荐。

在补全的过程中，会出现 thrift ast parse 失败的情况。因此需要保留该文件最新可以 parse 的 ast 版本，通过改版本的 ast 进行补全计算.


# fault tolerate parser

https://eyalkalderon.com/blog/nom-error-recovery/
https://github.com/ebkalderon/example-fault-tolerant-parser
https://arxiv.org/pdf/1806.11150.pdf
https://github.com/mna/pigeon
https://github.com/mna/pigeon/wiki

pigeon 实现的 thrift parser: https://github.com/samuel/go-thrift/blob/master/parser/grammar.peg
写 parser 的教程: https://tiarkrompf.github.io/notes/?/just-write-the-parser/
生成好的语法错误信息提示：https://research.swtch.com/yyerror
