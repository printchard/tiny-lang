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

type ParserError struct {
	lexer.Token
	Msg string
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("[Line :%d:%d]: %s", e.Token.Line, e.Token.Column, e.Msg)
}

func (e *ParserError) Format(fileName string) string {
	return fmt.Sprintf("[%s:%d:%d]: %s", fileName, e.Token.Line, e.Token.Column, e.Msg)
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) error(msg string) error {
	tok := p.peekToken()
	return &ParserError{Msg: msg, Token: tok}
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
	case lexer.IfToken:
		return p.parseIfStatement()
	case lexer.WhileToken:
		return p.parseWhileStatement()
	case lexer.FunctionToken:
		return p.parseFunctionStatement()
	case lexer.ReturnToken:
		return p.parseReturnStatement()
	case lexer.IdentToken:
		ident := p.peekToken()
		p.match(lexer.IdentToken)
		if p.peek() == lexer.AssignToken {
			return p.parseAssignStatement(ident)
		} else if p.peek() == lexer.LeftParenToken {
			expr, err := p.parseFunctionCall(ident)
			if err != nil {
				return nil, err
			}
			return ExpressionStatement{expr}, nil
		}
		fallthrough
	default:
		expr, err := p.parseLogicalExpression()
		if err != nil {
			return nil, err
		}

		return ExpressionStatement{expr}, nil
	}
}

func (p *Parser) parseDeclareStatement() (Statement, error) {
	letToken := p.peekToken()
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
		Identifier: &Identifier{identToken},
		Value:      exp,
		LetToken:   letToken,
	}, nil
}

func (p *Parser) parseAssignStatement(ident lexer.Token) (Statement, error) {
	if p.peek() == lexer.LeftBracketToken {
		if err := p.match(lexer.LeftBracketToken); err != nil {
			return nil, err
		}
		index, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.match(lexer.RightBracketToken); err != nil {
			return nil, err
		}
		assignToken := p.peekToken()
		if err := p.match(lexer.AssignToken); err != nil {
			return nil, err
		}
		exp, err := p.parseLogicalExpression()
		if err != nil {
			return nil, err
		}
		return &IndexAssignmentStatement{
			Left:        &Identifier{ident},
			Index:       index,
			Value:       exp,
			AssignToken: assignToken,
		}, nil
	}

	assignToken := p.peekToken()
	if err := p.match(lexer.AssignToken); err != nil {
		return nil, err
	}
	exp, err := p.parseLogicalExpression()
	if err != nil {
		return nil, err
	}
	return &AssignmentStatement{
		Identifier:  &Identifier{ident},
		Value:       exp,
		AssignToken: assignToken,
	}, nil
}

func (p *Parser) parseIfStatement() (Statement, error) {
	ifToken := p.peekToken()
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
			IfToken:   ifToken,
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
		IfToken:   ifToken,
	}, nil
}

func (p *Parser) parseWhileStatement() (Statement, error) {
	whileToken := p.peekToken()
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
		Condition:  cond,
		Body:       body,
		WhileToken: whileToken,
	}, nil
}

