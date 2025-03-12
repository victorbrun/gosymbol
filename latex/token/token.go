package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456

	// Operators
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	APEX     = "^"
	FRAC     = "FRAC"
	LOG      = "LOG"
	EXP      = "EXP"

	// Delimiters
	BACKSLASH = "\\"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LSPAREN   = "["
	RSPAREN   = "]"

	//math constants
	E  = "e"
	PI = "\\pi"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"\\frac": FRAC,
	"\\log":  LOG,
	"\\exp":  EXP,
}
var mathConstants = map[string]TokenType{
	"e":    E,
	"\\pi": PI,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	if tok, ok := mathConstants[ident]; ok {
		return tok
	}
	return IDENT
}
