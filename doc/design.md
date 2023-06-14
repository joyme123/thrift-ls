##

名词说明:

- cache: 全局缓存，可以在同一个 server 下不同的 session 中共享
- session: 会话。指一个  connection 接入的会话。一个 session 中会有多个 view
- view: 视图。一个 workspace 对应一个 view。
- workspace: 指一个工程
- snapshot: 指 view 的一个快照。保存了当前 view 中的状态，比如：包含的文件，目录等等
