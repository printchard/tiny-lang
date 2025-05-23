package main

import (
	"fmt"

	"github.com/printchard/tiny-lang/lexer"
)

func main() {
	lex := lexer.New(`let x := -5 + y * (3 + 2)
print x`)
	tokens := lex.Tokenize()
	fmt.Println(tokens)
}
