package lexer

import (
	"strconv"
	"unicode"
)

type Lexer struct {
	input    string
	position int
	line     int
}

func New(input string) *Lexer {
	return &Lexer{
		input: input,
	}
}

func (l *Lexer) error() {
	if l.position >= len(l.input) {
		return
	}
	panic("unexpected character at line " + strconv.Itoa(l.line) + ": " + string(l.peek()))
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
	}
	l.position++
	return char
}

func (l *Lexer) skipWhitespace() {
	for l.position < len(l.input) && (l.peek() == ' ' || l.peek() == '\n') {
		if l.peek() == '\n' {
			l.line++
		}
		l.position++
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

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	if l.position >= len(l.input) {
		return Token{Type: EOFToken}
	}

	switch l.peek() {
	case '=':
		l.next()
		return NewToken(AssignToken)
	case ':':
		l.next()
		if l.peek() != '=' {
			l.error()
		}
		l.next()
		return NewToken(DeclareToken)
	case '+':
		l.next()
		return NewToken(PlusToken)
	case '-':
		l.next()
		return NewToken(MinusToken)
	case '*':
		l.next()
		return NewToken(MultiplyToken)
	case '/':
		l.next()
		return NewToken(DivideToken)
	case '(':
		l.next()
		return NewToken(LeftParenToken)
	case ')':
		l.next()
		return NewToken(RightParenToken)
	}

	if unicode.IsLetter(l.peek()) {
		literal := l.readLiteral()
		switch literal {
		case "let":
			return NewToken(LetToken)
		case "print":
			return NewToken(PrintToken)
		default:
			return Token{Type: IdentToken, Literal: literal}
		}
	}
	if unicode.IsDigit(l.peek()) {
		literal := l.readNumber()
		return Token{Type: NumberToken, Literal: literal}
	}

	l.error()
	return Token{}
}

func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for currToken := l.NextToken(); currToken.Type != EOFToken; currToken = l.NextToken() {
		tokens = append(tokens, currToken)
	}
	return tokens
}
