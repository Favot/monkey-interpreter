package parser

import (
	"fmt"
	"strconv"

	"github.com/Favot/monkey-interpreter/abstractSyntaxTree"
	"github.com/Favot/monkey-interpreter/lexer"
	"github.com/Favot/monkey-interpreter/token"
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	currentToken token.Token
	lookahead    token.Token

	prefixParseFunctions map[token.TokenType]prefixParseFunction
	infixParseFunctions  map[token.TokenType]infixParseFunction
}

type (
	prefixParseFunction func() abstractSyntaxTree.Expression
	infixParseFunction  func(abstractSyntaxTree.Expression) abstractSyntaxTree.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESS_GREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQUALS:       EQUALS,
	token.NOT_EQUALS:   EQUALS,
	token.LESS_THAN:    LESS_GREATER,
	token.GREATER_THAN: LESS_GREATER,
	token.ADD:          SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
}

func NewParser(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer, errors: []string{}}

	parser.nextToken()
	parser.nextToken()

	parser.prefixParseFunctions = make(map[token.TokenType]prefixParseFunction)
	parser.regiesterPrefix(token.IDENT, parser.parseIndifier)
	parser.regiesterPrefix(token.INT, parser.parseIntegerLiteral)
	parser.regiesterPrefix(token.BANG, parser.parsePrefixExpression)
	parser.regiesterPrefix(token.MINUS, parser.parsePrefixExpression)

	parser.infixParseFunctions = make(map[token.TokenType]infixParseFunction)
	parser.registerInfix(token.ADD, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQUALS, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQUALS, parser.parseInfixExpression)
	parser.registerInfix(token.LESS_THAN, parser.parseInfixExpression)
	parser.registerInfix(token.GREATER_THAN, parser.parseInfixExpression)

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
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseLetStatement() *abstractSyntaxTree.LetStatement {
	letStatement := &abstractSyntaxTree.LetStatement{Token: parser.currentToken}

	if !parser.expectPeek(token.IDENT) {
		return nil
	}

	letStatement.Name = &abstractSyntaxTree.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

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

func (parser *Parser) regiesterPrefix(tokenType token.TokenType, function prefixParseFunction) {
	parser.prefixParseFunctions[tokenType] = function
}

func (parser *Parser) registerInfix(tokenType token.TokenType, function infixParseFunction) {
	parser.infixParseFunctions[tokenType] = function
}

func (parser *Parser) parseExpressionStatement() *abstractSyntaxTree.ExpressionStatement {
	statement := &abstractSyntaxTree.ExpressionStatement{Token: parser.currentToken}

	statement.Expression = parser.parseExpression(LOWEST)

	if parser.peekNextTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpression(precedent int) abstractSyntaxTree.Expression {

	prefix := parser.prefixParseFunctions[parser.currentToken.Type]
	if prefix == nil {
		parser.noPrefixParseFunctionError(parser.currentToken.Type)
		return nil
	}

	leftExpression := prefix()

	for !parser.peekNextTokenIs(token.SEMICOLON) && precedent < parser.peekPrecedence() {
		infix := parser.infixParseFunctions[parser.lookahead.Type]

		if infix == nil {
			return leftExpression
		}

		parser.nextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (parser *Parser) parseIndifier() abstractSyntaxTree.Expression {
	return &abstractSyntaxTree.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() abstractSyntaxTree.Expression {
	literal := &abstractSyntaxTree.IntegerLiteral{Token: parser.currentToken}

	value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)

	if err != nil {
		message := fmt.Sprintf("coulf not parse %q as integer", parser.currentToken.Literal)
		parser.errors = append(parser.errors, message)
		return nil
	}

	literal.Value = value

	return literal
}

func (parser *Parser) noPrefixParseFunctionError(tokenType token.TokenType) {
	message := fmt.Sprintf("no prefix parse fuinction for %s found", tokenType)

	parser.errors = append(parser.errors, message)

}

func (parser *Parser) parsePrefixExpression() abstractSyntaxTree.Expression {

	expression := &abstractSyntaxTree.PrefixEpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
	}

	parser.nextToken()

	expression.Rigth = parser.parseExpression(PREFIX)

	return expression

}

func (parser *Parser) peekPrecedence() int {
	if precedence, ok := precedences[parser.lookahead.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (parser *Parser) currentPrecedence() int {
	if precedence, ok := precedences[parser.currentToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (parser *Parser) parseInfixExpression(left abstractSyntaxTree.Expression) abstractSyntaxTree.Expression {
	expression := &abstractSyntaxTree.InfixEpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
		Left:     left,
	}

	precedence := parser.currentPrecedence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)

	return expression
}
