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

func (p *Parser) parseStatement() {
	if p.peek() == lexer.LetToken {
		p.parseDeclareStatement()
	} else if p.peek() == lexer.AssignToken {
		p.parseAssignStatement()
	} else if p.peek() == lexer.PrintToken {
		p.parsePrintStatement()
	} else {
		p.error()
	}
}

func (p *Parser) parseDeclareStatement() {
	p.match(lexer.LetToken)
	p.match(lexer.IdentToken)
	p.match(lexer.DeclareToken)
	p.parseExpression()
}

func (p *Parser) parseAssignStatement() {
	p.match(lexer.IdentToken)
	p.match(lexer.AssignToken)
	p.parseExpression()
}

func (p *Parser) parsePrintStatement() {
	p.match(lexer.PrintToken)
	p.parseExpression()
}

func (p *Parser) parseExpression() {
	p.parseTerm()
	for p.peek() == lexer.PlusToken || p.peek() == lexer.MinusToken {
		p.match(p.peek())
		p.parseTerm()
	}
}

func (p *Parser) parseTerm() {
	p.parseUnary()
	for p.peek() == lexer.MultiplyToken || p.peek() == lexer.DivideToken {
		p.match(p.peek())
		p.parseUnary()
	}
}

func (p *Parser) parseUnary() {
	if p.peek() == lexer.MinusToken {
		p.match(lexer.MinusToken)
		p.parseUnary()
	} else {
		p.parseFactor()
	}
}

func (p *Parser) parseFactor() {
	if p.peek() == lexer.LeftParenToken {
		p.match(lexer.LeftParenToken)
		p.parseExpression()
		p.match(lexer.RightParenToken)
	} else if p.peek() == lexer.NumberToken {
		p.match(lexer.NumberToken)
	} else {
		p.match(lexer.IdentToken)
	}
}

func (p *Parser) parseProgram() {
	for p.current < len(p.tokens) {
		if p.peek() == lexer.EOFToken {
			p.error()
		}
		p.parseStatement()
	}
}

func (p *Parser) Parse() {
	p.parseProgram()
}
