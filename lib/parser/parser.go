package parser

import (
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[tokenType]int{
	EQ:       EQUALS,
	NOTEQ:    EQUALS,
	LT:       LESSGREATER,
	GT:       LESSGREATER,
	PLUS:     SUM,
	MINUS:    SUM,
	ASTERISK: PRODUCT,
	SLASH:    PRODUCT,
}

type (
	prefixparse func() expression
	infixparse  func(expression) expression
)

type Parser struct {
	lex    *lexer
	errors []string

	curtok  token
	nexttok token

	prefixparsefns map[tokenType]prefixparse
	infixparsefns  map[tokenType]infixparse
}

func NewParser(input string) *Parser {
	l := newlexer(input)
	temp := Parser{lex: l, errors: []string{}}

	temp.curtok = l.next()
	temp.nexttok = l.next()

	temp.prefixparsefns = make(map[tokenType]prefixparse)
	temp.registerprefix(IDENT, temp.parseident)
	temp.registerprefix(INT, temp.parseintliteral)
	temp.registerprefix(MINUS, temp.parseprefixexpr)
	temp.registerprefix(BANG, temp.parseprefixexpr)
	temp.registerprefix(TRUE, temp.parseboolexpr)
	temp.registerprefix(FALSE, temp.parseboolexpr)

	temp.infixparsefns = make(map[tokenType]infixparse)
	temp.registerinfix(PLUS, temp.parseinfixexpr)
	temp.registerinfix(MINUS, temp.parseinfixexpr)
	temp.registerinfix(ASTERISK, temp.parseinfixexpr)
	temp.registerinfix(SLASH, temp.parseinfixexpr)
	temp.registerinfix(LT, temp.parseinfixexpr)
	temp.registerinfix(GT, temp.parseinfixexpr)
	temp.registerinfix(EQ, temp.parseinfixexpr)
	temp.registerinfix(NOTEQ, temp.parseinfixexpr)

	return &temp
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) Parse() *program {
	prog := &program{}
	prog.statements = []statement{}

	for p.curtok.ttype != EOF {
		stmt := p.parsestatement()

		if stmt != nil {
			prog.statements = append(prog.statements, stmt)
		}
		p.next()
	}
	return prog
}

func (p *Parser) next() {
	p.curtok = p.nexttok
	p.nexttok = p.lex.next()
}

func (p *Parser) parsestatement() statement {
	switch p.curtok.ttype {
	case LET:
		return p.parselet()
	case RETURN:
		return p.parsereturn()
	default:
		return p.parseexprstatement()
	}
}

func (p *Parser) parselet() *letstatement {
	stmt := &letstatement{tok: p.curtok}

	if !p.expect(IDENT) {
		return nil
	}

	stmt.name = &identifier{tok: p.curtok, value: p.curtok.literal}

	if !p.expect(ASSIGN) {
		return nil
	}

	for !p.curtokis(SEMICOLON) {
		p.next()
	}
	return stmt
}

func (p *Parser) parsereturn() *returnstatement {
	stmt := &returnstatement{tok: p.curtok}
	p.next()
	for !p.curtokis(SEMICOLON) {
		p.next()
	}
	return stmt
}

func (p *Parser) parseexprstatement() *expressionstatement {
	stmt := &expressionstatement{tok: p.curtok}
	stmt.expr = p.parseexpr(LOWEST)

	if p.nexttokis(SEMICOLON) {
		p.next()
	}

	return stmt
}

func (p *Parser) parseident() expression {
	return &identifier{tok: p.curtok, value: p.curtok.literal}
}

func (p *Parser) parseexpr(precedence int) expression {
	prefix := p.prefixparsefns[p.curtok.ttype]

	if prefix == nil {
		p.noprefixfound(p.curtok.ttype)
		return nil
	}

	left := prefix()

	for !p.nexttokis(SEMICOLON) && precedence < p.peekprecedence() {
		infix := p.infixparsefns[p.nexttok.ttype]
		if infix == nil {
			return left
		}
		p.next()
		left = infix(left)
	}

	return left
}

func (p *Parser) parseintliteral() expression {
	lit := &intliteral{tok: p.curtok}
	value, err := strconv.ParseInt(p.curtok.literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curtok.literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.value = value
	return lit
}

func (p *Parser) parseprefixexpr() expression {
	expr := &prefixexpr{
		tok:      p.curtok,
		operator: p.curtok.literal,
	}

	p.next()
	expr.right = p.parseexpr(PREFIX)
	return expr
}

func (p *Parser) parseinfixexpr(l expression) expression {
	expr := &infixexpr{
		tok:      p.curtok,
		operator: p.curtok.literal,
		left:     l,
	}
	precedence := p.curprecedence()
	p.next()
	expr.right = p.parseexpr(precedence)
	return expr
}

func (p *Parser) parseboolexpr() expression {
	return &boolexpr{tok: p.curtok, value: p.curtokis(TRUE)}
}

func (p *Parser) curtokis(t tokenType) bool {
	return p.curtok.ttype == t
}

func (p *Parser) nexttokis(t tokenType) bool {
	return p.nexttok.ttype == t
}

func (p *Parser) expect(t tokenType) bool {
	if p.nexttokis(t) {
		p.next()
		return true
	}
	p.peekerror(t)
	return false
}

func (p *Parser) peekerror(t tokenType) {
	msg := fmt.Sprintf("expected next token is %s, got %s instead", t, p.nexttok.ttype)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noprefixfound(t tokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekprecedence() int {
	if p, ok := precedences[p.nexttok.ttype]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curprecedence() int {
	if p, ok := precedences[p.curtok.ttype]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerprefix(tok tokenType, fn prefixparse) {
	p.prefixparsefns[tok] = fn
}

func (p *Parser) registerinfix(tok tokenType, fn infixparse) {
	p.infixparsefns[tok] = fn
}
