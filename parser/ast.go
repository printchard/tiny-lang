package parser

import (
	"fmt"

	"github.com/printchard/tiny-lang/lexer"
)

type Environment struct {
	Variables map[string]float64
}

type Node interface {
	String() string
}

type Expression interface {
	Node
	Eval(env *Environment) float64
}

type Statement interface {
	Node
	Execute(env *Environment)
}

type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) String() string {
	return fmt.Sprintf("%f", n.Value)
}

func (n *NumberLiteral) Eval(env *Environment) float64 {
	return n.Value
}

type Identifier struct {
	Name string
}

func (i *Identifier) String() string {
	return i.Name
}

func (i *Identifier) Eval(env *Environment) float64 {
	if value, ok := env.Variables[i.Name]; ok {
		return value
	}
	panic(fmt.Sprintf("undefined variable: %s", i.Name))
}

type BinaryExpression struct {
	Left  Expression
	Op    lexer.TokenType
	Right Expression
}

func (b *BinaryExpression) String() string {
	return fmt.Sprintf("%s %s %s", b.Left.String(), b.Op, b.Right.String())
}

func (b *BinaryExpression) Eval(env *Environment) float64 {
	left := b.Left.Eval(env)
	right := b.Right.Eval(env)

	switch b.Op {
	case lexer.PlusToken:
		return left + right
	case lexer.MinusToken:
		return left - right
	case lexer.MultiplyToken:
		return left * right
	case lexer.DivideToken:
		return left / right
	default:
		panic(fmt.Sprintf("unknown operator: %s", b.Op))
	}
}

type UnaryExpression struct {
	Op    lexer.TokenType
	Right Expression
}

func (u *UnaryExpression) String() string {
	return fmt.Sprintf("%s %s", u.Op, u.Right.String())
}

func (u *UnaryExpression) Eval(env *Environment) float64 {
	switch u.Op {
	case lexer.MinusToken:
		return -u.Right.Eval(env)
	default:
		panic(fmt.Sprintf("unknown operator: %s", u.Op))
	}
}

type DeclarationStatement struct {
	Identifier *Identifier
	Value      Expression
}

func (d *DeclarationStatement) String() string {
	return fmt.Sprintf("let %s = %s", d.Identifier.String(), d.Value.String())
}
func (d *DeclarationStatement) Execute(env *Environment) {
	if _, ok := env.Variables[d.Identifier.Name]; ok {
		panic(fmt.Sprintf("variable already declared: %s", d.Identifier.Name))
	}
	env.Variables[d.Identifier.Name] = d.Value.Eval(env)
}

type AssignmentStatement struct {
	Identifier *Identifier
	Value      Expression
}

func (a *AssignmentStatement) String() string {
	return fmt.Sprintf("%s = %s", a.Identifier.String(), a.Value.String())
}

func (a *AssignmentStatement) Execute(env *Environment) {
	if _, ok := env.Variables[a.Identifier.Name]; !ok {
		panic(fmt.Sprintf("undefined variable: %s", a.Identifier.Name))
	}
	env.Variables[a.Identifier.Name] = a.Value.Eval(env)
}

type PrintStatement struct {
	Expression Expression
}

func (p *PrintStatement) String() string {
	return fmt.Sprintf("print %s", p.Expression.String())
}

func (p *PrintStatement) Execute(env *Environment) {
	fmt.Println(p.Expression.Eval(env))
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var result string
	for _, stmt := range p.Statements {
		result += stmt.String() + "\n"
	}
	return result
}

func (p *Program) Execute(env *Environment) {
	for _, stmt := range p.Statements {
		stmt.Execute(env)
	}
}
