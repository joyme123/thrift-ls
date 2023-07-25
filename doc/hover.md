# hover

hover 用来展示代码的一些属性。比如某个 identifier 的类型，某个引用常量的值等等。lsp 中还有一个 signatureHelp 和 hover 类似，
不过在 thrift 中不存在函数的调用，因此这里只考虑 hover 的实现。

以下类型的 ast node 上可以使用 hover:

1. FieldType：function type/struct field/union field/exception field/function argument/function throws field/typedef/const
2. ConstValue: field default value

hover 的实现:
   1. 找到 definition： 复用 goto definition 的查找实现
   2. 展示 definition: 复用 format 的实现
