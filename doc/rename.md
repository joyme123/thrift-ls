# Rename

rename 是基于 AST + 语义分析实现的 identifier 批量替换功能。

以 struct 为例，通过 rename struct name，可以同步批量修改所有引用该 struct name 的地方。

在实现上，Rename 依赖 codejump/find references 的能力。通过查找到 struct name 所有引用的地方，然后对指定位置进行替换即可

在 LSP 中，rename 有两步：

1. prepareRename：client 向 server 发送 prepareRename request，server 根据当前光标位置，计算出 rename 的范围。如果不支持 rename，则返回 error。
2. rename: client 向 server 发送 rename request，这个 request 包含当前的位置以及要变更后的 name，server 根据这两个信息，返回 WorkspaceEdit 响应。
   WorkspaceEdit 响应中包含了当前 workspace 在所有要 rename 的文件和位置.


