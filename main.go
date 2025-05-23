package main

import (
	"fmt"
	"os"

	"github.com/printchard/tiny-lang/lexer"
	"github.com/printchard/tiny-lang/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tiny-lang <input-file>")
		return
	}

	inputFile := os.Args[1]
	input, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	lex := lexer.New(string(input))
	tokens := lex.Tokenize()

	p := parser.New(tokens)
	p.Execute()
}
