package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/printchard/tiny-lang/formatter"
	"github.com/printchard/tiny-lang/lexer"
	"github.com/printchard/tiny-lang/parser"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		repl()
		return
	}

	path := os.Args[1]

	if path == "fmt" {
		handleFormat(os.Args[2:])
		return
	}

	input, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
		return
	}
	lex := lexer.New(string(input))
	tokens, err := lex.Tokenize()
	if err != nil {
		var lexerErr *lexer.LexerError
		if errors.As(err, &lexerErr) {
			log.Fatalln(lexerErr.Format(path))
		} else {
			log.Fatalf("Generic Error: %v\n", err)
		}
	}

	p := parser.New(tokens)
	env := parser.NewDefaultEnvironment()
	if err := p.Execute(env); err != nil {
		var parserErr *parser.ParserError
		var runtimeErr *parser.RuntimeError
		if errors.As(err, &runtimeErr) {
			log.Fatalln(runtimeErr.Format(path, string(input)))
		} else if errors.As(err, &parserErr) {
			log.Fatalln(parserErr.Format(path))
		} else {
			log.Fatalf("Generic Error: %v\n", err)
		}
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
	fs := flag.NewFlagSet("fmt", flag.ExitOnError)
	write := fs.Bool("w", false, "Set this flag to replace the content of the file")

	err := fs.Parse(args)
	if err != nil {
		log.Fatalf("Invalid arguments: %s\n", err)
	}

	path := fs.Arg(0)
	if path == "" {
		log.Fatalf("A path is needed to format")
	}

	source, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	fmtSource, err := formatter.Format(string(source))
	if err != nil {
		log.Fatalf("Error formatting source: %s", err)
	}
	if *write {
		os.WriteFile(path, []byte(fmtSource), 0644)
	} else {
		fmt.Print(fmtSource)
	}
}
