package cache

import (
	"context"
	"sync"

	"go.lsp.dev/uri"
)

type Snapshot struct {
	id int64

	view *View

	// ctx is used to cancel background job
	ctx context.Context

	refCount sync.WaitGroup

	files *FilesMap

	parsedCache map[uri.URI]ParsedFile
}

func (s *Snapshot) Acquire() func() {
	s.refCount.Add(1)
	return s.refCount.Done
}

func (s *Snapshot) ReadFile(ctx context.Context, uri uri.URI) (FileHandle, error) {
	s.view.MarkFileKnown(uri)

	if fh, ok := s.files.Get(uri); ok {
		return fh, nil
	}

	fh, err := s.view.fs.ReadFile(ctx, uri)
	if err != nil {
		return nil, err
	}
	s.files.Set(uri, fh)

	return fh, nil
}
