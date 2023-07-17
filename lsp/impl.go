package lsp

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/completion"
	"github.com/joyme123/thrift-ls/lsp/types"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func (s *Server) didOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	document := params.TextDocument
	if document.LanguageID != LanguageIDThrift {
		return nil
	}

	fileURI := document.URI
	change := &cache.FileChange{
		URI:     fileURI,
		Version: int(document.Version),
		Content: []byte(document.Text),
		From:    cache.FileChangeTypeDidOpen,
	}

	if err := s.session.UpdateOverlayFS(ctx, []*cache.FileChange{change}); err != nil {
		return err
	}

	if _, err := s.session.ViewOf(fileURI); err != nil {
		// create view for this folder
		filename := fileURI.Filename()
		dir := uri.New(path.Dir(filename))
		s.session.CreateView(dir)
	}

	view, _ := s.session.ViewOf(fileURI)
	view.FileChange(ctx, []*cache.FileChange{change}, func() {
		ss, release := view.Snapshot()
		defer release()
		err := s.diagnostic(ctx, ss, change)
		if err != nil {
			log.Error("diagnostic error", err)
		}
	})

	return nil
}

func (s *Server) didChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	changes := cache.FileChangeFromLSPDidChange(params)
	if err := s.session.UpdateOverlayFS(ctx, changes); err != nil {
		return err
	}

	document := params.TextDocument
	fileURI := document.URI
	view, err := s.session.ViewOf(fileURI)
	if err != nil {
		return err
	}

	view.FileChange(ctx, changes, func() {
		ss, release := view.Snapshot()
		defer release()
		for i := range changes {
			err := s.diagnostic(ctx, ss, changes[i])
			if err != nil {
				log.Error("diagnostic error", err)
			}
		}
	})

	return nil
}

func (s *Server) completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	snapshot, release, fh, err := s.getFileContext(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}
	defer release()

	items, err := completion.DefaultTokenCompletion.Completion(ctx, snapshot, &completion.CompletionRequest{
		TriggerKind: 0,
		Pos: types.Position{
			Line:      params.Position.Line,
			Character: params.Position.Character,
		},
		Fh: fh,
	})
	if err != nil {
		return nil, err
	}

	rng := protocol.Range{
		Start: protocol.Position{
			Line:      params.Position.Line,
			Character: params.Position.Character,
		},
		End: protocol.Position{
			Line:      params.Position.Line,
			Character: params.Position.Character,
		},
	}

	return toLspCompletionList(items, rng), nil
}

func toLspCompletionList(items []*completion.CompletionItem, rng protocol.Range) *protocol.CompletionList {
	list := &protocol.CompletionList{
		IsIncomplete: true,
	}
	for i := range items {
		item := protocol.CompletionItem{
			Label:  items[i].Label,
			Detail: items[i].Detail,
			Kind:   items[i].Kind,
			TextEdit: &protocol.TextEdit{
				NewText: items[i].InsertText,
				Range:   rng,
			},
			FilterText:       strings.TrimLeft(items[i].InsertText, "&*"),
			InsertTextFormat: protocol.InsertTextFormatPlainText,
			SortText:         fmt.Sprintf("%05d", i),
			Preselect:        i == 0,
			Deprecated:       items[i].Deprecated,
			Documentation:    items[i].Documentation,
		}
		list.Items = append(list.Items, item)
	}
	return list
}

func (s *Server) getFileContext(ctx context.Context, uri uri.URI) (ss *cache.Snapshot, release func(), fh cache.FileHandle, err error) {
	var view *cache.View
	view, err = s.session.ViewOf(uri)
	if err != nil {
		return
	}

	ss, release = view.Snapshot()

	fh, err = ss.ReadFile(ctx, uri)
	if err != nil {
		release()
		return
	}

	return
}
