.PHONY: build install install-for-mason test
build:
	go build -o bin/thriftls main.go
install:
	cp bin/thriftls /usr/local/bin/thriftls
install-for-mason:
	cp bin/thriftls ~/.local/share/nvim/mason/packages/thriftls/thriftls-darwin-arm64
test:
	@go test -gcflags=all=-l -gcflags=all=-d=checkptr=0 -race -coverpkg=./... -coverprofile=coverage.out $(shell go list ./...)
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'

e2e-test:
	@bash tests/e2e/run-e2e.sh
