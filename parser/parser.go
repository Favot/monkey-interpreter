package parser

import (
	"fmt"

	"github.com/Favot/monkey-interpreter/abstractSyntaxTree"
	"github.com/Favot/monkey-interpreter/lexer"
	"github.com/Favot/monkey-interpreter/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	lookahead    token.Token
	errors       []string
}

func NewParser(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer, errors: []string{}}

	parser.nextToken()
	parser.nextToken()

	return parser
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) peekError(nextToken token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", nextToken, parser.currentToken.Type)

	parser.errors = append(parser.errors, message)
}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.lookahead
	parser.lookahead = parser.lexer.NextToken()
}

func (parser *Parser) parseProgram() *abstractSyntaxTree.Program {
	program := &abstractSyntaxTree.Program{}
	program.Statements = []abstractSyntaxTree.Statement{}

	for parser.currentToken.Type != token.EOF {
		parserStatement := parser.parseStatement()
		if parserStatement != nil {
			program.Statements = append(program.Statements, parserStatement)
		}
		parser.nextToken()
	}

	return program

}

func (parser *Parser) parseStatement() abstractSyntaxTree.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return nil
	}
}

func (parser *Parser) parseLetStatement() *abstractSyntaxTree.LetStatement {
	letStatement := &abstractSyntaxTree.LetStatement{Token: parser.currentToken}

	if !parser.expectPeek(token.IDENT) {
		return nil
	}

	letStatement.Name = &abstractSyntaxTree.Identifier{Token: parser.currentToken, Value: parser.currentToken.Value}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}
	return letStatement
}

func (parser *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return parser.currentToken.Type == tokenType
}

func (parser *Parser) peekNextTokenIs(tokenType token.TokenType) bool {
	return parser.lookahead.Type == tokenType
}

func (parser *Parser) expectPeek(tokenType token.TokenType) bool {
	if parser.peekNextTokenIs(tokenType) {
		parser.nextToken()
		return true
	} else {
		parser.peekError(tokenType)
		return false
	}
}

func (parser *Parser) parseReturnStatement() *abstractSyntaxTree.ReturnStatement {

	statement := &abstractSyntaxTree.ReturnStatement{
		Token: parser.currentToken,
	}

	parser.nextToken()

	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}
