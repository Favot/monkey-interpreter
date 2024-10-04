package lexer

import "github.com/Favot/monkey-interpreter/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	currentChar  byte
}

func NewLexer(input string) *Lexer {
	lexer := &Lexer{input: input}

	lexer.readChar()

	return lexer
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.currentChar = 0
	} else {
		lexer.currentChar = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
	lexer.readPosition++
}

func (lexer *Lexer) NextToken() token.Token {
	var currentToken token.Token

	lexer.skipWhitespace()

	switch lexer.currentChar {
	case '=':
		if lexer.peekNextChar() == '=' {
			char := lexer.currentChar
			lexer.readChar()
			currentToken = token.Token{Type: token.EQUALS, Value: string(char) + string(lexer.currentChar)}
		} else {
			currentToken = newToken(token.ASSIGN, lexer.currentChar)
		}
	case '+':
		currentToken = newToken(token.ADD, lexer.currentChar)
	case '-':
		currentToken = newToken(token.MINUS, lexer.currentChar)
	case '!':
		if lexer.peekNextChar() == '=' {
			char := lexer.currentChar
			lexer.readChar()
			currentToken = token.Token{Type: token.NOT_EQUALS, Value: string(char) + string(lexer.currentChar)}
		} else {
			currentToken = newToken(token.BANG, lexer.currentChar)
		}
	case '*':
		currentToken = newToken(token.ASTERISK, lexer.currentChar)
	case '/':
		currentToken = newToken(token.SLASH, lexer.currentChar)
	case '<':
		currentToken = newToken(token.LESS_THAN, lexer.currentChar)
	case '>':
		currentToken = newToken(token.GREATER_THAN, lexer.currentChar)
	case '(':
		currentToken = newToken(token.LEFT_PARENTHESIS, lexer.currentChar)
	case ')':
		currentToken = newToken(token.RIGHT_PARENTHESIS, lexer.currentChar)
	case ',':
		currentToken = newToken(token.COMMA, lexer.currentChar)
	case ';':
		currentToken = newToken(token.SEMICOLON, lexer.currentChar)
	case '{':
		currentToken = newToken(token.LEFT_BRACE, lexer.currentChar)
	case '}':
		currentToken = newToken(token.RIGHT_BRACE, lexer.currentChar)
	case 0:
		currentToken.Value = ""
		currentToken.Type = token.EOF
	default:
		if isLetter(lexer.currentChar) {
			currentToken.Value = lexer.readIdentifer()
			currentToken.Type = token.LookupIdentifier(currentToken.Value)
			return currentToken
		} else if isDigit(lexer.currentChar) {
			currentToken.Value = lexer.readNumber()
			currentToken.Type = token.INT
			return currentToken
		} else {
			currentToken = newToken(token.ILLEGAL, lexer.currentChar)
		}
	}

	lexer.readChar()

	return currentToken

}

func (lexer *Lexer) readIdentifer() string {
	position := lexer.position
	for isLetter(lexer.currentChar) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Value: string(char)}
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position
	for isDigit(lexer.currentChar) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.currentChar == ' ' || lexer.currentChar == '\t' || lexer.currentChar == '\n' || lexer.currentChar == '\r' {
		lexer.readChar()
	}
}

func (lexer *Lexer) peekNextChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPosition]
	}
}
