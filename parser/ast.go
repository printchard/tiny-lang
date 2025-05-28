package parser

import (
	"fmt"

	"github.com/printchard/tiny-lang/lexer"
)

type Node interface {
	String() string
}

type Expression interface {
	Node
	Eval(env *Environment) (Value, error)
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

func (n *NumberLiteral) Eval(env *Environment) (Value, error) {
	return Value{Type: Number, Number: n.Value}, nil
}

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) String() string {
	return fmt.Sprintf("%q", s.Value)
}

func (s *StringLiteral) Eval(env *Environment) (Value, error) {
	return Value{Type: String, String: s.Value}, nil
}

type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) String() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *BooleanLiteral) Eval(env *Environment) (Value, error) {
	return Value{Type: Boolean, Boolean: b.Value}, nil
}

type Identifier struct {
	Name string
}

func (i *Identifier) String() string {
	return i.Name
}

func (i *Identifier) Eval(env *Environment) (Value, error) {
	if value, ok := env.Get(i.Name); ok {
		return value, nil
	}
	return Value{}, fmt.Errorf("undefined variable: %s", i.Name)
}

type BinaryExpression struct {
	Left  Expression
	Op    lexer.TokenType
	Right Expression
}

func (b *BinaryExpression) String() string {
	return fmt.Sprintf("%s %s %s", b.Left.String(), b.Op, b.Right.String())
}

func (b *BinaryExpression) Eval(env *Environment) (Value, error) {
	left, err := b.Left.Eval(env)
	if err != nil {
		return Value{}, err
	}
	right, err := b.Right.Eval(env)
	if err != nil {
		return Value{}, err
	}
	if left.Type != right.Type {
		return Value{}, fmt.Errorf("type mismatch: %s and %s", left.Type, right.Type)
	}

	switch left.Type {
	case Number:
		switch b.Op {
		case lexer.PlusToken:
			return Value{Type: Number, Number: left.Number + right.Number}, nil
		case lexer.MinusToken:
			return Value{Type: Number, Number: left.Number - right.Number}, nil
		case lexer.MultiplyToken:
			return Value{Type: Number, Number: left.Number * right.Number}, nil
		case lexer.DivideToken:
			if right.Number == 0 {
				return Value{}, fmt.Errorf("division by zero")
			}
			return Value{Type: Number, Number: left.Number / right.Number}, nil
		case lexer.EqualToken:
			return Value{Type: Boolean, Boolean: left.Number == right.Number}, nil
		case lexer.NotEqualToken:
			return Value{Type: Boolean, Boolean: left.Number != right.Number}, nil
		case lexer.LTToken:
			return Value{Type: Boolean, Boolean: left.Number < right.Number}, nil
		case lexer.LEQToken:
			return Value{Type: Boolean, Boolean: left.Number <= right.Number}, nil
		case lexer.GTToken:
			return Value{Type: Boolean, Boolean: left.Number > right.Number}, nil
		case lexer.GEQToken:
			return Value{Type: Boolean, Boolean: left.Number >= right.Number}, nil
		default:
			return Value{}, fmt.Errorf("unknown operator: %s", b.Op)
		}
	case String:
		switch b.Op {
		case lexer.PlusToken:
			return Value{Type: String, String: left.String + right.String}, nil
		case lexer.EqualToken:
			return Value{Type: Boolean, Boolean: left.String == right.String}, nil
		case lexer.NotEqualToken:
			return Value{Type: Boolean, Boolean: left.String != right.String}, nil
		default:
			return Value{}, fmt.Errorf("unknown operator for strings: %s", b.Op)
		}
	case Boolean:
		switch b.Op {
		case lexer.EqualToken:
			return Value{Type: Boolean, Boolean: left.Boolean == right.Boolean}, nil
		case lexer.NotEqualToken:
			return Value{Type: Boolean, Boolean: left.Boolean != right.Boolean}, nil
		case lexer.AndToken:
			return Value{Type: Boolean, Boolean: left.Boolean && right.Boolean}, nil
		case lexer.OrToken:
			return Value{Type: Boolean, Boolean: left.Boolean || right.Boolean}, nil
		default:
			return Value{}, fmt.Errorf("unknown operator for booleans: %s", b.Op)
		}
	default:
		return Value{}, fmt.Errorf("unsupported type for binary operation: %s", left.Type)
	}
}

