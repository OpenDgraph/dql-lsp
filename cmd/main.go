package main

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/OpenDgraph/dql-lsp/internal/lsp"
	"github.com/tliron/commonlog"

	_ "github.com/tliron/commonlog/simple"
)

func init() {
	isDebug := strings.ToLower(os.Getenv("DEBUG")) == "true"
	log.Println("[DEBUG]", isDebug)

	// Setup file logging
	if isDebug {
		logFile, err := os.OpenFile("./testing/dql-lsp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open log file:", err)
		}

		log.SetOutput(logFile)
	} else {
		log.Println("[Init] Logging only to stdout")
		log.SetOutput(io.Discard)
	}
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()

	commonlog.Configure(1, nil) // Default log verbosity

	lsp.RunServer()

	log.Println("==== Starting DQL Language Server ====")
}

// func completionOld(ctx *glsp.Context, params *protocol.CompletionParams) (any, error) {
// 	log.Println("[Completion] Triggered at:", params.TextDocument.URI, "Position:", params.Position)

// 	kind := protocol.CompletionItemKindKeyword

// 	items := []protocol.CompletionItem{
// 		{
// 			Label:      "query",
// 			Kind:       &kind,
// 			Detail:     ptr("DQL operation type"),
// 			InsertText: ptr("query"),
// 		},
// 		{
// 			Label:      "mutation",
// 			Kind:       &kind,
// 			Detail:     ptr("DQL operation type"),
// 			InsertText: ptr("mutation"),
// 		},
// 	}

// 	return protocol.CompletionList{
// 		IsIncomplete: false,
// 		Items:        items,
// 	}, nil
// }
