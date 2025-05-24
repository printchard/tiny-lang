package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/printchard/tiny-lang/lexer"
	"github.com/printchard/tiny-lang/parser"
)

func main() {
	if len(os.Args) < 2 {
		repl()
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
	p.Execute(nil)
}

func repl() {
	env := parser.Environment{
		Variables: make(map[string]float64),
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("tiny-lang> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		lex := lexer.New(input)
		tokens := lex.Tokenize()

		p := parser.New(tokens)
		p.Execute(&env)
	}
}
