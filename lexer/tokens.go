package lexer

type TokenType int

const (
	EOFToken TokenType = iota
	LetToken
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
	IfToken
	ElseToken
	LeftBraceToken
	RightBraceToken
	EqualToken
	NotEqualToken
	GTToken
	LTToken
	GEQToken
	LEQToken
	WhileToken
	TrueToken
	FalseToken
	OrToken
	AndToken
	NotToken
	StringToken
	LeftBracketToken
	RightBracketToken
	CommaToken
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
	case IfToken:
		return "IF"
	case ElseToken:
		return "ELSE"
	case LeftBraceToken:
		return "{"
	case RightBraceToken:
		return "}"
	case EqualToken:
		return "=="
	case NotEqualToken:
		return "!="
	case GTToken:
		return ">"
	case LTToken:
		return "<"
	case GEQToken:
		return ">="
	case LEQToken:
		return "<="
	case EOFToken:
		return "EOF"
	case WhileToken:
		return "WHILE"
	case TrueToken:
		return "TRUE"
	case FalseToken:
		return "FALSE"
	case OrToken:
		return "||"
	case AndToken:
		return "&&"
	case NotToken:
		return "!"
	case StringToken:
		return "STRING"
	case LeftBracketToken:
		return "["
	case RightBracketToken:
		return "]"
	case CommaToken:
		return ","
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
