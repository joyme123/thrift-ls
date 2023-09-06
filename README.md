# Thrift language server

[![Language grade: Go](https://img.shields.io/lgtm/grade/go/g/joyme123/thrift-ls.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/joyme123/thrift-ls/context:go)

![Go](https://github.com/joyme123/thrift-ls/workflows/Go/badge.svg?branch=main)

thrift-ls implements language server protocol

## Usages

### vim

use thriftls as a lsp provider for thrift

### neovim

You can use [mason](https://github.com/williamboman/mason.nvim) to install thriftls.
And use [nvim-lspconfig](https://github.com/neovim/nvim-lspconfig) to configure thriftls

`:LspInfo` to set lsp information. default log file location: `~/.local/state/nvim/lsp.log`.

![neovim](./doc/image/neovim.png)

### vscode

install thrift-language-server in extension market

![vscode](./doc/image/vscode.png)

## Configurations

config file default location:

- windows: `C:\Users\${user}\.thriftls\config.yaml`
- macos, linux: `~/.thriftls/config.yaml`

## ScreenShot
