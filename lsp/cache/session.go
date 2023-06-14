package cache

import (
	"fmt"
	"math/rand"
	"sync"

	"go.lsp.dev/uri"
)

type Session struct {
	id int64

	cache *Cache

	viewMu  sync.Mutex
	views   []*View
	viewMap map[uri.URI]*View // map of URI->best view

	*overlayFS
}

func NewSession(cache *Cache) *Session {
	sess := &Session{
		id:        rand.Int63(),
		cache:     cache,
		views:     make([]*View, 0),
		viewMap:   make(map[uri.URI]*View),
		overlayFS: newOverlayFS(cache),
	}

	return sess
}

func (s *Session) CreateView(folder uri.URI) {
	s.views = append(s.views, NewView(folder.Filename(), folder, s.overlayFS))
}

func (s *Session) ViewOf(fileURI uri.URI) (*View, error) {
	s.viewMu.Lock()
	defer s.viewMu.Unlock()

	if view, ok := s.viewMap[fileURI]; ok {
		return view, nil
	}

	if len(s.views) == 0 {
		return nil, fmt.Errorf("views is nil")
	}

	for i := range s.views {
		if s.views[i].ContainsFile(fileURI) {
			s.viewMap[fileURI] = s.views[i]
			return s.views[i], nil
		}
	}

	for i := range s.views {
		if s.views[i].FileKnown(fileURI) {
			s.viewMap[fileURI] = s.views[i]
			return s.views[i], nil
		}
	}

	return s.views[0], nil
}
