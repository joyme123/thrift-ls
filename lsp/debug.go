package lsp

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"go.lsp.dev/jsonrpc2"
)

func DebugReplier(reply jsonrpc2.Replier) jsonrpc2.Replier {
	return func(ctx context.Context, result interface{}, err error) error {
		res, _ := json.Marshal(result)
		log.Debug("jsonrpc reply debug: ", "result", string(res), "err", err)
		return reply(ctx, result, err)
	}
}

func DebugHandler(handler jsonrpc2.Handler) jsonrpc2.Handler {
	return func(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
		log.Debug("jsonrpc request debug: ", "req", req.Method(), "params", string(req.Params()))

		defer func() {
			if r := recover(); r != nil {
				log.Errorln("Recovered: ", r)
			}
		}()

		return handler(ctx, DebugReplier(reply), req)
	}
}
