package lsp

import (
	"log"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	// version string = "0.0.1"
	handler protocol.Handler
)

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	log.Println("[Initialize] Params:", params)

	capabilities := handler.CreateServerCapabilities()
	capabilities.TextDocumentSync = protocol.TextDocumentSyncKindFull
	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: []string{" ", "\n", "{", "(", ":", "\""},
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

	return result, nil
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
