## 调试

这里以 neovim 为例，在配置文件里用 lua 进行设置：

```lua
require('lspconfig').thriftls.setup{
  handlers=handlers,
  on_attach = on_attach,
  flags = {
    debounce_text_changes = 150,
  },
}

vim.lsp.set_log_level("debug")
```

修改 lspconfig 的代码，增加 thriftls 的配置。这里以 lazy 安装 lspconfig 为例:

path: ~/.local/share/nvim/lazy/nvim-lspconfig/lua/lspconfig/server_configurations/thriftls.lua

```
local util = require 'lspconfig.util'

return {
  default_config = {
    cmd = { 'thriftls' },
    filetypes = { 'thrift' },
    root_dir = function(fname)
      return util.root_pattern('.thrift')(fname)
    end,
    single_file_support = true,
  },
  docs = {
    description = [[
    thrift language server
    ]],
    default_config = {
      root_dir = [[root_pattern(".thrift")]],
    },
  },
}
```

`:LspInfo` 查看 lsp 的信息。一般日志的路径在 ~/.local/state/nvim/lsp.log。可以 tail -f 查看日志的输出进行调试