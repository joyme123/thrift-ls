package lsp

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/joyme123/thrift-ls/lsp/cache"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func (s *Server) initialize(ctx context.Context, params *protocol.InitializeParams) (result *protocol.InitializeResult, err error) {
	rootURI := params.RootURI
	if rootURI == "" {
		rootURI = uri.URI(params.RootPath)
	}

	folders := make([]uri.URI, 0)
	if rootURI != "" {
		folders = append(folders, rootURI)
	}

	for _, ws := range params.WorkspaceFolders {
		folders = append(folders, uri.URI(ws.URI))
	}

	log.Debugln("initialized folders: ", folders)
	if len(folders) > 0 {
		s.session.Initialize(func() {
			for i := range folders {
				s.walkFoldersThriftFile(folders[i])
			}
		})
	}

	return initializeResult(), nil
}

func (s *Server) walkFoldersThriftFile(folder uri.URI) {
	log.Debugln("walk dir 2: ", folder.Filename())
	// WalkDir walk files with lexical order
	filepath.WalkDir(folder.Filename(), func(path string, d fs.DirEntry, err error) error {
		log.Debugln("walk:", path)
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".thrift") {
			return nil
		}

		fileURI := uri.File(path)
		log.Debugln("file path:", fileURI)
		err = s.openFile(context.TODO(), &cache.FileChange{
			URI:     fileURI,
			Version: 0,
			Content: []byte{},
			From:    cache.FileChangeTypeInitialize,
		})
		if err != nil {
			log.Error("openFile err:", err)
			err = nil
		}

		// always return nil to continue parse
		return nil
	})
}

