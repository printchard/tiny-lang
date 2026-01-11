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

type LexerError struct {
	Msg    string
	line   int
	column int
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("[Line %d:%d]: %s", e.line, e.column, e.Msg)
}

func (e *LexerError) Format(fileName string) string {
	return fmt.Sprintf("[%s:%d:%d]: %s", fileName, e.line, e.column, e.Msg)
}

func New(input string) *Lexer {
	return &Lexer{
		input:  input,
		line:   1,
		column: 1,
	}
}

func (l *Lexer) error(msg string) error {
	return &LexerError{Msg: msg, line: l.line, column: l.column}
}

func (l *Lexer) newToken(t TokenType) Token {
	return Token{
		Type:    t,
		Literal: t.String(),
		Column:  l.column,
		Line:    l.line,
	}
}

func (l *Lexer) newTokenLiteral(t TokenType, literal string) Token {
	return Token{
		Type:    t,
		Literal: literal,
		Column:  l.column,
		Line:    l.line,
	}
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
		if l.peek() == '=' {
			l.next()
			return l.newToken(EqualToken), nil
		}
		return l.newToken(AssignToken), nil
	case ':':
		l.next()
		if l.peek() != '=' {
			return l.newToken(ColonToken), nil
		}
		l.next()
		return l.newToken(DeclareToken), nil
	case '+':
		l.next()
		return l.newToken(PlusToken), nil
	case '-':
		l.next()
		return l.newToken(MinusToken), nil
	case '*':
		l.next()
		return l.newToken(MultiplyToken), nil
	case '/':
		l.next()
		return l.newToken(DivideToken), nil
	case '(':
		l.next()
		return l.newToken(LeftParenToken), nil
	case ')':
		l.next()
		return l.newToken(RightParenToken), nil
	case '{':
		l.next()
		return l.newToken(LeftBraceToken), nil
	case '}':
		l.next()
		return l.newToken(RightBraceToken), nil
	case '<':
		l.next()
		if l.peek() == '=' {
			l.next()
			return l.newToken(LEQToken), nil
		}
		return l.newToken(LTToken), nil
	case '>':
		l.next()
		if l.peek() == '=' {
			l.next()
			return l.newToken(GEQToken), nil
		}
		return l.newToken(GTToken), nil
	case '!':
		l.next()
		if l.peek() == '=' {
			l.next()
			return l.newToken(NotEqualToken), nil
		}
		return l.newToken(NotToken), nil
	case '&':
		l.next()
		if l.peek() != '&' {
			return Token{}, l.error("expected '&' after '&'")
		}
		l.next()
		return l.newToken(AndToken), nil
	case '|':
		l.next()
		if l.peek() != '|' {
			return Token{}, l.error("expected '|' after '|'")
		}
		l.next()
		return l.newToken(OrToken), nil
	case '"':
		l.next()
		start := l.position
		for l.peek() != '"' && l.peek() != 0 {
			l.next()
		}
		if l.peek() == 0 {
			return Token{}, l.error("unterminated string literal")
		}
		l.next()
		return Token{
			Type:    StringToken,
			Literal: string(l.input[start : l.position-1]),
			Column:  l.column,
			Line:    l.line,
		}, nil
	case '[':
		l.next()
		return l.newToken(LeftBracketToken), nil
	case ']':
		l.next()
		return l.newToken(RightBracketToken), nil
	case ',':
		l.next()
		return l.newToken(CommaToken), nil
	}

	if unicode.IsLetter(l.peek()) {
		literal := l.readLiteral()
		switch literal {
		case "let":
			return l.newToken(LetToken), nil
		case "if":
			return l.newToken(IfToken), nil
		case "else":
			return l.newToken(ElseToken), nil
		case "while":
			return l.newToken(WhileToken), nil
		case "true":
			return l.newToken(TrueToken), nil
		case "false":
			return l.newToken(FalseToken), nil
		case "func":
			return l.newToken(FunctionToken), nil
		case "return":
			return l.newToken(ReturnToken), nil
		case "void":
			return l.newToken(VoidToken), nil
		default:
			return l.newTokenLiteral(IdentToken, literal), nil
		}
	}
	if unicode.IsDigit(l.peek()) {
		literal := l.readNumber()
		return l.newTokenLiteral(NumberToken, literal), nil
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
