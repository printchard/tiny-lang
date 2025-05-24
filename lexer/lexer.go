package lexer

import (
	"fmt"
	"unicode"
)

type Lexer struct {
	input    string
	position int
	line     int
	column   int
}

func New(input string) *Lexer {
	return &Lexer{
		input:  input,
		line:   1,
		column: 1,
	}
}

func (l *Lexer) error(msg string) error {
	return fmt.Errorf("lexer error at line %d, column %d: %s",
		l.line, l.column, msg)
}

func (l *Lexer) peek() rune {
	if l.position >= len(l.input) {
		return 0
	}
	return rune(l.input[l.position])
}

func (l *Lexer) next() rune {
	if l.position >= len(l.input) {
		return 0
	}
	char := rune(l.input[l.position])
	if char == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	l.position++
	return char
}

func (l *Lexer) skipWhitespace() {
	for l.position < len(l.input) && unicode.IsSpace(l.peek()) {
		l.next()
	}
}

func (l *Lexer) readLiteral() string {
	l.skipWhitespace()
	start := l.position
	for unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek()) || l.peek() == '_' {
		l.next()
	}
	return string(l.input[start:l.position])
}

func (l *Lexer) readNumber() string {
	l.skipWhitespace()
	start := l.position
	for unicode.IsDigit(l.peek()) {
		l.next()
	}
	if l.peek() == '.' {
		l.next()
		for unicode.IsDigit(l.peek()) {
			l.next()
		}
	}
	return string(l.input[start:l.position])
}

func (l *Lexer) NextToken() (Token, error) {
	l.skipWhitespace()
	if l.position >= len(l.input) {
		return Token{Type: EOFToken}, nil
	}

	switch l.peek() {
	case '=':
		l.next()
		return NewToken(AssignToken, l.column, l.line), nil
	case ':':
		l.next()
		if l.peek() != '=' {
			return Token{}, l.error("expected '=' after ':'")
		}
		l.next()
		return NewToken(DeclareToken, l.column, l.line), nil
	case '+':
		l.next()
		return NewToken(PlusToken, l.column, l.line), nil
	case '-':
		l.next()
		return NewToken(MinusToken, l.column, l.line), nil
	case '*':
		l.next()
		return NewToken(MultiplyToken, l.column, l.line), nil
	case '/':
		l.next()
		return NewToken(DivideToken, l.column, l.line), nil
	case '(':
		l.next()
		return NewToken(LeftParenToken, l.column, l.line), nil
	case ')':
		l.next()
		return NewToken(RightParenToken, l.column, l.line), nil
	}

	if unicode.IsLetter(l.peek()) {
		literal := l.readLiteral()
		switch literal {
		case "let":
			return NewToken(LetToken, l.column, l.line), nil
		case "print":
			return NewToken(PrintToken, l.column, l.line), nil
		default:
			return Token{Type: IdentToken, Literal: literal, Column: l.column, Line: l.line}, nil
		}
	}
	if unicode.IsDigit(l.peek()) {
		literal := l.readNumber()
		return Token{Type: NumberToken, Literal: literal, Column: l.column, Line: l.line}, nil
	}

	return Token{}, l.error("unexpected character")
}

func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token
	for currToken, err := l.NextToken(); currToken.Type != EOFToken; currToken, err = l.NextToken() {
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, currToken)
	}
	return tokens, nil
}
