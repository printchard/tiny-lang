package lexer

import "testing"

func TestTokenize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []TokenType
	}{
		{
			name:  "declaration",
			input: "let value := 5",
			want:  []TokenType{LetToken, IdentToken, DeclareToken, NumberToken},
		},
		{
			name:  "operators",
			input: "a >= b && c != d || !e",
			want:  []TokenType{IdentToken, GEQToken, IdentToken, AndToken, IdentToken, NotEqualToken, IdentToken, OrToken, NotToken, IdentToken},
		},
		{
			name:  "collection",
			input: "[1, false, void]",
			want:  []TokenType{LeftBracketToken, NumberToken, CommaToken, FalseToken, CommaToken, VoidToken, RightBracketToken},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := New(test.input).Tokenize()
			if err != nil {
				t.Fatal(err)
			}
			if len(tokens) != len(test.want) {
				t.Fatalf("got %d tokens, want %d", len(tokens), len(test.want))
			}
			for index, token := range tokens {
				if token.Type != test.want[index] {
					t.Errorf("token %d: got %s, want %s", index, token.Type, test.want[index])
				}
			}
		})
	}
}

func TestTokenizeErrors(t *testing.T) {
	for _, input := range []string{"@", "&", "|", `"unterminated`} {
		t.Run(input, func(t *testing.T) {
			if _, err := New(input).Tokenize(); err == nil {
				t.Fatal("expected lexer error")
			}
		})
	}
}
