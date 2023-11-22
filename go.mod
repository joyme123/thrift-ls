module github.com/joyme123/thrift-ls

go 1.19

require (
	github.com/cloudwego/thriftgo v0.2.11
	github.com/sergi/go-diff v1.3.1
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.7.0
	go.lsp.dev/jsonrpc2 v0.10.0
	go.lsp.dev/pkg v0.0.0-20210717090340-384b27a52fb2
	go.lsp.dev/protocol v0.12.0
	go.lsp.dev/uri v0.3.0
	go.uber.org/zap v1.21.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/segmentio/asm v1.1.3 // indirect
	github.com/segmentio/encoding v0.3.4 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace go.lsp.dev/protocol => github.com/joyme123/protocol v0.12.1-0.20230807112304-26cf0ace806b
