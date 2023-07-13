##

名词说明:

- cache: 全局缓存，可以在同一个 server 下不同的 session 中共享
- session: 会话。指一个  connection 接入的会话。一个 session 中会有多个 view
- view: 视图。一个 workspace 对应一个 view。
- workspace: 指一个工程
- snapshot: 指 view 的一个快照。保存了当前 view 中的状态，比如：包含的文件，目录等等

## 文件内容维护

当编辑器打开代码文件后，language server 就不能再从磁盘中读取文件了。此时编辑器会通过以下几个事件来通知 server:

1. textDocument/didOpen: 整个文件的文本内容会通过该事件一起传到 server。
2. textDocument/didChange: 编辑器通知 server 代码发生了变更。变更通知有三种方式：
   - 不更新
   - 全量更新
   - 增量更新
3. textDocument/didClose: 编辑器关闭了该文件。后续如果有需要该文件的地方，server 可以自己去磁盘中读取。

对于通过 include 引入的文件，server 需要自己去磁盘中读取该文件的内容进行解析。

对于未在编辑器中打开的文件，server 可以通过 fsnotify 进行文件内容变更检测，从而可以实时维护最新的文件内容。

## Parse 和 Analysis

Parse: 解析代码文件生成 AST，这个过程中可以发现语法上的错误

Analysis:  根据已有的 AST 进行分析，主要用来找出一些非语法错误。 比如：

1. Struct 在 field index 超出可用范围，field index 冲突
2. 类型不存在
3. 目标语言的关键字冲突。比如要根据 thrift 生成 go 语言，就需要在 thrift 代码里避免 go 语言的关键字。等等

Parse 和 Analysis 在以下情况下被触发:

1. 在 textDocument/didOpen 时，会进行 Parse 和 Analysis 两个步骤。 如果发现已经进行过 Parse，则会跳过 Parse 和 Analysis。
2. 在 textDocument/didChange 时，会仅针对当前文件进行 Parse，并进行全局的 Analysis。
3. 在检测到磁盘上的文件发生变更时(比如 git 切换分支，丢弃变更等操作)，会对变更的文件进行 Parse，以及全局的 Analysis。

## snapshot 

snapshot 是属于 view 的，在 view 的初始化时，会对 snapshot 进行初始化。为了防止 snapshot 被释放，view 需要自己持有一份 snapshot 的引用。当 view 不需要
该 snapshot 时主动释放即可。

snapshot 是当前 view 的快照，保存了当前 view 的所有代码的 parse 和 analysis 结果。在发生 `open` `change` 或磁盘上的文件变更时，需要根据最新的文件内容进行 Parse 和 Analysis。
因此这里可以参考上面的 `Parse 和 Analysis` 章节所说的，新的 snapshot 并不是完全重新 Parse 和 Analysis 所有的代码，而是继承之前的 Parse 和 Analysis 结果，尽量小的计算出结果。
具体实现上，snapshot 对自己进行 clone，生成新的 snapshot 实例，然后再处理 changes。

## Promise

promise 提供了一种缓存机制，防止对同一个方法进行重复的调用，并可以异步的获取结果。例如：

1. 对 snapshot 中的某个文件进行 parse 时，该文件可能已经被其他事件触发，正在 parse 中。那么后面触发的这个 parse 就不需要重复进行，只需要等待之前的 parse 结果并获取即可。

## Mapper

1. 在 LSP 定义中，Line 和 Column 是从 0 开始的。在 Parser 中, Line 和 Column 是从 1 开始的。
2. 在 LSP 定义中，字符编码可以支持 utf8/utf16/utf32。默认支持 utf16。因此 language server 的实现中默认支持 utf16 即可。但是对于 parser 来说，是基于 rune 的。
因此对于 Position 的计算使用的编码是不一样的。

以下面的字符为例: a 的 position 是(0, 2)，在 parser 里的 position 是(1,2)
```
😀a
```


因此需要 有一个 Mapper 负责：根据不同的 Position 和编码定义，在两边实现灵活的转换。

