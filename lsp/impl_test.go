package lsp

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/memoize"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func Test_DidOpen(t *testing.T) {
	ctx := context.TODO()
	fileURI, err := uri.Parse("file:///opt/file.thrift")
	assert.NoError(t, err)
	fileContent := `
include "base.thrift"

struct Test {
	1: required string Name,
	2: optional i32 Age,
}`
	params := &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        fileURI,
			LanguageID: "thrift",
			Version:    0,
			Text:       fileContent,
		},
	}

	store := &memoize.Store{}
	cache := cache.New(store)
	srv := NewServer(cache, nil)
	err = srv.DidOpen(ctx, params)
	assert.NoError(t, err)

	assert.NotNil(t, srv.session)

	fh, err := srv.session.ReadFile(ctx, fileURI)
	assert.NoError(t, err)
	assert.Equal(t, int(fh.Version()), 0)
	gotContent, err := fh.Content()
	assert.NoError(t, err)
	assert.Equal(t, gotContent, []byte(fileContent))
}

func Test_DidChange(t *testing.T) {
	ctx := context.TODO()
	fileURI, err := uri.Parse("file:///opt/file.thrift")
	assert.NoError(t, err)
	fileContentInit := `
include "base.thrift"

struct Test {
	1: required string Name,
	2: optional i32 Age,
}`
	fileContent := `
include "base.thrift"

struct Test {
	1: required string Name,
	2: optional i32 Age,
	3: required string Email,

}`
	openParams := &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        fileURI,
			LanguageID: "thrift",
			Version:    0,
			Text:       fileContentInit,
		},
	}
	params := &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{
				URI: fileURI,
			},
			Version: 1,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{
				Text: fileContent,
			},
		},
	}

	store := &memoize.Store{}
	cache := cache.New(store)
	srv := NewServer(cache, nil)

	err = srv.DidOpen(ctx, openParams)

	err = srv.DidChange(ctx, params)
	assert.NoError(t, err)

	fh, err := srv.session.ReadFile(ctx, fileURI)
	assert.NoError(t, err)
	assert.Equal(t, int(fh.Version()), 1)
	gotContent, err := fh.Content()
	assert.NoError(t, err)
	assert.Equal(t, gotContent, []byte(fileContent))
}

func Test_Completion(t *testing.T) {
	ctx := context.TODO()
	fileURI, err := uri.Parse("file:///opt/file.thrift")
	assert.NoError(t, err)
	fileContent := `include "base.thrift"

struct Test {
	1: required string Name,
	2: optional i32 Age,
        3: required string N
}`
	openParams := &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        fileURI,
			LanguageID: "thrift",
			Version:    0,
			Text:       fileContent,
		},
	}

	store := &memoize.Store{}
	cache := cache.New(store)
	srv := NewServer(cache, nil)
	err = srv.DidOpen(ctx, openParams)

	completionParams := &protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: fileURI,
			},
			Position: protocol.Position{
				Line:      5, // line and character start with 0
				Character: 28,
			},
		},
		WorkDoneProgressParams: protocol.WorkDoneProgressParams{
			WorkDoneToken: &protocol.ProgressToken{},
		},
		PartialResultParams: protocol.PartialResultParams{
			PartialResultToken: &protocol.ProgressToken{},
		},
		Context: &protocol.CompletionContext{
			TriggerKind: protocol.CompletionTriggerKindInvoked,
		},
	}

	completionList, err := srv.Completion(ctx, completionParams)
	assert.NoError(t, err)

	expectCompletionList := &protocol.CompletionList{
		IsIncomplete: true,
		Items: []protocol.CompletionItem{
			{
				Detail:           "Name",
				Documentation:    "",
				FilterText:       "Name",
				InsertTextFormat: protocol.InsertTextFormatPlainText,
				Kind:             protocol.CompletionItemKindText,
				Label:            "Name",
				Preselect:        true,
				SortText:         fmt.Sprintf("%05d", 0),
				TextEdit: &protocol.TextEdit{
					NewText: "Name",
					Range: protocol.Range{
						// To insert text into a document create a range where start == end.
						Start: protocol.Position{
							Line:      5,
							Character: 28,
						},
						End: protocol.Position{
							Line:      5,
							Character: 28,
						},
					},
				},
			},
		},
	}
	assert.Equal(t, len(expectCompletionList.Items), len(completionList.Items))
	assert.Equal(t, expectCompletionList.Items[0].TextEdit, completionList.Items[0].TextEdit)
	assert.Equal(t, expectCompletionList, completionList)
}

func Test_WalkDir(t *testing.T) {
	filepath.WalkDir("/Users/jiang/projects/bytedance/vke-guosen/pkg/server/idl", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			fmt.Println("dir: ", path)
		} else {
			fmt.Println("file: ", path)
		}
		return nil
	})
}
