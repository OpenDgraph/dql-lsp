package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	_ "github.com/tliron/commonlog/simple"
)

const lsName = "dql-lsp"

var (
	version string = "0.0.1"
	handler protocol.Handler
)

func init() {
	// Setup file logging
	logFile, err := os.OpenFile("dql-lsp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	log.SetOutput(logFile)
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()

	commonlog.Configure(1, nil) // Default log verbosity

	// Register protocol handlers
	handler = protocol.Handler{
		Initialize:                      initialize,
		Initialized:                     initialized,
		Shutdown:                        shutdown,
		SetTrace:                        setTrace,
		TextDocumentCompletion:          completion, // Register completion
		TextDocumentHover:               hover,      // Register hover
		TextDocumentDidOpen:             didOpen,    // Register file open
		TextDocumentDidChange:           didChange,  // Register file change
		TextDocumentDidClose:            didClose,   // Register file close
		WorkspaceDidChangeConfiguration: didChangeConfiguration,
		WorkspaceDidChangeWatchedFiles:  didChangeWatchedFiles,
	}

	// Start LSP server using stdio
	server := server.NewServer(&handler, lsName, false)
	if err := server.RunStdio(); err != nil {
		log.Println("[Error] Server stopped with error:", err)
	}
	log.Println("==== Starting DQL Language Server ====")
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	log.Println("[Initialize] Params:", params)

	capabilities := handler.CreateServerCapabilities()
	capabilities.TextDocumentSync = protocol.TextDocumentSyncKindFull
	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: []string{" ", "\n"},
	}
	capabilities.HoverProvider = true

	result := protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}

	log.Println("[Initialize] Sending capabilities response")

	return result, nil // <-- VERY IMPORTANT
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	log.Println("[Initialized] Received initialized notification")
	return nil
}

func shutdown(context *glsp.Context) error {
	log.Println("[Shutdown] Client requested shutdown")
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	log.Println("[SetTrace] Trace level set to:", params.Value)
	protocol.SetTraceValue(params.Value)
	return nil
}

func didOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	log.Println("[DidOpen] File opened:", params.TextDocument.URI)
	log.Println("[DidOpen] Content:\n", params.TextDocument.Text)
	return nil
}

func didChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	log.Println("[DidChange] File changed:", params.TextDocument.URI)
	for _, change := range params.ContentChanges {
		if contentChange, ok := change.(protocol.TextDocumentContentChangeEvent); ok {
			log.Println("[DidChange] Change content:\n", contentChange.Text)
			log.Println("[DidChange] Change detected")
		}
	}
	return nil
}

func didClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	log.Println("[DidClose] File closed:", params.TextDocument.URI)
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

func completion(ctx *glsp.Context, params *protocol.CompletionParams) (any, error) {
	log.Println("[Completion] Triggered at:", params.TextDocument.URI, "Position:", params.Position)

	kind := protocol.CompletionItemKindKeyword

	items := []protocol.CompletionItem{
		{
			Label:      "query",
			Kind:       &kind,
			Detail:     ptr("GraphQL operation type"),
			InsertText: ptr("query"),
		},
		{
			Label:      "mutation",
			Kind:       &kind,
			Detail:     ptr("GraphQL operation type"),
			InsertText: ptr("mutation"),
		},
		{
			Label:      "subscription",
			Kind:       &kind,
			Detail:     ptr("GraphQL operation type"),
			InsertText: ptr("subscription"),
		},
	}

	return protocol.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}, nil
}

func hover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	log.Println("[Hover] Hover requested at:", params.TextDocument.URI, "Position:", params.Position)

	word := "query"
	content := fmt.Sprintf("Information about the word: `%s`", word)

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: content,
		},
	}, nil
}

// Helper to create *string from string
func ptr(s string) *string {
	return &s
}
