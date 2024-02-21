package lexer

import (
	"testing"

	"github.com/daneofmanythings/diceroni/pkg/interpreter/token"
)

func TestNextToken(t *testing.T) {
	input := `(4 + 3d5ma4) - d6mq2mi2 * d20 - 4/2;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LPAREN, "("},
		{token.INT, "4"},
		{token.PLUS, "+"},
		{token.INT, "3"},
		{token.DICE, "5"},
		{token.DICEMAX, "4"},
		{token.RPAREN, ")"},
		{token.MINUS, "-"},
		{token.DICE, "6"},
		{token.DICEQUANT, "2"},
		{token.DICEMIN, "2"},
		{token.ASTERISK, "*"},
		{token.DICE, "20"},
		{token.MINUS, "-"},
		{token.INT, "4"},
		{token.SLASH, "/"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
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
