.PHONY: build install test
build:
	go build -o bin/thriftls main.go
install:
	cp bin/thriftls /usr/local/bin/thriftls
test:
	@go test -gcflags=all=-l -gcflags=all=-d=checkptr=0 -race -coverpkg=./... -coverprofile=coverage.out $(shell go list ./...)
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'