type UnaryExpression struct {
	Op    lexer.TokenType
	Right Expression
}

func (u *UnaryExpression) String() string {
	return fmt.Sprintf("%s %s", u.Op, u.Right.String())
}

func (u *UnaryExpression) Eval(env *Environment) (Value, error) {
	value, err := u.Right.Eval(env)
	if err != nil {
		return Value{}, err
	}
	switch value.Type {
	case Number:
		switch u.Op {
		case lexer.MinusToken:
			return Value{Type: Number, Number: -value.Number}, nil
		default:
			return Value{}, fmt.Errorf("unknown unary operator: %s", u.Op)
		}
	case Boolean:
		switch u.Op {
		case lexer.NotToken:
			return Value{Type: Boolean, Boolean: !value.Boolean}, nil
		default:
			return Value{}, fmt.Errorf("unknown unary operator for boolean: %s", u.Op)
		}
	default:
		return Value{}, fmt.Errorf("unsupported type for unary operation: %s", value.Type)
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
	if _, ok := env.Get(d.Identifier.Name); ok {
		return fmt.Errorf("variable already declared: %s", d.Identifier.Name)
	}
	value, err := d.Value.Eval(env)
	if err != nil {
		return err
	}

	env.Define(d.Identifier.Name, value)
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
	if _, ok := env.Get(a.Identifier.Name); !ok {
		return fmt.Errorf("undefined variable: %s", a.Identifier.Name)
	}
	value, err := a.Value.Eval(env)
	if err != nil {
		return err
	}
	env.Set(a.Identifier.Name, value)
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
	switch value.Type {
	case Number:
		fmt.Println(value.Number)
	case String:
		fmt.Println(value.String)
	case Boolean:
		fmt.Println(value.Boolean)
	}
	return nil
}

type IfStatement struct {
	Condition Expression
	Then      []Statement
	Else      []Statement
}

func (i *IfStatement) String() string {
	thenBody := ""
	for _, stmt := range i.Then {
		thenBody += stmt.String() + "\n"
	}
	elseBody := ""
	for _, stmt := range i.Else {
		elseBody += stmt.String() + "\n"
	}
	return fmt.Sprintf("if %s {\n%s} else {\n%s}", i.Condition.String(), thenBody, elseBody)
}

func (i *IfStatement) Execute(env *Environment) error {
	val, err := i.Condition.Eval(env)
	if err != nil {
		return err
	}

	if val.Type != Boolean {
		return fmt.Errorf("condition must evaluate to boolean, got %s", val.Type)
	}

	childEnv := NewEnvironment(env)

	if val.Boolean {
		for _, stmt := range i.Then {
			if err := stmt.Execute(childEnv); err != nil {
				return err
			}
		}
	} else {
		for _, stmt := range i.Else {
			if err := stmt.Execute(childEnv); err != nil {
				return err
			}
		}
	}
	return nil
}

type WhileStatement struct {
	Condition Expression
	Body      []Statement
}

func (w *WhileStatement) String() string {
	body := ""
	for _, stmt := range w.Body {
		body += stmt.String() + "\n"
	}
	return fmt.Sprintf("while %s {\n%s}", w.Condition.String(), body)
}

func (w *WhileStatement) Execute(env *Environment) error {
	val, err := w.Condition.Eval(env)
	if err != nil {
		return err
	}

	if val.Type != Boolean {
		return fmt.Errorf("condition must evaluate to boolean, got %s", val.Type)
	}

	for val.Boolean {
		childEnv := NewEnvironment(env)
		for _, stmt := range w.Body {
			if err := stmt.Execute(childEnv); err != nil {
				return err
			}
		}

		val, err = w.Condition.Eval(env)
		if err != nil {
			return err
		}

		if val.Type != Boolean {
			return fmt.Errorf("condition must evaluate to boolean, got %s", val.Type)
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
