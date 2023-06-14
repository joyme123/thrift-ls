package cache

import (
	"context"
	"sync"

	"go.lsp.dev/uri"
)

// An overlayFS is a source.FileSource that keeps track of overlays on top of a
// delegate FileSource.
type overlayFS struct {
	delegate FileSource

	mu       sync.Mutex
	overlays map[uri.URI]*Overlay
}

func newOverlayFS(delegate FileSource) *overlayFS {
	return &overlayFS{
		delegate: delegate,
		overlays: make(map[uri.URI]*Overlay),
	}
}

// Overlays returns a new unordered array of overlays.
func (fs *overlayFS) Overlays() []*Overlay {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	overlays := make([]*Overlay, 0, len(fs.overlays))
	for _, overlay := range fs.overlays {
		overlays = append(overlays, overlay)
	}
	return overlays
}

func (fs *overlayFS) ReadFile(ctx context.Context, uri uri.URI) (FileHandle, error) {
	fs.mu.Lock()
	overlay, ok := fs.overlays[uri]
	fs.mu.Unlock()
	if ok {
		return overlay, nil
	}
	return fs.delegate.ReadFile(ctx, uri)
}

// An Overlay is a file open in the editor. It may have unsaved edits.
// It implements the source.FileHandle interface.
type Overlay struct {
	uri     uri.URI
	content []byte
	hash    Hash
	version int32

	// saved is true if a file matches the state on disk,
	// and therefore does not need to be part of the overlay sent to go/packages.
	saved bool
}

func (o *Overlay) URI() uri.URI { return o.uri }

func (o *Overlay) FileIdentity() FileIdentity {
	return FileIdentity{
		URI:  o.uri,
		Hash: o.hash,
	}
}

func (o *Overlay) Content() ([]byte, error) { return o.content, nil }
func (o *Overlay) Version() int32           { return o.version }
func (o *Overlay) Saved() bool              { return o.saved }
