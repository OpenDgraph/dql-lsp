package lsp

import (
	"fmt"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func hover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	doc, ok := GetDocument(params.TextDocument.URI)
	if !ok {
		return nil, fmt.Errorf("document not found")
	}

	offset, err := positionToOffset(doc.Content, int(params.Position.Line), int(params.Position.Character))
	if err != nil {
		return nil, err
	}

	word := extractWordAtOffset(doc.Content, offset)
	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: fmt.Sprintf("Information about `%s`", word),
		},
	}, nil
}
