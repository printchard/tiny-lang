package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/printchard/tiny-lang/formatter"
	"github.com/printchard/tiny-lang/lexer"
	"github.com/printchard/tiny-lang/parser"
)

func main() {
	if len(os.Args) < 2 {
		repl()
		return
	}

	path := os.Args[1]

	if path == "fmt" {
		handleFormat(os.Args)
		return
	}

	input, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s", err)
		return
	}
	lex := lexer.New(string(input))
	tokens, err := lex.Tokenize()
	if err != nil {
		var lexerErr *lexer.LexerError
		if errors.As(err, &lexerErr) {
			fmt.Fprintln(os.Stderr, lexerErr.Format(path))
		} else {
			fmt.Fprintf(os.Stderr, "Generic Error: %v\n", err)
		}
		os.Exit(1)
	}

	p := parser.New(tokens)
	env := parser.NewDefaultEnvironment()
	if err := p.Execute(env); err != nil {
		var parserErr *parser.ParserError
		var runtimeErr *parser.RuntimeError
		if errors.As(err, &runtimeErr) {
			fmt.Fprintln(os.Stderr, runtimeErr.Format(path, string(input)))
		} else if errors.As(err, &parserErr) {
			fmt.Fprintln(os.Stderr, parserErr.Format(path))
		} else {
			fmt.Fprintf(os.Stderr, "Generic Error: %v\n", err)
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
		stmts, err := p.Parse()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, stmt := range stmts {
			if expr, ok := stmt.(parser.ExpressionStatement); ok {
				val, err := expr.ExecuteValue(env)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(val)
			} else if err := stmt.Execute(env); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func handleFormat(args []string) {
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "Not enough arguments to format")
		os.Exit(1)
	}

	path := args[2]
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s", err)
		os.Exit(1)
	}

	fmtSource, err := formatter.Format(string(source))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting source: %s", err)
		os.Exit(1)
	}
	fmt.Print(fmtSource)
}
