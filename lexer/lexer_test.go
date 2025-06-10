package lexer

import (
	"testing"

	"github.com/OG-Open-Source/SDCL/token"
)

func TestNextToken(t *testing.T) {
	input := `---
key: "value"
another_key: 123
# this is a comment
---
你好: "世界"
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.FRONTMATTER_DELIMITER, "---"},
		{token.NEWLINE, "\n"},
		{token.KEY, "key"},
		{token.COLON, ":"},
		{token.STRING, "value"},
		{token.NEWLINE, "\n"},
		{token.KEY, "another_key"},
		{token.COLON, ":"},
		{token.INT, "123"},
		{token.NEWLINE, "\n"},
		{token.COMMENT, " this is a comment"},
		{token.NEWLINE, "\n"},
		{token.FRONTMATTER_DELIMITER, "---"},
		{token.NEWLINE, "\n"},
		{token.KEY, "你好"},
		{token.COLON, ":"},
		{token.STRING, "世界"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}