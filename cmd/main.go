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

func main() {
	logFile, err := os.OpenFile("dql-lsp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Erro ao abrir log de arquivo:", err)
	}
	log.SetOutput(logFile)
	log.Println("==== Iniciando LSP ====")

	commonlog.Configure(1, nil) // Verbosidade de log

	handler = protocol.Handler{
		Initialize:             initialize,
		Initialized:            initialized,
		Shutdown:               shutdown,
		SetTrace:               setTrace,
		TextDocumentCompletion: completion, // Completion registrado
		TextDocumentHover:      hover,      // Hover registrado
		TextDocumentDidOpen:    didOpen,    // Abrir arquivo registrado
		TextDocumentDidChange:  didChange,
		TextDocumentDidClose:   didClose,
	}

	server := server.NewServer(&handler, lsName, false)
	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()
	capabilities.TextDocumentSync = protocol.TextDocumentSyncKindFull
	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: []string{" ", "\n"},
	}
	capabilities.HoverProvider = true

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func didOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	log.Println("Arquivo aberto:", params.TextDocument.URI)
	log.Println("Conteúdo:\n", params.TextDocument.Text)
	return nil
}

func didChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	log.Println("File Changed:", params.TextDocument.URI)
	for _, change := range params.ContentChanges {
		if contentChange, ok := change.(protocol.TextDocumentContentChangeEvent); ok {
			log.Println("Change:\n", contentChange.Text)
		}
	}
	return nil
}

func didClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	log.Println("File closed:", params.TextDocument.URI)
	return nil
}

func completion(ctx *glsp.Context, params *protocol.CompletionParams) (any, error) {
	kind := protocol.CompletionItemKindKeyword

	items := []protocol.CompletionItem{
		{
			Label:      "query",
			Kind:       &kind, // Ponteiro correto
			Detail:     ptr("GraphQL operation type"),
			InsertText: ptr("query"),
		},
		{
			Label:      "mutation",
			Kind:       &kind, // Ponteiro correto
			Detail:     ptr("GraphQL operation type"),
			InsertText: ptr("mutation"),
		},
		{
			Label:      "subscription",
			Kind:       &kind, // Ponteiro correto
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
	word := "query" // Simulado
	content := fmt.Sprintf("Information about the word: `%s`", word)

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: content,
		},
	}, nil
}

// Helper para criar *string
func ptr(s string) *string {
	return &s
}
