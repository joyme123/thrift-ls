package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/joyme123/thrift-ls/format"
	tlog "github.com/joyme123/thrift-ls/log"
	"github.com/joyme123/thrift-ls/lsp"
	"github.com/joyme123/thrift-ls/parser"
	"github.com/joyme123/thrift-ls/utils/diff"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/pkg/fakenet"
	"gopkg.in/yaml.v2"
)

type Options struct {
	LogLevel int `yaml:"logLevel"` // 1: fatal, 2: error, 3: warn, 4: info, 5: debug, 6: trace
}

func main_format(opt format.Options, file string) error {
	if file == "" {
		err := errors.New("must specified a thrift file to format")
		fmt.Println(err)
		return err
	}

	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return err
	}

	thrift_file := filepath.Base(file)
	ast, err := parser.Parse(thrift_file, content)
	if err != nil {
		fmt.Println(err)
		return err
	}
	formated, err := format.FormatDocumentWithValidation(ast.(*parser.Document), true)

	if opt.Write {
		var perms os.FileMode
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Println(err)
			return err
		}
		perms = fileInfo.Mode() // 使用原文件的权限

		// overwrite
		err = os.WriteFile(file, []byte(formated), perms)
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		if opt.Diff {
			diffLines := diff.Diff("old", content, "new", []byte(formated))
			fmt.Print(string(diffLines))
		} else {
			fmt.Print(formated)
		}
		return err
	}

	return nil

}

func main() {
	rand.Seed(time.Now().UnixMilli())

	formatter := false
	formatFile := ""
	flag.BoolVar(&formatter, "format", false, "use thrift-ls as a format tool")
	flag.StringVar(&formatFile, "f", "", "file path to format")
	formatOpts := format.Options{}
	formatOpts.SetFlags()
	flag.Parse()
	formatOpts.InitDefault()

	opts := configInit()
	tlog.Init(opts.LogLevel)

	if formatter {
		main_format(formatOpts, formatFile)
		os.Exit(1)
		return
	}

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
	opts := &Options{}

	logLevel := -1
	flag.IntVar(&logLevel, "logLevel", -1, "set log level")
	flag.Parse()

	dir, err := os.UserHomeDir()
	if err != nil {
		dir = os.TempDir()
	}
	dir = dir + "/.thriftls"
	configFile := dir + "/config.yaml"

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
