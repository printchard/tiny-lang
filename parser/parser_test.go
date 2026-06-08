package parser

import (
	"testing"

	"github.com/printchard/tiny-lang/lexer"
)

func parseSource(t *testing.T, source string) []Statement {
	t.Helper()
	tokens, err := lexer.New(source).Tokenize()
	if err != nil {
		t.Fatalf("lexing %q: %v", source, err)
	}
	statements, err := New(tokens).Parse()
	if err != nil {
		t.Fatalf("parsing %q: %v", source, err)
	}
	return statements
}

func TestParseValidPrograms(t *testing.T) {
	tests := []string{
		"1 + 2 * 3",
		"let values := [1, false, void]",
		"values[0] = 2",
		"if value { print(value) } else if other { print(other) } else { return }",
		"while value { value = value - 1 }",
		"func add: left, right { return left + right }",
		"let add := func: left, right { return left + right }",
	}

	for _, source := range tests {
		t.Run(source, func(t *testing.T) {
			if len(parseSource(t, source)) == 0 {
				t.Fatal("expected at least one statement")
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []string{
		"1 = 2",
		"if value {",
		"let value 5",
		"func { return 1 }",
	}

	for _, source := range tests {
		t.Run(source, func(t *testing.T) {
			tokens, err := lexer.New(source).Tokenize()
			if err != nil {
				t.Fatal(err)
			}
			if _, err := New(tokens).Parse(); err == nil {
				t.Fatal("expected parser error")
			}
		})
	}
}

func TestParseExpressionStatement(t *testing.T) {
	statements := parseSource(t, "1 + 2 * 3")
	if _, ok := statements[0].(ExpressionStatement); !ok {
		t.Fatalf("got %T, want ExpressionStatement", statements[0])
	}
}
