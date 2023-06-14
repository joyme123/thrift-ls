package lsp

import (
	"context"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/memoize"
	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/pkg/event"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

type StreamServer struct {
	logger *zap.Logger

	cache *cache.Cache
}

func NewStreamServer() *StreamServer {
	logger, _ := zap.NewProduction()

	store := &memoize.Store{}

	return &StreamServer{
		cache:  cache.New(store),
		logger: logger,
	}
}

func (s *StreamServer) ServeStream(ctx context.Context, conn jsonrpc2.Conn) error {
	client := protocol.ClientDispatcher(conn, s.logger)

	server := NewServer(s.cache, client)
	// Clients may or may not send a shutdown message. Make sure the server is
	// shut down.
	// TODO(rFindley): this shutdown should perhaps be on a disconnected context.
	defer func() {
		if err := server.Shutdown(ctx); err != nil {
			event.Error(ctx, "error shutting down", err)
		}
	}()
	ctx = protocol.WithClient(ctx, client)
	conn.Go(ctx,
		DebugHandler(
			protocol.Handlers(
				protocol.ServerHandler(server, jsonrpc2.MethodNotFoundHandler))))
	<-conn.Done()
	return conn.Err()
}
