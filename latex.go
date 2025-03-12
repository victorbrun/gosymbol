package gosymbol

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/victorbrun/gosymbol/latex/lexer"
	"github.com/victorbrun/gosymbol/latex/token"
)

type (
	prefixParseFn func() Expr
	infixParseFn  func(Expr) Expr
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func parserNew(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseInt)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, func() Expr { return p.parseGroupedExpression(token.LPAREN) })
	p.registerPrefix(token.LSPAREN, func() Expr { return p.parseGroupedExpression(token.LSPAREN) })
	p.registerPrefix(token.LBRACE, func() Expr { return p.parseGroupedExpression(token.LBRACE) })
	p.registerPrefix(token.LOG, p.parsePrefixExpression)
	p.registerPrefix(token.EXP, p.parsePrefixExpression)
	p.registerPrefix(token.PI, p.parseIdentifier)
	p.registerPrefix(token.E, p.parseIdentifier)
	p.registerPrefix(token.FRAC, p.parseFrac)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.APEX, p.parseInfixExpression)
	return p
}

var operators = []token.TokenType{
	token.BACKSLASH,
	token.PLUS,
	token.MINUS,
	token.ASTERISK,
	token.APEX,
	token.SLASH,
}

var mathConstants = map[token.TokenType]Expr{
	token.E:  E,
	token.PI: PI,
}

func (p *Parser) getErrors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next tojen to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

const (
	_ int = iota
	precedenceLowest
	precedenceSum
	precedenceProduct
	precedencePower
	precedencePrefix
)

func (p *Parser) parseExpression(precedence int) Expr {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.EOF) && precedence < p.peekPrecedence() {
		if !p.curTokenIsOperator() && !p.peekTokenIsOperator() {
			p.nextToken()
			leftExp = p.parseImplicitMultiplication(leftExp)
		} else {
			infix := p.infixParseFns[p.peekToken.Type]
			if infix == nil {
				return leftExp
			}
			p.nextToken()
			leftExp = infix(leftExp)
		}
	}
	return leftExp
}

func (p *Parser) parseIdentifier() Expr {
	if mathConst, ok := mathConstants[p.curToken.Type]; ok {
		return mathConst
	}
	return Var(VarName(p.curToken.Literal))
}

func (p *Parser) parseInt() Expr {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	return Int(value)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() Expr {
	var prefixConstructor func(Expr) Expr
	if p.curToken.Type == token.MINUS {
		prefixConstructor = func(expr Expr) Expr { return Sub(Int(0), expr) }
	} else if p.curToken.Type == token.LOG {
		prefixConstructor = func(expr Expr) Expr { return Log(expr) }
	} else if p.curToken.Type == token.EXP {
		prefixConstructor = func(expr Expr) Expr { return Exp(expr) }
	}

	p.nextToken()
	expr := p.parseExpression(precedencePrefix)
	return prefixConstructor(expr)
}

var precedences = map[token.TokenType]int{
	token.PLUS:     precedenceSum,
	token.MINUS:    precedenceSum,
	token.SLASH:    precedenceProduct,
	token.ASTERISK: precedenceProduct,
	token.APEX:     precedencePower,
}

func (p *Parser) peekPrecedence() int {
	if !p.curTokenIsOperator() && !p.peekTokenIsOperator() {
		return precedenceProduct
	}
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return precedenceLowest
}

func (p *Parser) curTokenIsOperator() bool {
	preDelimiters := []token.TokenType{
		token.LPAREN,
		token.LSPAREN,
		token.LBRACE,
	}
	return slices.Contains(operators, p.curToken.Type) || slices.Contains(preDelimiters, p.curToken.Type)
}

func (p *Parser) peekTokenIsOperator() bool {
	postDelimiters := []token.TokenType{
		token.RPAREN,
		token.RSPAREN,
		token.RBRACE,
	}
	return slices.Contains(operators, p.peekToken.Type) || slices.Contains(postDelimiters, p.peekToken.Type)

}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return precedenceLowest
}

func (p *Parser) parseInfixExpression(left Expr) Expr {
	var infixConstructor func(Expr, Expr) Expr
	if p.curToken.Type == token.MINUS {
		infixConstructor = func(left Expr, right Expr) Expr { return Sub(left, right) }
	} else if p.curToken.Type == token.PLUS {
		infixConstructor = func(left Expr, right Expr) Expr { return Add(left, right) }
	} else if p.curToken.Type == token.ASTERISK {
		infixConstructor = func(left Expr, right Expr) Expr { return Mul(left, right) }
	} else if p.curToken.Type == token.SLASH {
		infixConstructor = Div
	} else if p.curToken.Type == token.APEX {
		infixConstructor = func(left Expr, right Expr) Expr { return Pow(left, right) }
	}

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)
	return infixConstructor(left, right)
}

func (p *Parser) parseImplicitMultiplication(left Expr) Expr {
	right := p.parseExpression(precedenceProduct)
	return Mul(left, right)
}

func (p *Parser) parseGroupedExpression(delimiter token.TokenType) Expr {
	p.nextToken()

	exp := p.parseExpression(precedenceLowest)
	if (delimiter == token.LPAREN && !(p.expectPeek(token.RPAREN))) || (delimiter == token.LSPAREN && !(p.expectPeek(token.RSPAREN))) || (delimiter == token.LBRACE && !(p.expectPeek(token.RBRACE))) {
		return nil
	}
	return exp
}

func (p *Parser) parseFrac() Expr {
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	numerator := p.parseGroupedExpression(p.curToken.Type)
	if numerator == nil {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	denominator := p.parseGroupedExpression(p.curToken.Type)
	if denominator == nil {
		return nil
	}
	p.nextToken()
	return Div(numerator, denominator)
}
