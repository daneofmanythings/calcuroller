package lexer

import (
	"testing"

	"github.com/daneofmanythings/diceroni/pkg/interpreter/token"
)

func TestNextToken(t *testing.T) {
	input := `4 + 3 - d6 * 2d20 - 4/2; adv`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "4"},
		{token.PLUS, "+"},
		{token.INT, "3"},
		{token.MINUS, "-"},
		{token.DICE, "6"},
		{token.ASTERISK, "*"},
		{token.INT, "2"},
		{token.DICE, "20"},
		{token.MINUS, "-"},
		{token.INT, "4"},
		{token.SLASH, "/"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.ADVANTAGE, "adv"},
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
