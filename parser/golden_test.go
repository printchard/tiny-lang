package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoldenPrograms(t *testing.T) {
	files, err := filepath.Glob("testdata/*.tiny")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no golden programs found")
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			source, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			expectedPath := strings.TrimSuffix(file, ".tiny") + ".expected"
			expected, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatal(err)
			}
			if got := executeSource(t, string(source)); got != string(expected) {
				t.Fatalf("output mismatch:\ngot  %q\nwant %q", got, expected)
			}
		})
	}
}
