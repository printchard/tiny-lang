package main

import (
	"fmt"
	"os"

	"github.com/printchard/tiny-lang/lexer"
	"github.com/printchard/tiny-lang/parser"
)

func main() {
	input, err := os.ReadFile("input.tiny")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	lex := lexer.New(string(input))
	tokens := lex.Tokenize()
	fmt.Println(tokens)

	p := parser.New(tokens)
	p.Parse()
}
