package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/joyme123/thrift-ls/lsp"
	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/pkg/fakenet"
	"k8s.io/klog/v2"
)

type Options struct{}

func main() {
	rand.Seed(time.Now().UnixMilli())

	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("log_file", "/tmp/thriftls.log")
	flag.Parse()
	defer klog.Flush()

	ctx := context.Background()
	// server := &lsp.Server{}
	// handler := protocol.ServerHandler(server, nil)
	//
	// streamServer := jsonrpc2.HandlerServer(handler)
	// if err := jsonrpc2.ListenAndServe(ctx, "tcp", "127.0.0.1:8000", streamServer, 60*time.Second); err != nil {
	// 	panic(err)
	// }

	ss := lsp.NewStreamServer()
	stream := jsonrpc2.NewStream(fakenet.NewConn("stdio", os.Stdin, os.Stdout))
	conn := jsonrpc2.NewConn(stream)
	err := ss.ServeStream(ctx, conn)
	if errors.Is(err, io.EOF) {
		return
	}
	panic(err)
}
