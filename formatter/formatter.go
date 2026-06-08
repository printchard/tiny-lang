package formatter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/printchard/tiny-lang/lexer"
	"github.com/printchard/tiny-lang/parser"
)

const (
	precedenceOr = iota + 1
	precedenceAnd
	precedenceComparison
	precedenceAdditive
	precedenceMultiplicative
	precedenceUnary
	precedencePostfix
	precedencePrimary
)

func binaryPrecedence(op lexer.TokenType) int {
	switch op {
	case lexer.OrToken:
		return precedenceOr
	case lexer.AndToken:
		return precedenceAnd
	case lexer.EqualToken, lexer.NotEqualToken,
		lexer.LTToken, lexer.LEQToken,
		lexer.GTToken, lexer.GEQToken:
		return precedenceComparison
	case lexer.PlusToken, lexer.MinusToken:
		return precedenceAdditive
	case lexer.MultiplyToken, lexer.DivideToken:
		return precedenceMultiplicative
	default:
		return precedencePrimary
	}
}

func expressionPrecedence(expr parser.Expression) int {
	switch expr := expr.(type) {
	case *parser.BinaryExpression:
		return binaryPrecedence(expr.Op)
	case *parser.UnaryExpression:
		return precedenceUnary
	case *parser.PostfixExpression:
		return precedencePostfix
	case parser.FunctionCallExpression:
		return precedencePostfix
	default:
		return precedencePrimary
	}
}

func needsParentheses(expr parser.Expression, parentPrecedence int, isRightChild bool) bool {
	childPrecedence := expressionPrecedence(expr)

	if childPrecedence < parentPrecedence {
		return true
	} else if childPrecedence > parentPrecedence {
		return false
	}

	return isRightChild
}

type Printer struct {
	builder strings.Builder
	indent  int
}

func (p *Printer) write(str string) {
	p.builder.WriteString(str)
}

func (p *Printer) writeIndent() {
	for range p.indent {
		p.builder.WriteByte(' ')
	}
}

func (p *Printer) Indent() {
	p.indent += 2
}

func (p *Printer) Unindent() {
	if p.indent > 1 {
		p.indent -= 2
	}
}

func (p *Printer) formatStatement(stmt parser.Statement) {
	p.writeIndent()

	switch stmt := stmt.(type) {
	case parser.ExpressionStatement:
		p.formatExpression(stmt.Expr, 0, false)
	case *parser.DeclarationStatement:
		p.write("let ")
		p.write(stmt.Identifier.String())
		p.write(" := ")
		p.formatExpression(stmt.Value, 0, false)
	case *parser.AssignmentStatement:
		p.write(stmt.Identifier.String())
		p.write(" = ")
		p.formatExpression(stmt.Value, 0, false)
	case *parser.IfStatement:
		p.formatIf(stmt)
	case *parser.WhileStatement:
		p.write("while ")
		p.formatExpression(stmt.Condition, 0, false)
		p.write(" {\n")
		p.Indent()
		p.formatStatements(stmt.Body, false)
		p.Unindent()
		p.writeIndent()
		p.write("}")
	case *parser.ReturnStatement:
		if stmt.Return == nil {
			p.write("return")
			return
		}
		p.write("return ")
		p.formatExpression(stmt.Return, 0, false)
	case *parser.IndexAssignmentStatement:
		p.write(stmt.Left.String())
		p.write("[")
		p.formatExpression(stmt.Index, 0, false)
		p.write("] = ")
		p.formatExpression(stmt.Value, 0, false)
	case *parser.FunctionDeclarationStatement:
		p.write("func ")
		p.write(stmt.Identifier.String())
		p.formatParameters(stmt.Fn.Params)
		p.formatFunctionBody(stmt.Fn.Body)
	default:
		panic("invalid format statement")
	}
}

