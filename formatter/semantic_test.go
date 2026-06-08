package formatter

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/printchard/tiny-lang/lexer"
	"github.com/printchard/tiny-lang/parser"
)

func executeOutput(t *testing.T, source string) string {
	t.Helper()
	tokens, err := lexer.New(source).Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	executeErr := parser.New(tokens).Execute(parser.NewDefaultEnvironment())
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

func TestFormattingPreservesProgramOutput(t *testing.T) {
	source := `
func add: left, right {
  return left + right
}

let values := [1, 2, 3]
print(add(values[0], values[1] * values[2]))
`

	formatted, err := Format(source)
	if err != nil {
		t.Fatal(err)
	}

	originalOutput := executeOutput(t, source)
	formattedOutput := executeOutput(t, formatted)
	if originalOutput != formattedOutput {
		t.Fatalf("formatted program changed output: original %q, formatted %q", originalOutput, formattedOutput)
	}
}
