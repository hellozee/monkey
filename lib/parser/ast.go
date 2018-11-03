package parser

import (
	"bytes"
)

type node interface {
	tokenliteral() string
	tostring() string
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

func (p *program) tostring() string {
	var out bytes.Buffer

	for _, s := range p.statements {
		out.WriteString(s.tostring())
	}

	return out.String()
}

type letstatement struct {
	tok   token
	name  *identifier
	value expression
}

func (l *letstatement) statementnode()       {}
func (l *letstatement) tokenliteral() string { return l.tok.literal }

func (l *letstatement) tostring() string {
	var out bytes.Buffer

	out.WriteString(l.tokenliteral() + " ")
	out.WriteString(l.name.tostring())
	out.WriteString(" = ")

	if l.value != nil {
		out.WriteString(l.value.tostring())
	}

	out.WriteString(";")

	return out.String()
}

type returnstatement struct {
	tok   token
	value expression
}

func (r *returnstatement) statementnode()       {}
func (r *returnstatement) tokenliteral() string { return r.tok.literal }

func (r *returnstatement) tostring() string {
	var out bytes.Buffer
	out.WriteString(r.tokenliteral() + " ")
	if r.value != nil {
		out.WriteString(r.value.tostring())
	}

	out.WriteString(";")

	return out.String()
}

type expressionstatement struct {
	tok  token
	expr expression
}

func (e *expressionstatement) statementnode()       {}
func (e *expressionstatement) tokenliteral() string { return e.tok.literal }

func (e *expressionstatement) tostring() string {
	if e.expr != nil {
		return e.expr.tostring()
	}
	return ""
}

type identifier struct {
	tok   token
	value string
}

func (i *identifier) expressionnode()      {}
func (i *identifier) tokenliteral() string { return i.tok.literal }
func (i *identifier) tostring() string     { return i.value }

type intliteral struct {
	tok   token
	value int64
}

func (i *intliteral) expressionnode()      {}
func (i *intliteral) tokenliteral() string { return i.tok.literal }
func (i *intliteral) tostring() string     { return i.tok.literal }

type prefixexpr struct {
	tok      token
	operator string
	right    expression
}

func (p *prefixexpr) expressionnode()      {}
func (p *prefixexpr) tokenliteral() string { return p.tok.literal }

func (p *prefixexpr) tostring() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.operator)
	out.WriteString(p.right.tostring())
	out.WriteString(")")
	return out.String()
}

type infixexpr struct {
	tok      token
	left     expression
	operator string
	right    expression
}

func (i *infixexpr) expressionnode()      {}
func (i *infixexpr) tokenliteral() string { return i.tok.literal }

func (i *infixexpr) tostring() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.operator)
	out.WriteString(i.right.tostring())
	out.WriteString(")")
	return out.String()
}
