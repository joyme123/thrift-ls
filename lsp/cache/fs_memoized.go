package cache

import (
	"context"
	"os"
	"sync"
	"time"

	"go.lsp.dev/uri"
)

// A memoizedFS is a file source that memoizes reads, to reduce IO.
type memoizedFS struct {
	mu sync.Mutex

	// filesByID maps existing file inodes to the result of a read.
	// (The read may have failed, e.g. due to EACCES or a delete between stat+read.)
	// Each slice is a non-empty list of aliases: different URIs.
	filesByID map[FileID][]*DiskFile
}

func newMemoizedFS() *memoizedFS {
	return &memoizedFS{filesByID: make(map[FileID][]*DiskFile)}
}

// A DiskFile is a file on the filesystem, or a failure to read one.
// It implements the source.FileHandle interface.
type DiskFile struct {
	uri     uri.URI
	modTime time.Time
	content []byte
	hash    Hash
	err     error
}

func (h *DiskFile) URI() uri.URI { return h.uri }

func (h *DiskFile) FileIdentity() FileIdentity {
	return FileIdentity{
		URI:  h.uri,
		Hash: h.hash,
	}
}

func (h *DiskFile) Saved() bool              { return true }
func (h *DiskFile) Version() int32           { return 0 }
func (h *DiskFile) Content() ([]byte, error) { return h.content, h.err }

// ReadFile stats and (maybe) reads the file, updates the cache, and returns it.
func (fs *memoizedFS) ReadFile(ctx context.Context, uri uri.URI) (FileHandle, error) {
	id, mtime, err := GetFileID(uri.Filename())
	if err != nil {
		// file does not exist
		return &DiskFile{
			err: err,
			uri: uri,
		}, nil
	}

	// We check if the file has changed by comparing modification times. Notably,
	// this is an imperfect heuristic as various systems have low resolution
	// mtimes (as much as 1s on WSL or s390x builders), so we only cache
	// filehandles if mtime is old enough to be reliable, meaning that we don't
	// expect a subsequent write to have the same mtime.
	//
	// The coarsest mtime precision we've seen in practice is 1s, so consider
	// mtime to be unreliable if it is less than 2s old. Capture this before
	// doing anything else.
	recentlyModified := time.Since(mtime) < 2*time.Second

	fs.mu.Lock()
	fhs, ok := fs.filesByID[id]
	if ok && fhs[0].modTime.Equal(mtime) {
		var fh *DiskFile
		// We have already seen this file and it has not changed.
		for _, h := range fhs {
			if h.uri == uri {
				fh = h
				break
			}
		}
		// No file handle for this exact URI. Create an alias, but share content.
		if fh == nil {
			newFH := *fhs[0]
			newFH.uri = uri
			fh = &newFH
			fhs = append(fhs, fh)
			fs.filesByID[id] = fhs
		}
		fs.mu.Unlock()
		return fh, nil
	}
	fs.mu.Unlock()

	// Unknown file, or file has changed. Read (or re-read) it.
	fh, err := readFile(ctx, uri, mtime) // ~25us
	if err != nil {
		return nil, err // e.g. cancelled (not: read failed)
	}

	fs.mu.Lock()
	if !recentlyModified {
		fs.filesByID[id] = []*DiskFile{fh}
	} else {
		delete(fs.filesByID, id)
	}
	fs.mu.Unlock()
	return fh, nil
}

// ioLimit limits the number of parallel file reads per process.
var ioLimit = make(chan struct{}, 128)

func readFile(ctx context.Context, uri uri.URI, mtime time.Time) (*DiskFile, error) {
	select {
	case ioLimit <- struct{}{}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	defer func() { <-ioLimit }()

	// It is possible that a race causes us to read a file with different file
	// ID, or whose mtime differs from the given mtime. However, in these cases
	// we expect the client to notify of a subsequent file change, and the file
	// content should be eventually consistent.
	content, err := os.ReadFile(uri.Filename()) // ~20us
	if err != nil {
		content = nil // just in case
	}
	return &DiskFile{
		modTime: mtime,
		uri:     uri,
		content: content,
		hash:    HashOf(content),
		err:     err,
	}, nil
}
