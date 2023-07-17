# Diagnostic

诊断是基于 AST 的语义分析，用来发现明显的错误。在 thrift 中，一些可以通过诊断发现这些问题：

1. field index 错误
   a. 超出边界
   b. 重复
2. 标识符冲突。比如定义了两个相同的 struct 等等
3. 循环依赖。因为大多数语言都不允许循环依赖，因此在 thrift 的代码中可以进行循环依赖的检测
4. include 实际不存在
5. 引用无法正常解析

## 实现

### 触发时机

诊断在 LSP 的定义中，是由 Server 主动发送给 Client 的，因此诊断的触发时机由 Server 来决定。
对开发者来说，当其打开 thrift 文件，编辑 thrift 文件后，都希望由 IDE 进行错误检查并提示。因此可以将诊断的触发
时机放在：

1. 用户第一次打开某个 thrift 文件。对这个文件以及相关联的文件进行诊断。
2. 用户编辑了某个 thrift 文件。对这个文件，以及引用此文件的关联文件进行诊断。用户在编辑文本时，会不断的触发 didChange 事件，
为了性能考虑，可以设计一个冷却周期。对于某一个文件，一个诊断后会进入冷却。在此期间内的 didChange 事件不会重复触发诊断。
并且 didChange 事件会重置冷却进度。直到用户停止输入一段时间（比如 2s），诊断会重新开始。

### 诊断队列

诊断是一个耗时的过程，并且单个文件修改会触发多个关联文件的诊断。因此要尽可能的避免无效的诊断触发。可以为每个文件创建一个诊断队列。同一个
队列内，诊断是同步进行的。且新版本 snapshot 的文件诊断可以覆盖旧版本的文件诊断，从而尽可能的减少无效诊断。

