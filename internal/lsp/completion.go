package lsp

import (
	"log"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

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
