package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/Favot/monkey-interpreter/repl"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Monkey programming language.\n", cases.Title(language.English).String(user.Username))

	fmt.Println("Let's get started!")

	repl.StartRepl(os.Stdin, os.Stdout)
}
