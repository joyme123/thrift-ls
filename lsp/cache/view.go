package cache

import (
	"context"
	"encoding/json"
	"math/rand"
	"strings"
	"sync"

	"github.com/joyme123/thrift-ls/lsp/memoize"
	log "github.com/sirupsen/logrus"
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
	snapshotMu      sync.Mutex
	snapshot        *Snapshot // latest snapshot; nil after shutdown has been called
	snapshotRelease func()
}

func NewView(name string, folder uri.URI, fs FileSource, store *memoize.Store) *View {
	view := &View{
		id:         rand.Int63(),
		name:       name,
		folder:     folder,
		fs:         fs,
		knownFiles: make(map[uri.URI]bool),
	}

	view.snapshot = NewSnapshot(view, store)

	view.snapshotRelease = view.snapshot.Acquire()

	asyncRelease := view.snapshot.Acquire()
	go func() {
		defer asyncRelease()

		view.snapshotMu.Lock()
		view.snapshot.Initialize(context.Background())
		view.snapshotMu.Unlock()
		return
	}()

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

func (v *View) FileChange(ctx context.Context, changes []*FileChange, postFns ...func()) {
	for _, change := range changes {
		v.MarkFileKnown(change.URI)
	}

	// snapshot clone
	newSnapshot, release := v.snapshot.clone()
	// release previous snapshot
	v.snapshotRelease()
	v.snapshotMu.Lock()
	v.snapshot = newSnapshot
	for _, change := range changes {
		v.snapshot.ForgetFile(change.URI)
	}
	v.snapshotMu.Unlock()
	v.snapshotRelease = release

	asyncRelease := v.snapshot.Acquire()
	// handle current snapshot

	go func() {
		defer asyncRelease()
		uris := make(map[uri.URI]struct{})
		for _, change := range changes {
			uris[change.URI] = struct{}{}
		}
		for uri := range uris {
			pf, err := v.snapshot.Parse(ctx, uri)
			if err != nil {
				log.Errorf("parse error: %v", err)
			} else {
				ast, _ := json.MarshalIndent(pf.ast, "", "  ")
				log.Debugln("parsed ast: ", string(ast))
			}
		}

		for i := range postFns {
			postFns[i]()
		}
	}()
	return
}

func (v *View) Snapshot() (*Snapshot, func()) {
	v.snapshotMu.Lock()
	defer v.snapshotMu.Unlock()

	if v.snapshot == nil {

	}

	return v.snapshot, v.snapshot.Acquire()
}
