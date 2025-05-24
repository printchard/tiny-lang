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
	Eval(env *Environment) (float64, error)
}

type Statement interface {
	Node
	Execute(env *Environment) error
}

type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) String() string {
	return fmt.Sprintf("%f", n.Value)
}

func (n *NumberLiteral) Eval(env *Environment) (float64, error) {
	return n.Value, nil
}

type Identifier struct {
	Name string
}

func (i *Identifier) String() string {
	return i.Name
}

func (i *Identifier) Eval(env *Environment) (float64, error) {
	if value, ok := env.Variables[i.Name]; ok {
		return value, nil
	}
	return 0, fmt.Errorf("undefined variable: %s", i.Name)
}

type BinaryExpression struct {
	Left  Expression
	Op    lexer.TokenType
	Right Expression
}

func (b *BinaryExpression) String() string {
	return fmt.Sprintf("%s %s %s", b.Left.String(), b.Op, b.Right.String())
}

func (b *BinaryExpression) Eval(env *Environment) (float64, error) {
	left, err := b.Left.Eval(env)
	if err != nil {
		return 0, err
	}
	right, err := b.Right.Eval(env)
	if err != nil {
		return 0, err
	}

	switch b.Op {
	case lexer.PlusToken:
		return left + right, nil
	case lexer.MinusToken:
		return left - right, nil
	case lexer.MultiplyToken:
		return left * right, nil
	case lexer.DivideToken:
		return left / right, nil
	case lexer.EqualToken:
		if left == right {
			return 1, nil
		}
		return 0, nil
	case lexer.NotEqualToken:
		if left != right {
			return 1, nil
		}
		return 0, nil
	case lexer.GTToken:
		if left > right {
			return 1, nil
		}
		return 0, nil
	case lexer.LTToken:
		if left < right {
			return 1, nil
		}
		return 0, nil
	case lexer.GEQToken:
		if left >= right {
			return 1, nil
		}
		return 0, nil
	case lexer.LEQToken:
		if left <= right {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", b.Op)
	}
}

type UnaryExpression struct {
	Op    lexer.TokenType
	Right Expression
}

func (u *UnaryExpression) String() string {
	return fmt.Sprintf("%s %s", u.Op, u.Right.String())
}

func (u *UnaryExpression) Eval(env *Environment) (float64, error) {
	switch u.Op {
	case lexer.MinusToken:
		value, err := u.Right.Eval(env)
		if err != nil {
			return 0, err
		}
		return -value, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", u.Op)
	}
}

type DeclarationStatement struct {
	Identifier *Identifier
	Value      Expression
}

func (d *DeclarationStatement) String() string {
	return fmt.Sprintf("let %s = %s", d.Identifier.String(), d.Value.String())
}
func (d *DeclarationStatement) Execute(env *Environment) error {
	if _, ok := env.Variables[d.Identifier.Name]; ok {
		return fmt.Errorf("variable already declared: %s", d.Identifier.Name)
	}
	value, err := d.Value.Eval(env)
	if err != nil {
		return err
	}

	env.Variables[d.Identifier.Name] = value
	return nil
}

type AssignmentStatement struct {
	Identifier *Identifier
	Value      Expression
}

func (a *AssignmentStatement) String() string {
	return fmt.Sprintf("%s = %s", a.Identifier.String(), a.Value.String())
}

func (a *AssignmentStatement) Execute(env *Environment) error {
	if _, ok := env.Variables[a.Identifier.Name]; !ok {
		return fmt.Errorf("undefined variable: %s", a.Identifier.Name)
	}
	value, err := a.Value.Eval(env)
	if err != nil {
		return err
	}
	env.Variables[a.Identifier.Name] = value
	return nil
}

type PrintStatement struct {
	Expression Expression
}

func (p *PrintStatement) String() string {
	return fmt.Sprintf("print %s", p.Expression.String())
}

func (p *PrintStatement) Execute(env *Environment) error {
	value, err := p.Expression.Eval(env)
	if err != nil {
		return err
	}
	fmt.Println(value)
	return nil
}

type IfStatement struct {
	Condition Expression
	Then      []Statement
	Else      []Statement
}

func (i *IfStatement) String() string {
	return fmt.Sprintf("if %s {}", i.Condition.String())
}

func (i *IfStatement) Execute(env *Environment) error {
	val, err := i.Condition.Eval(env)
	if err != nil {
		return err
	}

	if val != 0 {
		for _, stmt := range i.Then {
			if err := stmt.Execute(env); err != nil {
				return err
			}
		}
	} else {
		for _, stmt := range i.Else {
			if err := stmt.Execute(env); err != nil {
				return err
			}
		}
	}
	return nil
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

func (p *Program) Execute(env *Environment) error {
	for _, stmt := range p.Statements {
		if err := stmt.Execute(env); err != nil {
			return err
		}
	}
	return nil
}
