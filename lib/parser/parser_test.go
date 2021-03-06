package parser

import (
	"fmt"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
			let y = 10;
			let foo = 838383;`

	p := NewParser(input)
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
	letstmt, ok := s.(*letstatement)

	if !ok {
		t.Errorf("s not *letstatement. got=%T", s)
		return false
	}

	if letstmt.name.value != name {
		t.Errorf("letstmt.name.value not '%s'. got=%s", name, letstmt.name.value)
		return false
	}

	if letstmt.name.tokenliteral() != name {
		t.Errorf("s.name not '%s'. got=%s", name, letstmt.name)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	input := `return 5;
			  return 10;
			  return 90234820;`

	p := NewParser(input)
	program := p.Parse()
	checkparseerrors(t, p)

	if program == nil {
		t.Fatalf("Parse() returned nil")
	}

	if len(program.statements) != 3 {
		t.Fatalf("program.statements does not contain 3 statements, got %d", len(program.statements))
	}

	for _, stmt := range program.statements {
		returnstmt, ok := stmt.(*returnstatement)

		if !ok {
			t.Errorf("stmt not *ast.returnstatement. got=%T", stmt)
			continue
		}

		if returnstmt.tokenliteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnstmt.tokenliteral())
		}
	}
}

func TestString(t *testing.T) {
	program := &program{
		statements: []statement{
			&letstatement{
				tok: token{ttype: LET, literal: "let"},
				name: &identifier{
					tok:   token{ttype: IDENT, literal: "foo"},
					value: "foo",
				},
				value: &identifier{
					tok:   token{ttype: IDENT, literal: "bar"},
					value: "bar",
				},
			},
		},
	}

	if program.tostring() != "let foo = bar;" {
		t.Errorf("program.tostring() wrong. got=%q", program.tostring())
	}
}

func TestIdenfierExpressions(t *testing.T) {
	input := "foo"
	p := NewParser(input)
	prog := p.Parse()
	checkparseerrors(t, p)

	if len(prog.statements) != 1 {
		t.Fatalf("program doesn't have enough statements, got %d", len(prog.statements))
	}

	stmt, ok := prog.statements[0].(*expressionstatement)

	if !ok {
		t.Fatalf("prog.statement[0] is not an expression statement, got %T", prog.statements[0])
	}

	testIdent(t, stmt.expr, input)
}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"
	p := NewParser(input)
	prog := p.Parse()
	checkparseerrors(t, p)

	if len(prog.statements) != 1 {
		t.Fatalf("program doesn't have enough statements, got %d", len(prog.statements))
	}

	stmt, ok := prog.statements[0].(*expressionstatement)

	if !ok {
		t.Fatalf("prog.statement[0] is not an expression statement, got %T", prog.statements[0])
	}

	literal, ok := stmt.expr.(*intliteral)

	if !ok {
		t.Fatalf("expression.expr is not an integer literal, got %T", literal.value)
	}

	if literal.value != 5 {
		t.Errorf("literal.value not %d, got %d", 5, literal.value)
	}

	if literal.tokenliteral() != "5" {
		t.Errorf("literal.tokenliteral() not %d got %s", 5, literal.tokenliteral())
	}
}

func testIdent(t *testing.T, expr expression, value string) {
	ident, ok := expr.(*identifier)
	if !ok {
		t.Errorf("expr not *identifier. got=%T", expr)
	}
	if ident.value != value {
		t.Errorf("ident.value not %s. got=%s", value, ident.value)
	}
	if ident.tokenliteral() != value {
		t.Errorf("ident.tokenliteral() not %s. got=%s", value, ident.tokenliteral())
	}
}

func testLiteralExpression(t *testing.T, expr expression, expected interface{}) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, expr, int64(v))
	case int64:
		testIntegerLiteral(t, expr, v)
	case string:
		testIdent(t, expr, v)
	case bool:
		testBoolLiteral(t, expr, v)
	default:
		t.Errorf("type of expr not handled. got=%T", expr)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		intval   int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range tests {
		p := NewParser(tt.input)
		prog := p.Parse()
		checkparseerrors(t, p)

		if len(prog.statements) != 1 {
			t.Fatalf("prog.statements doesn't contain %d statements, got %d\n", 1, len(prog.statements))
		}

		stmt, ok := prog.statements[0].(*expressionstatement)
		if !ok {
			t.Fatalf("prog.statements[0] is not a expression statement, got %T\n", stmt)
		}

		expr, ok := stmt.expr.(*prefixexpr)
		if !ok {
			t.Fatalf("stmt is not prefixexpr. got=%T", stmt.expr)
		}

		if expr.operator != tt.operator {
			t.Fatalf("expr.operator is not '%s'. got=%s", tt.operator, expr.operator)
		}

		testLiteralExpression(t, expr.right, tt.intval)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		left     int64
		operator string
		right    int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range tests {
		p := NewParser(tt.input)
		prog := p.Parse()
		checkparseerrors(t, p)

		if len(prog.statements) != 1 {
			t.Fatalf("prog.statements does not contain %d statements. got=%d\n", 1, len(prog.statements))
		}

		stmt, ok := prog.statements[0].(*expressionstatement)
		if !ok {
			t.Fatalf("prog.statements[0] is not expressionstatement. got=%T",
				prog.statements[0])
		}
		expr, ok := stmt.expr.(*infixexpr)
		if !ok {
			t.Fatalf("expr is not infixexpr. got=%T", stmt.expr)
		}
		testLiteralExpression(t, expr.left, tt.left)
		if expr.operator != tt.operator {
			t.Fatalf("expr.operator is not '%s'. got=%s", tt.operator, expr.operator)
		}

		testLiteralExpression(t, expr.right, tt.right)
	}
}

func testIntegerLiteral(t *testing.T, i expression, value int64) bool {
	integer, ok := i.(*intliteral)
	if !ok {
		t.Errorf("i not intliteral. got=%T", i)
		return false
	}

	if integer.value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integer.value)
		return false
	}
	if integer.tokenliteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integer.tokenliteral())
		return false
	}
	return true
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		p := NewParser(tt.input)
		prog := p.Parse()
		checkparseerrors(t, p)
		actual := prog.tostring()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBoolExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		p := NewParser(tt.input)
		prog := p.Parse()
		checkparseerrors(t, p)

		if len(prog.statements) != 1 {
			t.Fatalf("prog.statements does not contain %d statements. got=%d\n", 1, len(prog.statements))
		}

		stmt, ok := prog.statements[0].(*expressionstatement)
		if !ok {
			t.Fatalf("prog.statements[0] is not expressionstatement. got=%T", prog.statements[0])
		}
		testLiteralExpression(t, stmt.expr, tt.expected)
	}
}

func testBoolLiteral(t *testing.T, expr expression, value bool) {
	bexpr, ok := expr.(*boolexpr)
	if !ok {
		t.Errorf("expr is not boolexpr. got=%T", expr)
	}
	if bexpr.value != value {
		t.Errorf("expr.value is not '%t'. got=%t", value, bexpr.value)
	}
	if bexpr.tokenliteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bexpr.tokenliteral() not %t. got=%s", value, bexpr.tokenliteral())
	}
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
