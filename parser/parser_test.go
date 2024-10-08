package parser

import (
	"testing"

	"github.com/Favot/monkey-interpreter/abstractSyntaxTree"
	"github.com/Favot/monkey-interpreter/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x  5;
		let  = 10;
		let 838383;
	`

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.parseProgram()
	checkParserError(t, parser)

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
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}
func testLetStatement(t *testing.T, statement abstractSyntaxTree.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		return false
	}
	letStatement, ok := statement.(*abstractSyntaxTree.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", statement)
		return false
	}
	if letStatement.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStatement.Name.Value)
		return false
	}
	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStatement.Name)
		return false
	}
	return true
}

func checkParserError(t *testing.T, parser *Parser) {
	errors := parser.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, message := range errors {
		t.Errorf("parser error: %q", message)
	}

	t.FailNow()

}
