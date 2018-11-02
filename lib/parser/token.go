package parser

type tokenType string

type token struct {
	ttype   tokenType
	literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LT        = "<"
	GT        = ">"
	BANG      = "!"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	EQ    = "=="
	NOTEQ = "!="
)

var keywords = map[string]tokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
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
