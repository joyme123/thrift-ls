.PHONY: build
build:
	go build -o bin/thriftls main.go
install:
	cp bin/thriftls /usr/local/bin/thriftls
