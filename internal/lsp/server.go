package lsp

import (
	"log"

	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "dql-lsp"

var version = "0.0.1"

func RunServer() {
	handler := protocol.Handler{
		Initialize:                      initialize,
		Initialized:                     initialized,
		Shutdown:                        shutdown,
		SetTrace:                        setTrace,
		TextDocumentCompletion:          completion,
		TextDocumentHover:               hover,
		TextDocumentDidOpen:             didOpen,
		TextDocumentDidChange:           didChange,
		TextDocumentDidClose:            didClose,
		TextDocumentDidSave:             didSave,
		WorkspaceDidChangeConfiguration: didChangeConfiguration,
		WorkspaceDidChangeWatchedFiles:  didChangeWatchedFiles,
	}

	server := server.NewServer(&handler, lsName, false)
	if err := server.RunStdio(); err != nil {
		log.Fatal("[Error] Server stopped with error:", err)
	}
}
