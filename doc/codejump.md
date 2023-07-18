# codejump

codejump 这里指代一切设计到代码跳转的特性。LSP 中支持的有:

1. Goto Declaration: 在 thrift 中不需要支持
2. Goto Definition:  在 thrift 中可以用来支持 struct/union/enum/exception 异常的跳转
3. Goto Type Definition: 在 thrift 中可以用来支持 struct/union 的跳转。和 Goto Definition 的工作方式稍有不同，
当用户在标识符上使用时，需要推断该标识符的类型，并跳转到类型定义上。
4. Find References: 查找引用。在 thrift 中主要用来做 struct/union/enum/exception 的引用查找

## 实现


