package parser

import (
	"fmt"
)

type Parser struct {
	lex     *lexer
	curtok  token
	nexttok token
	errors  []string
}

func NewParser(l *lexer) *Parser {
	temp := Parser{lex: l, errors: []string{}}
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
	default:
		return nil
	}
}

func (p *Parser) parselet() *let {
	stmt := &let{tok: p.curtok}

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
