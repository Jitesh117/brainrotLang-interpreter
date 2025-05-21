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

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
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
		returnStmt, ok := stmt.(*ast.ReturnStatement)
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