func (p *Parser) parseLogicalExpression() (Expression, error) {
	left, err := p.parseLogicalTerm()
	if err != nil {
		return nil, err
	}

	for p.peek() == lexer.OrToken {
		opToken := p.peekToken()
		if err := p.match(lexer.OrToken); err != nil {
			return nil, err
		}
		right, err := p.parseLogicalTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{
			Left:    left,
			Op:      lexer.OrToken,
			Right:   right,
			OpToken: opToken,
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
		opToken := p.peekToken()
		if err := p.match(lexer.AndToken); err != nil {
			return nil, err
		}
		right, err := p.parseLogicalUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{
			Left:    left,
			Op:      lexer.AndToken,
			Right:   right,
			OpToken: opToken,
		}
	}
	return left, nil
}

func (p *Parser) parseLogicalUnary() (Expression, error) {
	if p.peek() == lexer.NotToken {
		opToken := p.peekToken()
		if err := p.match(lexer.NotToken); err != nil {
			return nil, err
		}
		right, err := p.parseLogicalUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpression{
			Op:      lexer.NotToken,
			Right:   right,
			OpToken: opToken,
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
		opToken := p.peekToken()
		if err := p.match(op); err != nil {
			return nil, err
		}
		right, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return &BinaryExpression{
			Left:    left,
			Op:      op,
			Right:   right,
			OpToken: opToken,
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
		opToken := p.peekToken()
		if err := p.match(op); err != nil {
			return nil, err
		}
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{
			Left:    left,
			Op:      op,
			Right:   right,
			OpToken: opToken,
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
		opToken := p.peekToken()
		if err := p.match(op); err != nil {
			return nil, err
		}
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{
			Left:    left,
			Op:      op,
			Right:   right,
			OpToken: opToken,
		}
	}
	return left, nil
}

func (p *Parser) parseUnary() (Expression, error) {
	if p.peek() == lexer.MinusToken {
		opToken := p.peekToken()
		if err := p.match(lexer.MinusToken); err != nil {
			return nil, err
		}
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpression{
			Op:      lexer.MinusToken,
			Right:   right,
			OpToken: opToken,
		}, nil
	} else {
		return p.parseFactor()
	}
}

func (p *Parser) parseFactor() (Expression, error) {
	primary, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	if p.peek() == lexer.LeftBracketToken {
		bracketToken := p.peekToken()
		if err := p.match(lexer.LeftBracketToken); err != nil {
			return nil, err
		}
		index, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.match(lexer.RightBracketToken); err != nil {
			return nil, err
		}
		return &PostfixExpression{
			Left:         primary,
			Index:        index,
			BracketToken: bracketToken,
		}, nil
	}
	return primary, nil
}

func (p *Parser) parsePrimary() (Expression, error) {
	switch p.peek() {
	case lexer.LeftParenToken:
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
	case lexer.NumberToken:
		token := p.peekToken()
		if err := p.match(lexer.NumberToken); err != nil {
			return nil, err
		}
		value, err := strconv.ParseFloat(token.Literal, 64)
		if err != nil {
			return nil, err
		}
		return &NumberLiteral{Value: value, Token: token}, nil
	case lexer.StringToken:
		token := p.peekToken()
		if err := p.match(lexer.StringToken); err != nil {
			return nil, err
		}
		return &StringLiteral{Value: token.Literal, Token: token}, nil
	case lexer.TrueToken, lexer.FalseToken:
		token := p.peekToken()
		if err := p.match(token.Type); err != nil {
			return nil, err
		}
		return &BooleanLiteral{
			Value: token.Type == lexer.TrueToken,
			Token: token,
		}, nil
	case lexer.IdentToken:
		token := p.peekToken()
		p.match(lexer.IdentToken)
		if p.peek() == lexer.LeftParenToken {
			return p.parseFunctionCall(token)
		}
		return &Identifier{token}, nil
	case lexer.LeftBracketToken:
		return p.parseArrayLiteral()
	default:
		return nil, p.error(fmt.Sprintf("unexpected token in primary expression: %s", p.peek()))
	}
}

func (p *Parser) parseArrayLiteral() (Expression, error) {
	bracketToken := p.peekToken()
	if err := p.match(lexer.LeftBracketToken); err != nil {
		return nil, err
	}
	elements := []Expression{}
	for p.peek() != lexer.RightBracketToken {
		exp, err := p.parseLogicalExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, exp)
		if p.peek() == lexer.CommaToken {
			if err := p.match(lexer.CommaToken); err != nil {
				return nil, err
			}
		}
	}
	if err := p.match(lexer.RightBracketToken); err != nil {
		return nil, err
	}
	return &ArrayLiteral{Elements: elements, Token: bracketToken}, nil
}

func (p *Parser) parseFunctionStatement() (Statement, error) {
	var funcStmt FunctionStatement
	funcToken := p.peekToken()
	if err := p.match(lexer.FunctionToken); err != nil {
		return nil, err
	}
	ident := p.peekToken()
	if err := p.match(lexer.IdentToken); err != nil {
		return nil, err
	}
	funcStmt.Name = &Identifier{ident}
	funcStmt.FuncToken = funcToken
	if p.peek() == lexer.ColonToken {
		args, err := p.parseArgumentStatement()
		if err != nil {
			return nil, err
		}
		funcStmt.Args = args
	}

	if err := p.match(lexer.LeftBraceToken); err != nil {
		return nil, err
	}

	for p.peek() != lexer.RightBraceToken {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		funcStmt.Body = append(funcStmt.Body, stmt)
	}

	p.match(lexer.RightBraceToken)
	return funcStmt, nil
}

func (p *Parser) parseArgumentStatement() ([]*Identifier, error) {
	p.match(lexer.ColonToken)
	var decls []*Identifier
	ident := p.peekToken()
	if err := p.match(lexer.IdentToken); err != nil {
		return nil, err
	}
	decls = append(decls, &Identifier{ident})
	for p.peek() == lexer.CommaToken {
		p.match(lexer.CommaToken)
		ident := p.peekToken()
		if err := p.match(lexer.IdentToken); err != nil {
			return nil, err
		}
		decls = append(decls, &Identifier{ident})
	}

	return decls, nil
}

func (p *Parser) parseFunctionCall(ident lexer.Token) (Expression, error) {
	var fnCall FunctionCallExpression
	fnCall.Name = &Identifier{ident}
	leftParen := p.peekToken()
	if err := p.match(lexer.LeftParenToken); err != nil {
		return nil, err
	}
	fnCall.LeftParen = leftParen

	if p.peek() != lexer.RightParenToken {
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		fnCall.Args = append(fnCall.Args, arg)
		for p.peek() == lexer.CommaToken {
			p.match(lexer.CommaToken)
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			fnCall.Args = append(fnCall.Args, expr)
		}
	}

	if err := p.match(lexer.RightParenToken); err != nil {
		return nil, p.error("expected right closing paren.")
	}
	return fnCall, nil
}

func (p *Parser) parseReturnStatement() (Statement, error) {
	returnToken := p.peekToken()
	p.match(lexer.ReturnToken)

	if p.peek() == lexer.RightBraceToken || p.peek() == lexer.EOFToken {
		return &ReturnStatement{Return: nil, ReturnToken: returnToken}, nil
	}
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	return &ReturnStatement{Return: expr, ReturnToken: returnToken}, nil
}