func initializeResult() *protocol.InitializeResult {
	res := &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: true,
				// full is easy to implement. consider to use incremental for performance
				Change:            protocol.TextDocumentSyncKindFull,
				WillSave:          true,
				WillSaveWaitUntil: true,
				Save: &protocol.SaveOptions{
					IncludeText: true,
				},
			},
			CompletionProvider: &protocol.CompletionOptions{
				ResolveProvider: true,
				/**
				 * The additional characters, beyond the defaults provided by the client (typically
				 * [a-zA-Z]), that should automatically trigger a completion request. For example
				 * `.` in JavaScript represents the beginning of an object property or method and is
				 * thus a good candidate for triggering a completion request.
				 *
				 * Most tools trigger a completion request automatically without explicitly
				 * requesting it using a keyboard shortcut (e.g. Ctrl+Space). Typically they
				 * do so when the user starts to type an identifier. For example if the user
				 * types `c` in a JavaScript file code complete will automatically pop up
				 * present `console` besides others as a completion item. Characters that
				 * make up identifiers don't need to be listed here.
				 */
				TriggerCharacters: []string{"."},
			},
			HoverProvider: &protocol.HoverOptions{
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
					WorkDoneProgress: true,
				},
			},
			SignatureHelpProvider: &protocol.SignatureHelpOptions{
				TriggerCharacters:   []string{},
				RetriggerCharacters: []string{},
			},
			DeclarationProvider: &protocol.DeclarationRegistrationOptions{
				DeclarationOptions: protocol.DeclarationOptions{
					WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
						WorkDoneProgress: true,
					},
				},
				TextDocumentRegistrationOptions: protocol.TextDocumentRegistrationOptions{
					DocumentSelector: []*protocol.DocumentFilter{
						{
							Language: "thrift",
						},
					},
				},
				StaticRegistrationOptions: protocol.StaticRegistrationOptions{
					ID: "thriftls",
				},
			},
			DefinitionProvider: &protocol.DefinitionOptions{
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
					WorkDoneProgress: true,
				},
			},
			TypeDefinitionProvider: &protocol.TypeDefinitionRegistrationOptions{
				TextDocumentRegistrationOptions: protocol.TextDocumentRegistrationOptions{
					DocumentSelector: []*protocol.DocumentFilter{
						{
							Language: "thrift",
						},
					},
				},
				TypeDefinitionOptions: protocol.TypeDefinitionOptions{
					WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
						WorkDoneProgress: true,
					},
				},
				StaticRegistrationOptions: protocol.StaticRegistrationOptions{
					ID: "thriftls",
				},
			},
			ReferencesProvider: &protocol.ReferenceOptions{
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
					WorkDoneProgress: true,
				},
			},
			DocumentHighlightProvider: false,
			DocumentSymbolProvider: &protocol.DocumentSymbolOptions{
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
					WorkDoneProgress: true,
				},
				Label: "thriftls",
			},
			CodeActionProvider: &protocol.CodeActionOptions{
				// TODO(jpf): should support code actions
				CodeActionKinds: []protocol.CodeActionKind{},
				ResolveProvider: false,
			},
			CodeLensProvider: &protocol.CodeLensOptions{
				ResolveProvider: false,
			},
			DocumentLinkProvider: &protocol.DocumentLinkOptions{
				ResolveProvider: false,
			},
			ColorProvider: false,
			WorkspaceSymbolProvider: &protocol.WorkspaceSymbolOptions{
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
					WorkDoneProgress: true,
				},
			},
			DocumentFormattingProvider: &protocol.DocumentFormattingOptions{
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
					WorkDoneProgress: true,
				},
			},
			DocumentRangeFormattingProvider: &protocol.DocumentRangeFormattingOptions{
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
					WorkDoneProgress: true,
				},
			},
			DocumentOnTypeFormattingProvider: &protocol.DocumentOnTypeFormattingOptions{
				FirstTriggerCharacter: "}",
				MoreTriggerCharacter:  []string{},
			},
			RenameProvider: &protocol.RenameOptions{
				PrepareProvider: false,
			},
			ExecuteCommandProvider: &protocol.ExecuteCommandOptions{
				Commands: []string{},
			},
			CallHierarchyProvider:      false,
			LinkedEditingRangeProvider: false,
			SemanticTokensProvider: &protocol.SemanticTokensRegistrationOptions{
				TextDocumentRegistrationOptions: protocol.TextDocumentRegistrationOptions{
					DocumentSelector: []*protocol.DocumentFilter{
						{
							Language: "thrift",
						},
					},
				},
				SemanticTokensOptions: protocol.SemanticTokensOptions{
					WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
						WorkDoneProgress: true,
					},
					Legend: protocol.SemanticTokensLegend{
						TokenTypes:     []protocol.SemanticTokenTypes{},
						TokenModifiers: []protocol.SemanticTokenModifiers{},
					},
				},
				StaticRegistrationOptions: protocol.StaticRegistrationOptions{
					ID: "thriftls",
				},
			},
			Workspace: &protocol.ServerCapabilitiesWorkspace{
				WorkspaceFolders: &protocol.ServerCapabilitiesWorkspaceFolders{
					Supported:           true,
					ChangeNotifications: true,
				},
				FileOperations: &protocol.ServerCapabilitiesWorkspaceFileOperations{
					DidCreate: &protocol.FileOperationRegistrationOptions{
						Filters: []protocol.FileOperationFilter{},
					},
					WillCreate: &protocol.FileOperationRegistrationOptions{
						Filters: []protocol.FileOperationFilter{},
					},
					DidRename: &protocol.FileOperationRegistrationOptions{
						Filters: []protocol.FileOperationFilter{},
					},
					WillRename: &protocol.FileOperationRegistrationOptions{
						Filters: []protocol.FileOperationFilter{},
					},
					DidDelete: &protocol.FileOperationRegistrationOptions{
						Filters: []protocol.FileOperationFilter{},
					},
					WillDelete: &protocol.FileOperationRegistrationOptions{
						Filters: []protocol.FileOperationFilter{},
					},
				},
			},
			MonikerProvider: nil,
			Experimental:    nil,
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    ServerName,
			Version: ServerVersion,
		},
	}

	return res
}
