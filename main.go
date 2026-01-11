package main

import (
	"bufio"
	"errors"
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

	path := os.Args[1]
	input, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	lex := lexer.New(string(input))
	tokens, err := lex.Tokenize()
	if err != nil {
		var lexerErr *lexer.LexerError
		if errors.As(err, &lexerErr) {
			fmt.Println(lexerErr.Format(path))
		} else {
			fmt.Printf("Generic Error: %v\n", err)
		}
		os.Exit(1)
	}

	p := parser.New(tokens)
	env := parser.NewDefaultEnvironment()
	if err := p.Execute(env); err != nil {
		var parserErr *parser.ParserError
		var runtimeErr *parser.RuntimeError
		if errors.As(err, &runtimeErr) {
			fmt.Println(runtimeErr.Format(path, string(input)))
		} else if errors.As(err, &parserErr) {
			fmt.Println(parserErr.Format(path))
		} else {
			fmt.Printf("Generic Error: %v\n", err)
		}
		os.Exit(1)
	}
}

func repl() {
	env := parser.NewDefaultEnvironment()
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
		tokens, err := lex.Tokenize()
		if err != nil {
			fmt.Println(err)
			continue
		}

		p := parser.New(tokens)
		if err := p.Execute(env); err != nil {
			fmt.Println(err)
			continue
		}
	}
}
