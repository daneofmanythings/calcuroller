package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var Keywords = map[string]TokenType{
	"mi": DICEMIN,
	"ma": DICEMAX,
	"kh": DICEKEEPHIGHEST,
	"kl": DICEKEEPLOWEST,
	"qu": DICEQUANT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}

var DiceMods []TokenType = []TokenType{
	METATAG,
	DICEQUANT,
	DICEMAX,
	DICEMIN,
	DICEKEEPLOWEST,
	DICEKEEPHIGHEST,
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT"
	INT   = "INT"
	DICE  = "DICE"
	// 1343456

	// Diceroll Modifiers
	METATAG         = "TAG"
	DICEQUANT       = "QUANT"
	DICEMIN         = "MIN"
	DICEMAX         = "MAX"
	DICEKEEPHIGHEST = "HIGHEST"
	DICEKEEPLOWEST  = "LOWEST"

	// Operators
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	// Delimiters
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"

	// Keywords
	// ADVANTAGE    = "ADVANTAGE"
	// DISADVANDAGE = "DISADVANDAGE"
)
