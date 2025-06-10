package parser

import (
	"testing"

	"github.com/Jitesh117/brainrotLang-interpreter/ast"
	"github.com/Jitesh117/brainrotLang-interpreter/lexer"
)

func TestYeetStatements(t *testing.T) {
	input := `
	yeet x = 5;
	yeet y = 10;
	yeet foobar = 117;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testYeetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testYeetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "yeet" {
		t.Errorf("s.TokenLiteral not 'yeet'. got=%q", s.TokenLiteral())
		return false
	}

	yeetStmt, ok := s.(*ast.YeetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if yeetStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, yeetStmt.Name.Value)
		return false
	}
	if yeetStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, yeetStmt.Name)
		return false
	}
	return true
}

func TestSlayStatements(t *testing.T) {
	input := `
		slay 5;
		slay 10;
		slay 993322;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.SlayStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "slay" {
			t.Errorf("returnStmt.TokenLiteral not 'slay', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "117;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expt not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.TokenLiteral() != "117" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "117", literal.TokenLiteral())
	}
}