func (p *Printer) formatIf(stmt *parser.IfStatement) {
	p.write("if ")
	p.formatExpression(stmt.Condition, 0, false)
	p.write(" {\n")
	p.Indent()
	p.formatStatements(stmt.Then, false)
	p.Unindent()
	p.writeIndent()
	p.write("}")

	if len(stmt.Else) == 0 {
		return
	}
	if len(stmt.Else) == 1 {
		if elseIf, ok := stmt.Else[0].(*parser.IfStatement); ok {
			p.write(" else ")
			p.formatIf(elseIf)
			return
		}
	}

	p.write(" else {\n")
	p.Indent()
	p.formatStatements(stmt.Else, false)
	p.Unindent()
	p.writeIndent()
	p.write("}")
}

func (p *Printer) formatStatements(statements []parser.Statement, topLevel bool) {
	for index, stmt := range statements {
		isFunction := isFunctionDeclaration(stmt)
		if topLevel && isFunction && index > 0 && !isFunctionDeclaration(statements[index-1]) {
			p.write("\n")
		}
		p.formatStatement(stmt)
		p.write("\n")
		if topLevel && isFunction && index < len(statements)-1 {
			p.write("\n")
		}
	}
}

func isFunctionDeclaration(stmt parser.Statement) bool {
	_, ok := stmt.(*parser.FunctionDeclarationStatement)
	return ok
}

func (p *Printer) formatParameters(params []*parser.Identifier) {
	if len(params) == 0 {
		return
	}
	p.write(": ")
	for index, param := range params {
		if index > 0 {
			p.write(", ")
		}
		p.write(param.String())
	}
}

func (p *Printer) formatFunctionBody(body []parser.Statement) {
	if len(body) == 0 {
		p.write(" {}")
		return
	}
	p.write(" {\n")
	p.Indent()
	p.formatStatements(body, false)
	p.Unindent()
	p.writeIndent()
	p.write("}")
}

func (p *Printer) formatExpression(expr parser.Expression, parentPrecedence int, isRightChild bool) {
	precedence := expressionPrecedence(expr)
	parenthesized := needsParentheses(expr, parentPrecedence, isRightChild)
	if parenthesized {
		p.write("(")
	}

	switch expr := expr.(type) {
	case *parser.NumberLiteral:
		p.write(strconv.FormatFloat(expr.Value, 'g', -1, 64))
	case *parser.StringLiteral:
		p.write(strconv.Quote(expr.Value))
	case *parser.BooleanLiteral:
		p.write(strconv.FormatBool(expr.Value))
	case *parser.Identifier:
		p.write(expr.Token.Literal)
	case parser.VoidLiteral:
		p.write("void")
	case *parser.ArrayLiteral:
		p.builder.WriteRune('[')
		first := true
		for _, elem := range expr.Elements {
			if first {
				first = false
			} else {
				p.write(", ")
			}
			p.formatExpression(elem, 0, false)
		}
		p.builder.WriteRune(']')
	case *parser.UnaryExpression:
		p.write(expr.Op.String())
		p.formatExpression(expr.Right, precedence, false)
	case *parser.BinaryExpression:
		p.formatExpression(expr.Left, precedence, false)
		p.write(" ")
		p.write(expr.Op.String())
		p.write(" ")
		p.formatExpression(expr.Right, precedence, true)
	case *parser.PostfixExpression:
		p.formatExpression(expr.Left, precedence, false)
		p.write("[")
		p.formatExpression(expr.Index, 0, false)
		p.write("]")
	case parser.FunctionCallExpression:
		p.write(expr.Name.String())
		p.write("(")
		first := true
		for _, arg := range expr.Args {
			if first {
				p.formatExpression(arg, 0, false)
				first = false
				continue
			}
			p.builder.WriteString(", ")
			p.formatExpression(arg, 0, false)
		}
		p.write(")")
	case *parser.FunctionLiteral:
		p.write("func")
		p.formatParameters(expr.Params)
		p.formatFunctionBody(expr.Body)

	default:
		panic(fmt.Sprintf("unsupported expression: %T", expr))
	}
	if parenthesized {
		p.write(")")
	}
}

func Format(source string) (string, error) {
	tokens, err := lexer.New(source).Tokenize()
	if err != nil {
		return "", err
	}

	stmts, err := parser.New(tokens).Parse()
	if err != nil {
		return "", err
	}

	var printer Printer
	printer.formatStatements(stmts, true)

	return printer.builder.String(), nil
}
