package parser

import (
	"strconv"

	"github.com/printchard/tiny-lang/lexer"
)

type Parser struct {
	tokens  []lexer.Token
	current int
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) error() {
	if p.current < len(p.tokens) {
		panic("unexpected token: " + p.tokens[p.current].Literal + " at " + strconv.Itoa(p.current))
	} else {
		panic("unexpected end of input")
	}
}

func (p *Parser) peek() lexer.TokenType {
	if p.current >= len(p.tokens) {
		return lexer.TokenType(0)
	}
	return p.tokens[p.current].Type
}

func (p *Parser) match(expected lexer.TokenType) bool {
	if p.peek() == expected {
		p.current++
		return true
	}
	p.error()
	return false
}

func (p *Parser) parseStatement() Statement {
	if p.peek() == lexer.LetToken {
		return p.parseDeclareStatement()
	} else if p.peek() == lexer.IdentToken {
		return p.parseAssignStatement()
	} else if p.peek() == lexer.PrintToken {
		return p.parsePrintStatement()
	} else {
		p.error()
	}
	return nil
}

func (p *Parser) parseDeclareStatement() Statement {
	p.match(lexer.LetToken)
	identToken := p.tokens[p.current]
	p.match(lexer.IdentToken)
	p.match(lexer.DeclareToken)
	return &DeclarationStatement{
		Identifier: &Identifier{Name: identToken.Literal},
		Value:      p.parseExpression(),
	}
}

func (p *Parser) parseAssignStatement() Statement {
	identToken := p.tokens[p.current]
	p.match(lexer.IdentToken)
	p.match(lexer.AssignToken)
	return &AssignmentStatement{
		Identifier: &Identifier{Name: identToken.Literal},
		Value:      p.parseExpression(),
	}
}

func (p *Parser) parsePrintStatement() Statement {
	p.match(lexer.PrintToken)
	return &PrintStatement{
		Expression: p.parseExpression(),
	}
}

func (p *Parser) parseExpression() Expression {
	left := p.parseTerm()
	for p.peek() == lexer.PlusToken || p.peek() == lexer.MinusToken {
		op := p.peek()
		p.match(op)
		right := p.parseTerm()
		left = &BinaryExpression{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseTerm() Expression {
	left := p.parseUnary()
	for p.peek() == lexer.MultiplyToken || p.peek() == lexer.DivideToken {
		op := p.peek()
		p.match(op)
		right := p.parseUnary()
		left = &BinaryExpression{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseUnary() Expression {
	if p.peek() == lexer.MinusToken {
		p.match(lexer.MinusToken)
		return &UnaryExpression{
			Op:    lexer.MinusToken,
			Right: p.parseUnary(),
		}
	} else {
		return p.parseFactor()
	}
}

func (p *Parser) parseFactor() Expression {
	if p.peek() == lexer.LeftParenToken {
		p.match(lexer.LeftParenToken)
		expr := p.parseExpression()
		p.match(lexer.RightParenToken)
		return expr
	} else if p.peek() == lexer.NumberToken {
		p.match(lexer.NumberToken)
		value, err := strconv.ParseFloat(p.tokens[p.current-1].Literal, 64)
		if err != nil {
			p.error()
		}
		return &NumberLiteral{Value: value}
	} else {
		p.match(lexer.IdentToken)
		return &Identifier{Name: p.tokens[p.current-1].Literal}
	}
}

func (p *Parser) parseProgram() []Statement {
	statements := []Statement{}
	for p.current < len(p.tokens) {
		if p.peek() == lexer.EOFToken {
			p.error()
		}
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	return statements
}

func (p *Parser) Parse() []Statement {
	return p.parseProgram()
}

func (p *Parser) Execute() {
	env := &Environment{Variables: make(map[string]float64)}
	for _, stmt := range p.Parse() {
		stmt.Execute(env)
	}
}
