//go:build windows
// +build windows

package cache

import (
	"syscall"
	"time"
)

// Note: it may be convenient to have this helper return fs.FileInfo, but
// implementing this is actually quite involved on Windows. Since we only
// currently use mtime, keep it simple.
func getFileID(filename string) (FileID, time.Time, error) {
	filename16, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return FileID{}, time.Time{}, err
	}
	h, err := syscall.CreateFile(filename16, 0, 0, nil, syscall.OPEN_EXISTING, uint32(syscall.FILE_FLAG_BACKUP_SEMANTICS), 0)
	if err != nil {
		return FileID{}, time.Time{}, err
	}
	defer syscall.CloseHandle(h)
	var i syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(h, &i); err != nil {
		return FileID{}, time.Time{}, err
	}
	mtime := time.Unix(0, i.LastWriteTime.Nanoseconds())
	return FileID{
		device: uint64(i.VolumeSerialNumber),
		inode:  uint64(i.FileIndexHigh)<<32 | uint64(i.FileIndexLow),
	}, mtime, nil
}
