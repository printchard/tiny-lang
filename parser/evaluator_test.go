package parser

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/printchard/tiny-lang/lexer"
)

func executeSource(t *testing.T, source string) string {
	t.Helper()

	tokens, err := lexer.New(source).Tokenize()
	if err != nil {
		t.Fatalf("lexing: %v", err)
	}

	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer

	executeErr := New(tokens).Execute(NewDefaultEnvironment())
	writer.Close()
	os.Stdout = oldStdout

	var output bytes.Buffer
	if _, err := io.Copy(&output, reader); err != nil {
		t.Fatal(err)
	}
	reader.Close()

	if executeErr != nil {
		t.Fatalf("execution: %v", executeErr)
	}
	return output.String()
}

func evalExpression(t *testing.T, source string) Value {
	t.Helper()
	tokens, err := lexer.New(source).Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	statements, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	expression, ok := statements[0].(ExpressionStatement)
	if !ok {
		t.Fatalf("got %T, want ExpressionStatement", statements[0])
	}
	value, err := expression.ExecuteValue(NewDefaultEnvironment())
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func TestEvaluateExpressions(t *testing.T) {
	tests := []struct {
		source string
		want   string
	}{
		{"1 + 2 * 3", "7"},
		{"10 / 2 - 1", "4"},
		{"true && !false", "true"},
		{`"tiny" + " lang"`, "tiny lang"},
	}

	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			if got := evalExpression(t, test.source).String(); got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}

func TestExecutePrograms(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{
			name:   "arrays and builtins",
			source: "let values := [1, 2]\nprint(len(values))\nprint(push(values, 3))",
			want:   "2\n[1 2 3]\n",
		},
		{
			name:   "function call",
			source: "func add: left, right { return left + right }\nprint(add(2, 3))",
			want:   "5\n",
		},
		{
			name:   "control flow",
			source: "let value := 0\nwhile value < 3 { print(value) value = value + 1 }",
			want:   "0\n1\n2\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := executeSource(t, test.source); got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}
