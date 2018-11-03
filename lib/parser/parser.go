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

func NewParser(l *lexer) *Parser {
	temp := Parser{lex: l, errors: []string{}}
	temp.curtok = l.next()
	temp.nexttok = l.next()
	temp.prefixparsefns = make(map[tokenType]prefixparse)
	temp.registerprefix(IDENT, temp.parseident)
	temp.registerprefix(INT, temp.parseintliteral)
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
		return nil
	}

	left := prefix()
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

func (p *Parser) registerprefix(tok tokenType, fn prefixparse) {
	p.prefixparsefns[tok] = fn
}

func (p *Parser) registerinfix(tok tokenType, fn infixparse) {
	p.infixparsefns[tok] = fn
}
