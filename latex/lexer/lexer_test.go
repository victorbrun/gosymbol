package lexer

import (
	"github.com/victorbrun/gosymbol/latex/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `5e + 5\pi  + ab \ciao\log(a)`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "5"},
		{token.E, "e"},
		{token.PLUS, "+"},
		{token.INT, "5"},
		{token.PI, "\\pi"},
		{token.PLUS, "+"},
		{token.IDENT, "a"},
		{token.IDENT, "b"},
		{token.IDENT, "\\ciao"},
		{token.LOG, "\\log"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.RPAREN, ")"},
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
