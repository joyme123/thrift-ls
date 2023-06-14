# completion 实现

completion 要求在用户输入时实时的提供补全的候选项。lsp client 会发送请求到 server，提供用户当前打开的文件以及输入的光标位置。server 在接收到请求后，根据用户输入的位置进行补全项的计算。

补全项的计算主要考虑以下几点：

1. 用户当前输入位置的期望类型。比如是 Header 还是 Definition。如果是 Struct 类型的 Definition，那么是 Identifier 还是 Field。
2. 根据类型进行输入的补全推荐。
- 是 Identifier 的话，可以根据已存在的 Identifier 已经补全推荐。
- 是 Field 中的 FieldType 的话，则提供 thrift 中常见类型的补全推荐。

