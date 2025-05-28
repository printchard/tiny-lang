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

func (p *Parser) Parse() ([]Statement, error) {
	return p.parseProgram()
}

func (p *Parser) Execute(env *Environment) error {
	if env == nil {
		env = NewEnvironment(nil)
	}

	stmts, err := p.Parse()
	if err != nil {
		return err
	}
	for _, stmt := range stmts {
		if err := stmt.Execute(env); err != nil {
			return err
		}
	}
	return nil
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

func (p *Parser) parseStatement() (Statement, error) {
	switch p.peek() {
	case lexer.LetToken:
		return p.parseDeclareStatement()
	case lexer.IdentToken:
		return p.parseAssignStatement()
	case lexer.PrintToken:
		return p.parsePrintStatement()
	case lexer.IfToken:
		return p.parseIfStatement()
	case lexer.WhileToken:
		return p.parseWhileStatement()
	default:
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

	exp, err := p.parseLogicalExpression()
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
	exp, err := p.parseLogicalExpression()
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
	exp, err := p.parseLogicalExpression()
	if err != nil {
		return nil, err
	}
	return &PrintStatement{
		Expression: exp,
	}, nil
}

func (p *Parser) parseIfStatement() (Statement, error) {
	if err := p.match(lexer.IfToken); err != nil {
		return nil, err
	}
	cond, err := p.parseLogicalExpression()
	if err != nil {
		return nil, err
	}
	if err := p.match(lexer.LeftBraceToken); err != nil {
		return nil, err
	}

	thenBlock := []Statement{}

	for p.peek() != lexer.RightBraceToken {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		thenBlock = append(thenBlock, stmt)
	}
	if err := p.match(lexer.RightBraceToken); err != nil {
		return nil, err
	}

	if p.peek() != lexer.ElseToken {
		return &IfStatement{
			Condition: cond,
			Then:      thenBlock,
		}, nil
	}

	elseBlock := []Statement{}
	if err := p.match(lexer.ElseToken); err != nil {
		return nil, err
	}
	if p.peek() == lexer.IfToken {
		elseIf, err := p.parseIfStatement()
		if err != nil {
			return nil, err
		}
		elseBlock = append(elseBlock, elseIf)
	} else {
		if err := p.match(lexer.LeftBraceToken); err != nil {
			return nil, err
		}
		for p.peek() != lexer.RightBraceToken {
			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			elseBlock = append(elseBlock, stmt)
		}
		if err := p.match(lexer.RightBraceToken); err != nil {
			return nil, err
		}
	}

	return &IfStatement{
		Condition: cond,
		Then:      thenBlock,
		Else:      elseBlock,
	}, nil
}

func (p *Parser) parseWhileStatement() (Statement, error) {
	if err := p.match(lexer.WhileToken); err != nil {
		return nil, err
	}

	cond, err := p.parseLogicalExpression()
	if err != nil {
		return nil, err
	}
	if err := p.match(lexer.LeftBraceToken); err != nil {
		return nil, err
	}
	body := []Statement{}
	for p.peek() != lexer.RightBraceToken {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}
	if err := p.match(lexer.RightBraceToken); err != nil {
		return nil, err
	}
	return &WhileStatement{
		Condition: cond,
		Body:      body,
	}, nil
}

func (p *Parser) parseLogicalExpression() (Expression, error) {
	left, err := p.parseLogicalTerm()
	if err != nil {
		return nil, err
	}

	for p.peek() == lexer.OrToken {
		if err := p.match(lexer.OrToken); err != nil {
			return nil, err
		}
		right, err := p.parseLogicalTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{
			Left:  left,
			Op:    lexer.OrToken,
			Right: right,
		}
	}
	return left, nil
}

func (p *Parser) parseLogicalTerm() (Expression, error) {
	left, err := p.parseLogicalUnary()
	if err != nil {
		return nil, err
	}

	for p.peek() == lexer.AndToken {
		if err := p.match(lexer.AndToken); err != nil {
			return nil, err
		}
		right, err := p.parseLogicalUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{
			Left:  left,
			Op:    lexer.AndToken,
			Right: right,
		}
	}
	return left, nil
}

func (p *Parser) parseLogicalUnary() (Expression, error) {
	if p.peek() == lexer.NotToken {
		if err := p.match(lexer.NotToken); err != nil {
			return nil, err
		}
		right, err := p.parseLogicalUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpression{
			Op:    lexer.NotToken,
			Right: right,
		}, nil
	} else {
		return p.parseLogicalFactor()
	}
}

func (p *Parser) parseLogicalFactor() (Expression, error) {
	if p.peek() == lexer.LeftParenToken {
		if err := p.match(lexer.LeftParenToken); err != nil {
			return nil, err
		}
		expr, err := p.parseLogicalExpression()
		if err != nil {
			return nil, err
		}
		if err := p.match(lexer.RightParenToken); err != nil {
			return nil, err
		}
		return expr, nil
	} else if p.peek() == lexer.TrueToken || p.peek() == lexer.FalseToken {
		token := p.peekToken()
		if err := p.match(token.Type); err != nil {
			return nil, err
		}
		return &BooleanLiteral{
			Value: token.Type == lexer.TrueToken,
		}, nil
	} else {
		return p.parseComparison()
	}
}

func (p *Parser) parseComparison() (Expression, error) {
	left, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	switch p.peek() {
	case lexer.EqualToken, lexer.NotEqualToken, lexer.GTToken, lexer.LTToken, lexer.GEQToken, lexer.LEQToken:
		op := p.peek()
		if err := p.match(op); err != nil {
			return nil, err
		}
		right, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return &BinaryExpression{
			Left:  left,
			Op:    op,
			Right: right,
		}, nil
	}

	return left, nil
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
	} else if p.peek() == lexer.StringToken {
		token := p.peekToken()
		if err := p.match(lexer.StringToken); err != nil {
			return nil, err
		}
		return &StringLiteral{Value: token.Literal}, nil
	} else {
		token := p.peekToken()
		if err := p.match(lexer.IdentToken); err != nil {
			return nil, err
		}
		return &Identifier{Name: token.Literal}, nil
	}
}
