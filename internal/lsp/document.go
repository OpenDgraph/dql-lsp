package lsp

import protocol "github.com/tliron/glsp/protocol_3_16"

type DocumentState struct {
	Content  string
	Modified bool
}

var documents = make(map[protocol.DocumentUri]DocumentState)

// Funções de manipulação de documentos
func OpenDocument(uri protocol.DocumentUri, text string) {
	documents[uri] = DocumentState{Content: text, Modified: false}
}

func UpdateDocument(uri protocol.DocumentUri, text string) {
	documents[uri] = DocumentState{Content: text, Modified: true}
}

func SaveDocument(uri protocol.DocumentUri, text string) {
	documents[uri] = DocumentState{Content: text, Modified: false}
}

func CloseDocument(uri protocol.DocumentUri) {
	delete(documents, uri)
}

func GetDocument(uri protocol.DocumentUri) (DocumentState, bool) {
	doc, ok := documents[uri]
	return doc, ok
}
