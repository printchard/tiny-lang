package lexer

type TokenType int

const (
	LetToken TokenType = iota
	PrintToken
	IdentToken
	NumberToken
	AssignToken
	DeclareToken
	PlusToken
	MinusToken
	MultiplyToken
	DivideToken
	LeftParenToken
	RightParenToken
	EOFToken
)

func (t TokenType) String() string {
	switch t {
	case LetToken:
		return "LET"
	case PrintToken:
		return "PRINT"
	case IdentToken:
		return "IDENT"
	case NumberToken:
		return "NUMBER"
	case AssignToken:
		return "="
	case DeclareToken:
		return ":="
	case PlusToken:
		return "+"
	case MinusToken:
		return "-"
	case MultiplyToken:
		return "*"
	case DivideToken:
		return "/"
	case LeftParenToken:
		return "("
	case RightParenToken:
		return ")"
	case EOFToken:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Type    TokenType
	Literal string
	Column  int
	Line    int
}

func NewToken(t TokenType, col, line int) Token {
	return Token{
		Type:    t,
		Literal: t.String(),
		Column:  col,
		Line:    line,
	}
}
