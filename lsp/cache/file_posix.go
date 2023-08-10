//go:build !windows && !plan9
// +build !windows,!plan9

package cache

import (
	"os"
	"syscall"
	"time"
)

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
