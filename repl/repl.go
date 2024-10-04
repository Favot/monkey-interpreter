package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Favot/monkey-interpreter/lexer"
	"github.com/Favot/monkey-interpreter/token"
)

const PROMPT = "Monkey >> "

func StartRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			break
		}

		line := scanner.Text()
		lexer := lexer.NewLexer(line)

		for currentToken := lexer.NextToken(); currentToken.Type != token.EOF; currentToken = lexer.NextToken() {
			fmt.Printf("%+v\n", currentToken)
		}
	}
}
