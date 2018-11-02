package parser

type Parser struct {
	lex     *lexer
	curtok  token
	nexttok token
}

func NewParser(l *lexer) *Parser {
	temp := Parser{lex: l}
	return &temp
}

func (p *Parser) next() {
	p.curtok = p.nexttok
	p.nexttok = p.lex.next()
}

func (p *Parser) Parse() *program {
	return nil
}
