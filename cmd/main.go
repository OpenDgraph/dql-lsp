package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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

type DocumentState struct {
	Content  string
	Modified bool
}

var documents = make(map[protocol.DocumentUri]DocumentState)

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
		TextDocumentDidSave:             didSave,
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
	documents[params.TextDocument.URI] = DocumentState{
		Content:  params.TextDocument.Text,
		Modified: false,
	}
	return nil
}

func didSave(ctx *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	log.Println("[DidSave] File saved:", params.TextDocument.URI)

	if params.Text != nil {
		documents[params.TextDocument.URI] = DocumentState{
			Content:  *params.Text,
			Modified: false,
		}
		log.Println("[DidSave] Updated in-memory content on save")
	}
	return nil
}

func didChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	log.Println("[DidChange] File changed:", params.TextDocument.URI)

	for _, changeRaw := range params.ContentChanges {
		if change, ok := changeRaw.(protocol.TextDocumentContentChangeEventWhole); ok {
			documents[params.TextDocument.URI] = DocumentState{
				Content:  change.Text,
				Modified: true,
			}
			log.Println("[DidChange] Full change updated in-memory content:\n", change.Text)
		} else {
			log.Println("[DidChange] Unexpected change type:", changeRaw)
		}
	}

	return nil
}

func didClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	log.Println("[DidClose] File closed:", params.TextDocument.URI)
	delete(documents, params.TextDocument.URI)
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

func completionOld(ctx *glsp.Context, params *protocol.CompletionParams) (any, error) {
	log.Println("[Completion] Triggered at:", params.TextDocument.URI, "Position:", params.Position)

	kind := protocol.CompletionItemKindKeyword

	items := []protocol.CompletionItem{
		{
			Label:      "query",
			Kind:       &kind,
			Detail:     ptr("DQL operation type"),
			InsertText: ptr("query"),
		},
		{
			Label:      "mutation",
			Kind:       &kind,
			Detail:     ptr("DQL operation type"),
			InsertText: ptr("mutation"),
		},
	}

	return protocol.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}, nil
}

func completion(ctx *glsp.Context, params *protocol.CompletionParams) (any, error) {
	uri := params.TextDocument.URI
	position := params.Position

	doc := documents[uri]                            // Acessa o DocumentState
	context := analyzeContext(doc.Content, position) // Usa o campo Content

	log.Println("[Completion] Triggered at:", uri, "Position:", position, "Context:", context)

	var items []protocol.CompletionItem

	switch context {
	case "empty":
		items = append(items, createSnippet("query { }", "Create a new query block"))
	case "insideQuery":
		items = append(items, createSnippet("me (func: ) { }", "Add query block"))
	case "insideFunc":
		items = append(items, createKeyword("eq()", "Equality function"))
		items = append(items, createKeyword("lt()", "Less than function"))
	case "insideEq":
		items = append(items, createKeyword("name", "Field name"))
	}

	return protocol.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}, nil
}

func hover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {

	doc, ok := documents[params.TextDocument.URI]
	if !ok {
		return nil, fmt.Errorf("document not found")
	}

	offset, err := positionToOffset(doc.Content, int(params.Position.Line), int(params.Position.Character))
	if err != nil {
		return nil, err
	}

	word := extractWordAtOffset(doc.Content, offset)
	if word == "" {
		word = "(no word)"
	}

	log.Println("[Hover] Hover requested at:", params.TextDocument.URI, "Position:", params.Position, "Word:", word)

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

func positionToOffset(text string, line int, character int) (int, error) {
	lines := strings.Split(text, "\n")
	if line < 0 || line >= len(lines) {
		return 0, fmt.Errorf("line out of range")
	}

	lineText := lines[line]
	if character < 0 || character > len(lineText) {
		return 0, fmt.Errorf("character out of range")
	}

	offset := 0
	for i := 0; i < line; i++ {
		offset += len(lines[i]) + 1
	}
	offset += character

	return offset, nil
}

func extractWordAtOffset(text string, offset int) string {
	if offset >= len(text) {
		offset = len(text) - 1
	}
	if offset < 0 {
		offset = 0
	}

	start := offset
	for start > 0 {
		if isWordSeparator(text[start-1]) {
			break
		}
		start--
	}

	end := offset
	for end < len(text) {
		if isWordSeparator(text[end]) {
			break
		}
		end++
	}

	return text[start:end]
}

func isWordSeparator(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '(' || ch == ')' || ch == '{' || ch == '}' || ch == ',' || ch == ';'
}

func createSnippet(insert string, detail string) protocol.CompletionItem {
	kind := protocol.CompletionItemKindSnippet
	insertTextFormat := protocol.InsertTextFormatSnippet

	return protocol.CompletionItem{
		Label:            insert,
		Kind:             &kind,
		Detail:           ptr(detail),
		InsertText:       ptr(insert),
		InsertTextFormat: &insertTextFormat,
	}
}
func createKeyword(keyword string, detail string) protocol.CompletionItem {
	kind := protocol.CompletionItemKindKeyword
	return protocol.CompletionItem{
		Label:      keyword,
		Kind:       &kind,
		Detail:     ptr(detail),
		InsertText: ptr(keyword),
	}
}

func analyzeContext(text string, position protocol.Position) string {
	if len(strings.TrimSpace(text)) == 0 {
		return "empty"
	}

	lines := strings.Split(text, "\n")
	if int(position.Line) >= len(lines) {
		return "unknown"
	}

	textUpToCursor := strings.Join(lines[:position.Line+1], "\n")

	if strings.Contains(textUpToCursor, "query") {
		if strings.Contains(textUpToCursor, "func:") {
			return "insideFunc"
		}
		return "insideQuery"
	}

	if len(strings.TrimSpace(textUpToCursor)) == 0 {
		return "empty"
	}

	return "unknown"
}
