package lsp

import (
	"context"

	"go.lsp.dev/jsonrpc2"
	"k8s.io/klog/v2"
)

func DebugReplier(reply jsonrpc2.Replier) jsonrpc2.Replier {
	return func(ctx context.Context, result interface{}, err error) error {
		klog.InfoS("jsonrpc reply debug", "result", result, "err", err)
		return reply(ctx, result, err)
	}
}

func DebugHandler(handler jsonrpc2.Handler) jsonrpc2.Handler {
	return func(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
		klog.InfoS("jsonrpc request debug", "req", req.Method(), "params", string(req.Params()))
		return handler(ctx, DebugReplier(reply), req)
	}
}
