package parser

type lexer struct {
	input   string
	pos     int
	readPos int
	char    byte
}

func newlexer(data string) *lexer {
	temp := lexer{input: data}
	temp.read()
	return &temp
}

func (l *lexer) read() {
	if l.readPos >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPos]
	}

	l.pos = l.readPos
	l.readPos++
}

func (l *lexer) next() token {
	var tok token

	l.skipspace()

	switch l.char {
	case '=':
		tok = newtoken(ASSIGN, l.char)
	case ';':
		tok = newtoken(SEMICOLON, l.char)
	case '(':
		tok = newtoken(LPAREN, l.char)
	case ')':
		tok = newtoken(RPAREN, l.char)
	case '{':
		tok = newtoken(LBRACE, l.char)
	case '}':
		tok = newtoken(RBRACE, l.char)
	case ',':
		tok = newtoken(COMMA, l.char)
	case '+':
		tok = newtoken(PLUS, l.char)
	case 0:
		tok.literal = ""
		tok.ttype = EOF
	default:
		if isletter(l.char) {
			tok.literal = l.readidentifier()
			tok.ttype = lookupIdent(tok.literal)
			return tok
		} else if isdigit(l.char) {
			tok.literal = l.readnumber()
			tok.ttype = INT
			return tok
		}
		tok = newtoken(ILLEGAL, l.char)
	}

	l.read()
	return tok
}

func (l *lexer) readidentifier() string {
	pos := l.pos
	for isletter(l.char) {
		l.read()
	}
	return l.input[pos:l.pos]
}

func (l *lexer) readnumber() string {
	pos := l.pos
	for isdigit(l.char) {
		l.read()
	}
	return l.input[pos:l.pos]
}

func (l *lexer) skipspace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.read()
	}
}

func isletter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isdigit(char byte) bool {
	return '0' <= char && char <= '9'
}
