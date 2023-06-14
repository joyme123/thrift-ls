package cache

import (
	"reflect"
	"strconv"
	"sync/atomic"

	"github.com/joyme123/thrift-ls/lsp/memoize"
)

type Cache struct {
	id string

	store *memoize.Store

	*memoizedFS
}

var cacheIndex int64

func New(store *memoize.Store) *Cache {
	index := atomic.AddInt64(&cacheIndex, 1)

	if store == nil {
		store = &memoize.Store{}
	}

	c := &Cache{
		id:         strconv.FormatInt(index, 10),
		store:      store,
		memoizedFS: &memoizedFS{filesByID: map[FileID][]*DiskFile{}},
	}
	return c
}

func (c *Cache) ID() string                     { return c.id }
func (c *Cache) MemStats() map[reflect.Type]int { return c.store.Stats() }
