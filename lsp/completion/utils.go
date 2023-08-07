package completion

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/joyme123/thrift-ls/lsp/constants"
	log "github.com/sirupsen/logrus"
)

func ListDirAndFiles(dir, prefix string) (res []string, err error) {
	// handle prefix list ../../us
	if len(prefix) > 0 {
		prefix = filepath.Clean(prefix)
	}
	up := strings.Count(prefix, "../")

	pathItems := strings.Split(dir, "/")
	if len(pathItems) < up {
		return
	}

	pathItems = pathItems[0 : len(pathItems)-up]

	dir, filePrefix := filepath.Split(strings.TrimLeft(prefix, "../"))
	baseDir := strings.Join(pathItems, "/") + "/" + dir

	log.Debugf("include completion: walk dir %s with prefix %s", baseDir, filePrefix)

	filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || baseDir == path {
			return nil
		}
		log.Debugf("include completion: name: %s, prefix: %s", d.Name(), filePrefix)
		if strings.HasPrefix(d.Name(), filePrefix) {
			if d.IsDir() {
				res = append(res, d.Name()+"/")
			} else if strings.HasSuffix(d.Name(), constants.ThriftExtension) {
				res = append(res, d.Name())
			}
		}

		if d.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})

	return
}
