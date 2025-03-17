package lsp

import (
	"fmt"
	"strings"
)

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
