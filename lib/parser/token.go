package parser

type tokenType string

type token struct {
	ttype   tokenType
	literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"

	// 1343456
	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

var keywords = map[string]tokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func newtoken(tokentype tokenType, ch byte) token {
	return token{ttype: tokentype, literal: string(ch)}
}

func lookupIdent(ident string) tokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
