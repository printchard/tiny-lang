package parser

import (
	"fmt"
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

func (p *Parser) error(msg string) error {
	tok := p.peekToken()
	return fmt.Errorf("Parse error at line %d, column %d: %s",
		tok.Line, tok.Column, msg)
}

func (p *Parser) peek() lexer.TokenType {
	if p.current >= len(p.tokens) {
		return lexer.TokenType(0)
	}
	return p.tokens[p.current].Type
}

func (p *Parser) peekToken() lexer.Token {
	if p.current >= len(p.tokens) {
		return lexer.Token{}
	}
	return p.tokens[p.current]
}

func (p *Parser) match(expected lexer.TokenType) error {
	if p.current >= len(p.tokens) {
		return p.error("unexpected EOF")
	}
	if p.peek() == expected {
		p.current++
		return nil
	}
	return p.error(fmt.Sprintf("expected %s, found %s", expected, p.peek()))
}

func (p *Parser) parseStatement() (Statement, error) {
	if p.peek() == lexer.LetToken {
		return p.parseDeclareStatement()
	} else if p.peek() == lexer.IdentToken {
		return p.parseAssignStatement()
	} else if p.peek() == lexer.PrintToken {
		return p.parsePrintStatement()
	} else {
		return nil, p.error("unexpected token")
	}
}

func (p *Parser) parseDeclareStatement() (Statement, error) {
	if err := p.match(lexer.LetToken); err != nil {
		return nil, err
	}
	identToken := p.peekToken()
	if err := p.match(lexer.IdentToken); err != nil {
		return nil, err
	}
	if err := p.match(lexer.DeclareToken); err != nil {
		return nil, err
	}

	exp, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	return &DeclarationStatement{
		Identifier: &Identifier{Name: identToken.Literal},
		Value:      exp,
	}, nil
}

func (p *Parser) parseAssignStatement() (Statement, error) {
	identToken := p.peekToken()
	if err := p.match(lexer.IdentToken); err != nil {
		return nil, err
	}
	if err := p.match(lexer.AssignToken); err != nil {
		return nil, err
	}
	exp, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	return &AssignmentStatement{
		Identifier: &Identifier{Name: identToken.Literal},
		Value:      exp,
	}, nil
}

func (p *Parser) parsePrintStatement() (Statement, error) {
	if err := p.match(lexer.PrintToken); err != nil {
		return nil, err
	}
	exp, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	return &PrintStatement{
		Expression: exp,
	}, nil
}

func (p *Parser) parseExpression() (Expression, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for p.peek() == lexer.PlusToken || p.peek() == lexer.MinusToken {
		op := p.peek()
		if err := p.match(op); err != nil {
			return nil, err
		}
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}
	return left, nil
}

func (p *Parser) parseTerm() (Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for p.peek() == lexer.MultiplyToken || p.peek() == lexer.DivideToken {
		op := p.peek()
		if err := p.match(op); err != nil {
			return nil, err
		}
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}
	return left, nil
}

func (p *Parser) parseUnary() (Expression, error) {
	if p.peek() == lexer.MinusToken {
		if err := p.match(lexer.MinusToken); err != nil {
			return nil, err
		}
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpression{
			Op:    lexer.MinusToken,
			Right: right,
		}, nil
	} else {
		return p.parseFactor()
	}
}

func (p *Parser) parseFactor() (Expression, error) {
	if p.peek() == lexer.LeftParenToken {
		if err := p.match(lexer.LeftParenToken); err != nil {
			return nil, err
		}
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.match(lexer.RightParenToken); err != nil {
			return nil, err
		}
		return expr, nil
	} else if p.peek() == lexer.NumberToken {
		if err := p.match(lexer.NumberToken); err != nil {
			return nil, err
		}
		value, err := strconv.ParseFloat(p.tokens[p.current-1].Literal, 64)
		if err != nil {
			return nil, err
		}
		return &NumberLiteral{Value: value}, nil
	} else {
		if err := p.match(lexer.IdentToken); err != nil {
			return nil, err
		}
		return &Identifier{Name: p.tokens[p.current-1].Literal}, nil
	}
}

func (p *Parser) parseProgram() ([]Statement, error) {
	statements := []Statement{}
	for p.current < len(p.tokens) {
		if p.peek() == lexer.EOFToken {
			return nil, p.error("unexpected EOF")
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	return statements, nil
}

func (p *Parser) Parse() ([]Statement, error) {
	return p.parseProgram()
}

func (p *Parser) Execute(env *Environment) error {
	if env == nil {
		env = &Environment{Variables: make(map[string]float64)}
	}

	stmts, err := p.Parse()
	if err != nil {
		return err
	}
	for _, stmt := range stmts {
		stmt.Execute(env)
	}
	return nil
}
