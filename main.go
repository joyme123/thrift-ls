package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/joyme123/thrift-ls/log"
	"github.com/joyme123/thrift-ls/lsp"
	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/pkg/fakenet"
	"gopkg.in/yaml.v2"
)

type Options struct {
	LogLevel int `yaml:"logLevel"` // 1: fatal, 2: error, 3: warn, 4: info, 5: debug, 6: trace
}

func main() {
	rand.Seed(time.Now().UnixMilli())

	opts := configInit()
	log.Init(opts.LogLevel)

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

func configInit() *Options {
	logLevel := -1
	flag.IntVar(&logLevel, "logLevel", -1, "set log level")
	flag.Parse()

	dir, err := os.UserHomeDir()
	if err != nil {
		dir = os.TempDir()
	}
	dir = dir + "/.thriftls"
	configFile := dir + "/config.yaml"
	opts := &Options{}

	data, err := os.ReadFile(configFile)
	if err == nil {
		yaml.Unmarshal(data, opts)
	}

	if logLevel >= 0 {
		opts.LogLevel = logLevel // flag can override config file
	}
	if opts.LogLevel == 0 {
		opts.LogLevel = 3
	}

	return opts
}
