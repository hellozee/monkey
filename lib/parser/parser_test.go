package parser

import "testing"

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
let y = 10;
let foo = 838383;`

	l := newlexer(input)
	p := NewParser(l)

	program := p.Parse()
	checkparseerrors(t, p)

	if program == nil {
		t.Fatalf("Parse() returned nil")
	}

	if len(program.statements) != 3 {
		t.Fatalf("program.statements does not contain 3 statements, got %d", len(program.statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		stmt := program.statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s statement, name string) bool {
	if s.tokenliteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.tokenliteral())
		return false
	}
	letstmt, ok := s.(*let)

	if !ok {
		t.Errorf("s not *ast.letstatement. got=%T", s)
		return false
	}

	if letstmt.name.value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letstmt.name.value)
		return false
	}

	if letstmt.name.tokenliteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letstmt.name)
		return false
	}
	return true
}

func checkparseerrors(t *testing.T, p *Parser) {
	errors := p.errors

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
