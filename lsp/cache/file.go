package cache

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"syscall"
	"time"

	"go.lsp.dev/uri"
)

// A FileID uniquely identifies a file in the file system.
//
// If GetFileID(name1) returns the same ID as GetFileID(name2), the two file
// names denote the same file.
// A FileID is comparable, and thus suitable for use as a map key.
type FileID struct {
	device, inode uint64
}

// GetFileID returns the file system's identifier for the file, and its
// modification time.
// Like os.Stat, it reads through symbolic links.
func GetFileID(filename string) (FileID, time.Time, error) { return getFileID(filename) }

func getFileID(filename string) (FileID, time.Time, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return FileID{}, time.Time{}, err
	}
	stat := fi.Sys().(*syscall.Stat_t)
	return FileID{
		device: uint64(stat.Dev), // (int32 on darwin, uint64 on linux)
		inode:  stat.Ino,
	}, fi.ModTime(), nil
}

type Hash [sha256.Size]byte

// HashOf returns the hash of some data.
func HashOf(data []byte) Hash {
	return Hash(sha256.Sum256(data))
}

// Hashf returns the hash of a printf-formatted string.
func Hashf(format string, args ...interface{}) Hash {
	// Although this looks alloc-heavy, it is faster than using
	// Fprintf on sha256.New() because the allocations don't escape.
	return HashOf([]byte(fmt.Sprintf(format, args...)))
}

// String returns the digest as a string of hex digits.
func (h Hash) String() string {
	return fmt.Sprintf("%64x", [sha256.Size]byte(h))
}

// Less returns true if the given hash is less than the other.
func (h Hash) Less(other Hash) bool {
	return bytes.Compare(h[:], other[:]) < 0
}

// XORWith updates *h to *h XOR h2.
func (h *Hash) XORWith(h2 Hash) {
	// Small enough that we don't need crypto/subtle.XORBytes.
	for i := range h {
		h[i] ^= h2[i]
	}
}

// FileIdentity uniquely identifies a file at a version from a FileSystem.
type FileIdentity struct {
	URI  uri.URI
	Hash Hash // digest of file contents
}

func (id FileIdentity) String() string {
	return fmt.Sprintf("%s%s", id.URI, id.Hash)
}

// A FileHandle represents the URI, content, hash, and optional
// version of a file tracked by the LSP session.
//
// File content may be provided by the file system (for Saved files)
// or from an overlay, for open files with unsaved edits.
// A FileHandle may record an attempt to read a non-existent file,
// in which case Content returns an error.
type FileHandle interface {
	// URI is the URI for this file handle.
	// TODO(rfindley): this is not actually well-defined. In some cases, there
	// may be more than one URI that resolve to the same FileHandle. Which one is
	// this?
	URI() uri.URI
	// FileIdentity returns a FileIdentity for the file, even if there was an
	// error reading it.
	FileIdentity() FileIdentity
	// Saved reports whether the file has the same content on disk:
	// it is false for files open on an editor with unsaved edits.
	Saved() bool
	// Version returns the file version, as defined by the LSP client.
	// For on-disk file handles, Version returns 0.
	Version() int32
	// Content returns the contents of a file.
	// If the file is not available, returns a nil slice and an error.
	Content() ([]byte, error)
}

// A FileSource maps URIs to FileHandles.
type FileSource interface {
	// ReadFile returns the FileHandle for a given URI, either by
	// reading the content of the file or by obtaining it from a cache.
	ReadFile(ctx context.Context, uri uri.URI) (FileHandle, error)
}

// FilesMap holds files on disk and overlay files
type FilesMap struct {
	files    map[uri.URI]FileHandle
	overlays map[uri.URI]*Overlay
}

func (m *FilesMap) Get(key uri.URI) (FileHandle, bool) {
	fh, ok := m.files[key]
	return fh, ok
}

func (m *FilesMap) Set(key uri.URI, file FileHandle) {
	m.files[key] = file
	if o, ok := file.(*Overlay); ok {
		m.overlays[key] = o
	}
}

func (m *FilesMap) Clone() *FilesMap {
	newMap := &FilesMap{
		files:    make(map[uri.URI]FileHandle),
		overlays: make(map[uri.URI]*Overlay),
	}
	for key := range m.files {
		newMap.files[key] = m.files[key]
	}
	for key := range m.overlays {
		newMap.overlays[key] = m.overlays[key]
	}

	return newMap
}

func (m *FilesMap) Destroy() {
	m.files = nil
	m.overlays = nil
}

type FileChangeType string

const (
	FileChangeTypeDidOpen   FileChangeType = "DidOpen"
	FileChangeTypeDidChange FileChangeType = "DidChange"
	FileChangeTypeDidSave   FileChangeType = "DidSave"
)

type FileChange struct {
	URI     uri.URI
	Version int
	Content []byte
	From    FileChangeType
}
