package parser

import (
	"fmt"
	"testing"

	"github.com/Favot/monkey-interpreter/abstractSyntaxTree"
	"github.com/Favot/monkey-interpreter/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.parseProgram()
	checkParserErrors(t, parser)

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

func checkParserErrors(t *testing.T, parser *Parser) {
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

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
		`
	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)

	program := parser.parseProgram()

	checkParserErrors(t, parser)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statemen, got=%d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*abstractSyntaxTree.ReturnStatement)
		if !ok {
			t.Errorf("Statement not *abstractSyntaxTree.ReturnStatement. got=%T", statement)
			continue
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral not 'return', got %q", returnStatement.TokenLiteral())
		}
	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)

	program := parser.parseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statement. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*abstractSyntaxTree.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	identifier, ok := statement.Expression.(*abstractSyntaxTree.Identifier)

	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", statement.Expression)
	}
	if identifier.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", identifier.Value)
	}
	if identifier.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			identifier.TokenLiteral())
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	lexer := lexer.NewLexer(input)

	parser := NewParser(lexer)

	program := parser.parseProgram()

	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enought staments. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*abstractSyntaxTree.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	litetal, ok := statement.Expression.(*abstractSyntaxTree.IntegerLiteral)

	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", statement.Expression)
	}

	if litetal.Value != 5 {
		t.Fatalf("literal.Value not %d. got=%d", 5, litetal.Value)
	}

	if litetal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", litetal.TokenLiteral())
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, prefixTest := range prefixTests {
		lexer := lexer.NewLexer(prefixTest.input)
		parser := NewParser(lexer)
		program := parser.parseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("progam.Statement does not conatain %d statement. got=%d", 1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*abstractSyntaxTree.ExpressionStatement)

		if !ok {
			t.Fatalf("programe.Statement[0] is not ast.ExpressionStement. got=%T", program.Statements[0])
		}

		expression, ok := statement.Expression.(*abstractSyntaxTree.PrefixEpression)

		if !ok {
			t.Fatalf("statement is not ast.PrefixEpression. got=%T", statement.Expression)
		}

		if expression.Operator != prefixTest.operator {
			t.Fatalf("expression.Operator is not '%s' got=%s", prefixTest.operator, expression.Operator)
		}

		if !testIntegerLiteral(t, expression.Rigth, prefixTest.integerValue) {
			return
		}

	}

}

func testIntegerLiteral(t *testing.T, integerLiteral abstractSyntaxTree.Expression, value int64) bool {

	integer, ok := integerLiteral.(*abstractSyntaxTree.IntegerLiteral)

	if !ok {
		t.Errorf("il not *abs.IntegerLiteral. got=%T", integerLiteral)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d, got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Fatalf("integer.TokenLiteral not %d. got=%s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixEpression(t *testing.T) {

	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, infixTest := range infixTests {

		lexer := lexer.NewLexer(infixTest.input)

		parser := NewParser(lexer)

		program := parser.parseProgram()

		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statement does not containt %d statement. got=%d\n", 1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*abstractSyntaxTree.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statemens[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		expression, ok := statement.Expression.(*abstractSyntaxTree.InfixExpression)

		if !ok {
			t.Fatalf("expression is not ast.InfixEpression. got=%T", statement.Expression)
		}

		if !testIntegerLiteral(t, expression.Left, infixTest.leftValue) {
			return
		}

		if expression.Operator != infixTest.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", infixTest.operator, expression.Operator)
		}

		if !testIntegerLiteral(t, expression.Right, infixTest.rightValue) {
			return
		}

	}

}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		}, {
			"!-a",
			"(!(-a))",
		}, {
			"a + b + c",
			"((a + b) + c)",
		}, {
			"a + b - c",
			"((a + b) - c)",
		}, {
			"a * b * c",
			"((a * b) * c)",
		}, {
			"a * b / c",
			"((a * b) / c)",
		}, {
			"a + b / c",
			"(a + (b / c))",
		}, {
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		}, {
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		}, {
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		}, {
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		}, {
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		}, {
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}
	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.parseProgram()

		checkParserErrors(t, p)
		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
