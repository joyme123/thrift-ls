package symbols

import (
	"context"
	"testing"

	"github.com/joyme123/protocol"
	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/uri"
)

func TestDocumentSymbolsIncludesEnum(t *testing.T) {
	file := uri.URI("file:///tmp/status.thrift")
	ss := cache.BuildSnapshotForTest([]*cache.FileChange{
		{
			URI:     file,
			Version: 0,
			Content: []byte(`enum Status {
  UNKNOWN = 0,
  READY = 1,
}`),
			From: cache.FileChangeTypeDidOpen,
		},
	})

	got := DocumentSymbols(context.TODO(), ss, file)
	if assert.Len(t, got, 1) {
		enum := got[0]
		assert.Equal(t, "Status", enum.Name)
		assert.Equal(t, "Enum", enum.Detail)
		assert.Equal(t, protocol.SymbolKindEnum, enum.Kind)

		if assert.Len(t, enum.Children, 2) {
			assert.Equal(t, "UNKNOWN", enum.Children[0].Name)
			assert.Equal(t, "0", enum.Children[0].Detail)
			assert.Equal(t, protocol.SymbolKindNumber, enum.Children[0].Kind)
			assert.Equal(t, "READY", enum.Children[1].Name)
			assert.Equal(t, "1", enum.Children[1].Detail)
			assert.Equal(t, protocol.SymbolKindNumber, enum.Children[1].Kind)
		}
	}
}
