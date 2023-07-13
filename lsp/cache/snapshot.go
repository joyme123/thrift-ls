package cache

import (
	"context"
	"math/rand"
	"sync"

	"github.com/joyme123/thrift-ls/lsp/memoize"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/uri"
)

type Snapshot struct {
	id int64

	view *View

	// ctx is used to cancel background job
	ctx context.Context

	refCount sync.WaitGroup

	files *FilesMap

	store *memoize.Store

	parsedCache *ParseCaches
}

func (s *Snapshot) Acquire() func() {
	s.refCount.Add(1)
	return s.refCount.Done
}

func (s *Snapshot) Initialize(ctx context.Context) {

}

func (s *Snapshot) ReadFile(ctx context.Context, uri uri.URI) (FileHandle, error) {
	log.Debugln("snapshot read file", uri)
	s.view.MarkFileKnown(uri)

	if fh, ok := s.files.Get(uri); ok {
		return fh, nil
	}

	log.Debugln("snapshot read from fs")
	fh, err := s.view.fs.ReadFile(ctx, uri)
	if err != nil {
		return nil, err
	}
	s.files.Set(uri, fh)

	return fh, nil
}

func (s *Snapshot) Parse(ctx context.Context, uri uri.URI) error {
	fh, err := s.ReadFile(ctx, uri)
	if err != nil {
		return err
	}

	// DEBUG
	content, _ := fh.Content()
	log.Debugln("parse content:", string(content))

	pf, err := Parse(fh)
	if err != nil {
		return err
	}

	s.parsedCache.Set(uri, pf)

	return nil
}

func (s *Snapshot) Tokens() map[string]struct{} {
	return s.parsedCache.Tokens()
}

func (s *Snapshot) GetParsedFile(uri uri.URI) *ParsedFile {
	return s.parsedCache.Get(uri)
}

func (s *Snapshot) clone() (*Snapshot, func()) {
	snap := &Snapshot{
		id:   rand.Int63(),
		view: s.view,
		ctx:  context.Background(),
		// TODO(jpf): file change 没有更新，导致读到旧的缓存
		// files:       s.files.Clone(),
		files: &FilesMap{
			files:    make(map[uri.URI]FileHandle),
			overlays: make(map[uri.URI]*Overlay),
		},
		parsedCache: s.parsedCache.Clone(),
	}

	return snap, snap.Acquire()
}
