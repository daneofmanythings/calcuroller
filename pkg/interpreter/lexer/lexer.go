package lexer

import (
	"github.com/daneofmanythings/calcuroller/pkg/interpreter/token"
)

type Lexer struct {
	input        string
	position     int  // points to current char
	peekPosition int  // points after current char
	ch           byte // current char being examined pointed to by position
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	// advances the position pointers and updates ch accordingly
	if l.peekPosition < len(l.input) {
		l.ch = l.input[l.peekPosition]
	} else {
		l.ch = 0
	}
	l.position = l.peekPosition
	l.peekPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '%':
		tok = newToken(token.MODULUS, l.ch)
	case '^':
		tok = newToken(token.CARET, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '[':
		l.readChar()
		tok.Literal = l.readTag()
		tok.Type = token.METATAG
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			if l.ch == 'd' && isDigit(l.peekChar()) { // identifying dice string
				return l.newDiceToken()
			}
			tok.Literal = l.readIdentifier()
			// WARN: Re-evaluate this logic. If there are ever more reserved identifiers added, this will break
			tok.Type = token.LookupIdent(tok.Literal)
			// checking if the identifier corresponds to a dicemod and adjusting accordingly
			if _, ok := token.Keywords[tok.Literal]; ok {
				tok.Literal = l.readNumber()
			}
			// TODO: Add lexing for tags. they are to be surrounded by [], and be read in as simply strings
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) newDiceToken() token.Token {
	var tok token.Token
	tok.Type = token.DICE
	l.readChar() // advance to the integer
	tok.Literal = l.readNumber()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readTag() string {
	position := l.position
	for isTag(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isTag(ch byte) bool {
	return ch != ']'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	if l.peekPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.peekPosition]
	}
}

func (l *Lexer) peekPeekChar() byte {
	if l.peekPosition >= len(l.input)-1 {
		return 0
	} else {
		return l.input[l.peekPosition+1]
	}
}
