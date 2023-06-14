package cache

import (
	"context"
	"math/rand"
	"strings"
	"sync"

	"go.lsp.dev/uri"
)

type View struct {
	id int64

	// name is the user-specified name of this view.
	name string

	// TODO(jpf): view 的设计并不合理
	// workspace folder
	folder uri.URI

	fs FileSource

	knownFilesMu sync.Mutex
	knownFiles   map[uri.URI]bool

	// Track the latest snapshot via the snapshot field, guarded by snapshotMu.
	//
	// Invariant: whenever the snapshot field is overwritten, destroy(snapshot)
	// is called on the previous (overwritten) snapshot while snapshotMu is held,
	// incrementing snapshotWG. During shutdown the final snapshot is
	// overwritten with nil and destroyed, guaranteeing that all observed
	// snapshots have been destroyed via the destroy method, and snapshotWG may
	// be waited upon to let these destroy operations complete.
	snapshotMu sync.Mutex
	snapshot   *Snapshot // latest snapshot; nil after shutdown has been called
}

func NewView(name string, folder uri.URI, fs FileSource) *View {
	view := &View{
		id:         rand.Int63(),
		name:       name,
		folder:     folder,
		fs:         fs,
		knownFiles: make(map[uri.URI]bool),
	}

	view.snapshot = &Snapshot{
		id:       rand.Int63(),
		view:     view,
		ctx:      context.Background(),
		refCount: sync.WaitGroup{},
		files: &FilesMap{
			files:    make(map[uri.URI]FileHandle),
			overlays: make(map[uri.URI]*Overlay),
		},
	}

	return view
}

func (v *View) ContainsFile(uri uri.URI) bool {
	// folder: file:///workdir/
	// file: file:///workdir/file.idl

	folder := v.folder.Filename()
	file := uri.Filename()

	if !strings.HasPrefix(file, folder) {
		return false
	}

	folder = strings.TrimRight(folder, "/")
	file = strings.TrimLeft(file, folder)

	if strings.HasPrefix(file, "/") {
		return true
	}

	return false
}

func (v *View) MarkFileKnown(fileURI uri.URI) {
	v.knownFilesMu.Lock()
	defer v.knownFilesMu.Unlock()

	if v.knownFiles == nil {
		v.knownFiles = make(map[uri.URI]bool)
	}

	v.knownFiles[fileURI] = true
}

func (v *View) FileKnown(uri uri.URI) bool {
	v.knownFilesMu.Lock()
	defer v.knownFilesMu.Unlock()

	return v.knownFiles[uri]
}

func (v *View) FileChange(ctx context.Context, changes []*FileChange) (*Snapshot, func()) {

	return nil, nil
}

func (v *View) Snapshot() (*Snapshot, func()) {
	v.snapshotMu.Lock()
	defer v.snapshotMu.Unlock()

	if v.snapshot == nil {

	}

	return v.snapshot, v.snapshot.Acquire()
}
