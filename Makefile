# 获取系统架构
UNAME_M := $(shell uname -m)

# 根据系统架构设置 ARCH 变量
ifeq ($(UNAME_M),x86_64)
	ARCH := amd64
else ifeq ($(UNAME_M),amd64)
	ARCH := amd64
else ifeq ($(UNAME_M),aarch64)
	ARCH := arm64
else ifeq ($(UNAME_M),arm64)
	ARCH := arm64
else
	ARCH := unknown
endif

.PHONY: build install install-for-mason test
build:
	go build -o bin/thriftls main.go
install:
	cp bin/thriftls /usr/local/bin/thriftls
install-for-mason:
	cp bin/thriftls ~/.local/share/nvim/mason/packages/thriftls/thriftls-darwin-$(ARCH)
test:
	@go test -gcflags=all=-l -gcflags=all=-d=checkptr=0 -race -coverpkg=./... -coverprofile=coverage.out $(shell go list ./...)
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'

e2e-test:
	@bash tests/e2e/run-e2e.sh
