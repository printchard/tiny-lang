package formatter

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFormatterIdempotency(t *testing.T) {
	source, err := os.ReadFile("../test.tiny")
	if err != nil {
		t.Fatal(err)
	}

	fmtSource, err := Format(string(source))
	if err != nil {
		t.Fatal(err)
	}

	fmt2Source, err := Format(fmtSource)
	if err != nil {
		t.Fatal(err)
	}

	if fmtSource != fmt2Source {
		t.Error(cmp.Diff(fmtSource, fmt2Source))
	}
}

func TestFormatterPrecedence(t *testing.T) {
	tests := []struct {
		source string
		want   string
	}{
		{
			source: "1+2+3",
			want:   "1 + 2 + 3",
		},
		{
			source: "(1+ 2) * 3",
			want:   "(1 + 2) * 3",
		},
		{
			source: "1 + (2 * 3)",
			want:   "1 + 2 * 3",
		},
		{
			source: "a - (b-c)",
			want:   "a - (b - c)",
		},
		{
			source: "(a - b) - c",
			want:   "a - b - c",
		},
		{
			source: "!(a&&b)",
			want:   "!(a && b)",
		},
		{
			source: "!(a) && b",
			want:   "!a && b",
		},
	}

	for _, test := range tests {
		format, err := Format(test.source)
		if err != nil {
			t.Error(err)
		}

		format = strings.TrimSpace(format)
		if format != test.want {
			t.Error(cmp.Diff(test.want, format))
		}
	}

}
