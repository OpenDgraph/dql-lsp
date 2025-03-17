package lsp

import (
	"log"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func didOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	log.Println("[DidOpen] File opened:", params.TextDocument.URI)
	OpenDocument(params.TextDocument.URI, params.TextDocument.Text)
	return nil
}

func didChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	log.Println("[DidChange] File changed:", params.TextDocument.URI)

	for _, changeRaw := range params.ContentChanges {
		if change, ok := changeRaw.(protocol.TextDocumentContentChangeEventWhole); ok {
			UpdateDocument(params.TextDocument.URI, change.Text)
		}
	}
	return nil
}

func didSave(ctx *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	log.Println("[DidSave] File saved:", params.TextDocument.URI)
	if params.Text != nil {
		SaveDocument(params.TextDocument.URI, *params.Text)
	}
	return nil
}

func didClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	log.Println("[DidClose] File closed:", params.TextDocument.URI)
	CloseDocument(params.TextDocument.URI)
	return nil
}

func didChangeConfiguration(ctx *glsp.Context, params *protocol.DidChangeConfigurationParams) error {
	log.Println("[Workspace] Configuration changed")
	return nil
}

func didChangeWatchedFiles(ctx *glsp.Context, params *protocol.DidChangeWatchedFilesParams) error {
	log.Println("[Workspace] Watched files changed")
	return nil
}
