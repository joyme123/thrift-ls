gen:
	which pigeon >/dev/null || go install github.com/mna/pigeon@master
	pigeon -o thrift.peg.go ./thrift.peg

test-parser:
	go test -v github.com/joyme123/thrift-ls/parser/test
