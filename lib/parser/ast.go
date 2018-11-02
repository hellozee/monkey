package parser

type node interface {
	tokenliteral() string
}

type statement interface {
	node
	statementnode()
}

type expression interface {
	node
	expressionnode()
}

type program struct {
	statements []statement
}

func (p *program) tokenliteral() string {
	if len(p.statements) > 0 {
		return p.statements[0].tokenliteral()
	}
	return ""
}

type let struct {
	tok   token
	name  *identifier
	value expression
}

func (l *let) statementnode()       {}
func (l *let) tokenliteral() string { return l.tok.literal }

type identifier struct {
	tok   token
	value string
}

func (i *identifier) expressionnode()      {}
func (i *identifier) tokenliteral() string { return i.tok.literal }
