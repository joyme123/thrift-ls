# completion 实现

## 概述

completion 要求在用户输入时实时的提供补全的候选项。lsp client 会发送请求到 server，提供用户当前打开的文件以及输入的光标位置。server 在接收到请求后，根据用户输入的位置进行补全项的计算。

补全项的计算主要考虑以下几点：

1. 用户当前输入位置的期望类型。比如是 Header 还是 Definition。如果是 Struct 类型的 Definition，那么是 Identifier 还是 Field。
2. 根据类型进行输入的补全推荐。
- 是 Identifier 的话，可以根据已存在的 Identifier 已经补全推荐。
- 是 Field 中的 FieldType 的话，则提供 thrift 中常见类型的补全推荐。

在补全的过程中，会出现 thrift ast parse 失败的情况。因此需要保留该文件最新可以 parse 的 ast 版本，通过该版本的 ast 进行补全计算.

## 交互细节设计

> 以 VIM 为例

1. vim 发送 CompletionRequest，参数有
   a. TriggerKind
   b. Pos: 行号Line, 当前行的偏移量Character 
   c: URI: 打开文件的 URI
   d: TriggerCharacter:  触发补全的字符。这个仅当 TriggerKind 是 TriggerCharacter 时有值

2. language server 收到请求后，开始计算 completion.
   1. 根据 URI 定位到文件，根据 Pos 计算得到 parser 的 Location. 
   2. 再根据 AST 中的位置信息，找到 Location 对应的 AST node 以及路径。这里的路径指的是从 root node 到当前的 AST node 所有的 node。
   3. 根据 AST node 和 path 信息，进行补全推荐。

基于语义的实现方式，对 parser 的要求较高，parser 需要在语法错误的情况下知道当前的输入位置的 node 属性。因此目前使用的是基于 token 的补全方案。

### 补全清单

1. Document -> BadHeader: 对输入进行 include snippet 补全
2. Document -> Header -> Literal: 对 Path 进行补全
3. Document -> BadDefinition: 对输入进行 struct/union/exception/servcie/Enum/Const/Typedef snippet 的补全
4. Document -> Definition -> Struct/Union/Exception -> Bad Field: 对 required/optional 以及基础/自定义类型进行补全
5. Document -> Definition -> Struct/Union/Exception -> Field -> Identifier: 根据已有的 Identifier 进行补全
6. Document -> Definition -> Service -> 


## fault tolerate parser

1. https://eyalkalderon.com/blog/nom-error-recovery/
2. https://github.com/ebkalderon/example-fault-tolerant-parser
3. https://arxiv.org/pdf/1806.11150.pdf
4. https://github.com/mna/pigeon
5. https://github.com/mna/pigeon/wiki
6. pigeon 实现的 thrift parser: https://github.com/samuel/go-thrift/blob/master/parser/grammar.peg
7. 写 parser 的教程: https://tiarkrompf.github.io/notes/?/just-write-the-parser/
8. 生成好的语法错误信息提示：https://research.swtch.com/yyerror
