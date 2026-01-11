package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/printchard/tiny-lang/lexer"
)

type Node interface {
	String() string
	GetToken() lexer.Token
}

type Expression interface {
	Node
	Eval(env *Environment) (Value, error)
}

type Statement interface {
	Node
	Execute(env *Environment) error
}

type ReturnSignal struct {
	Value
}

func (s *ReturnSignal) Error() string {
	return "return signal"
}

type NumberLiteral struct {
	Value float64
	lexer.Token
}

func (n *NumberLiteral) GetToken() lexer.Token {
	return n.Token
}

func (n *NumberLiteral) String() string {
	return fmt.Sprintf("%f", n.Value)
}

func (n *NumberLiteral) Eval(env *Environment) (Value, error) {
	return Value{Type: Number, Number: n.Value}, nil
}

type StringLiteral struct {
	Value string
	lexer.Token
}

func (s *StringLiteral) GetToken() lexer.Token {
	return s.Token
}

func (s *StringLiteral) String() string {
	return fmt.Sprintf("%q", s.Value)
}

func (s *StringLiteral) Eval(env *Environment) (Value, error) {
	return Value{Type: String, Str: s.Value}, nil
}

type BooleanLiteral struct {
	Value bool
	lexer.Token
}

func (b *BooleanLiteral) GetToken() lexer.Token {
	return b.Token
}

func (b *BooleanLiteral) String() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *BooleanLiteral) Eval(env *Environment) (Value, error) {
	return Value{Type: Boolean, Boolean: b.Value}, nil
}

type ArrayLiteral struct {
	Elements []Expression
	lexer.Token
}

func (a *ArrayLiteral) GetToken() lexer.Token {
	return a.Token
}

func (a *ArrayLiteral) String() string {
	var elements []string
	for _, elem := range a.Elements {
		elements = append(elements, elem.String())
	}
	fmt.Printf("ArrayLiteral: %s\n", strings.Join(elements, ", "))
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func (a *ArrayLiteral) Eval(env *Environment) (Value, error) {
	var values []Value
	for _, elem := range a.Elements {
		value, err := elem.Eval(env)
		if err != nil {
			return Value{}, err
		}
		values = append(values, value)
	}
	return Value{Type: Array, Array: values}, nil
}

type VoidLiteral lexer.Token

func (v VoidLiteral) Eval(env *Environment) (Value, error) {
	return Value{}, nil
}

func (v VoidLiteral) GetToken() lexer.Token {
	return lexer.Token(v)
}

func (v VoidLiteral) String() string {
	return "void"
}

type Identifier struct {
	Token lexer.Token
}

func (i *Identifier) GetToken() lexer.Token {
	return i.Token
}

func (i *Identifier) String() string {
	return i.Token.Literal
}

func (i *Identifier) Eval(env *Environment) (Value, error) {
	if value, ok := env.Get(i.String()); ok {
		return value, nil
	}
	return Value{}, NewRuntimeError(i, fmt.Sprintf("undefined variable: %s", i.String()))
}

type BinaryExpression struct {
	Left    Expression
	Op      lexer.TokenType
	Right   Expression
	OpToken lexer.Token
}

func (b *BinaryExpression) GetToken() lexer.Token {
	return b.OpToken
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
		return Value{}, NewRuntimeError(b, fmt.Sprintf("type mismatch: %s and %s", left.Type, right.Type))
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
				return Value{}, NewRuntimeError(b, "division by zero")
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
			return Value{}, NewRuntimeError(b, fmt.Sprintf("unknown operator: %s", b.Op))
		}
	case String:
		switch b.Op {
		case lexer.PlusToken:
			return Value{Type: String, Str: left.Str + right.Str}, nil
		case lexer.EqualToken:
			return Value{Type: Boolean, Boolean: left.Str == right.Str}, nil
		case lexer.NotEqualToken:
			return Value{Type: Boolean, Boolean: left.Str != right.Str}, nil
		default:
			return Value{}, NewRuntimeError(b, fmt.Sprintf("unknown operator for strings: %s", b.Op))
		}
	default:
		bLeft, bRight := left.AsBoolean(), right.AsBoolean()
		switch b.Op {
		case lexer.EqualToken:
			return Value{Type: Boolean, Boolean: bLeft == bRight}, nil
		case lexer.NotEqualToken:
			return Value{Type: Boolean, Boolean: bLeft != bRight}, nil
		case lexer.AndToken:
			return Value{Type: Boolean, Boolean: bLeft && bRight}, nil
		case lexer.OrToken:
			return Value{Type: Boolean, Boolean: bLeft || bRight}, nil
		default:
			return Value{}, NewRuntimeError(b, fmt.Sprintf("unsupported types for binary operations: %s, %s", left.Type, right.Type))
		}
	}
}

type UnaryExpression struct {
	Op      lexer.TokenType
	Right   Expression
	OpToken lexer.Token
}

func (u *UnaryExpression) GetToken() lexer.Token {
	return u.OpToken
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
			return Value{}, NewRuntimeError(u, fmt.Sprintf("unknown unary operator: %s", u.Op))
		}
	default:
		switch u.Op {
		case lexer.NotToken:
			return Value{Type: Boolean, Boolean: value.AsBoolean()}, nil
		default:
			return Value{}, NewRuntimeError(u, fmt.Sprintf("unknown unary operator for boolean: %s", u.Op))
		}
	}
}

type PostfixExpression struct {
	Left         Expression
	Index        Expression
	BracketToken lexer.Token
}

func (p *PostfixExpression) GetToken() lexer.Token {
	return p.BracketToken
}

func (p *PostfixExpression) String() string {
	return fmt.Sprintf("%s[%s]", p.Left.String(), p.Index.String())
}

func (p *PostfixExpression) Eval(env *Environment) (Value, error) {
	left, err := p.Left.Eval(env)
	if err != nil {
		return Value{}, err
	}
	index, err := p.Index.Eval(env)
	if err != nil {
		return Value{}, err
	}
	if left.Type != Array {
		return Value{}, NewRuntimeError(p, fmt.Sprintf("left side of postfix expression must be an array, got %s", left.Type))
	}
	if index.Type != Number {
		return Value{}, NewRuntimeError(p, fmt.Sprintf("index must be a number, got %s", index.Type))
	}
	if int(index.Number) < 0 || int(index.Number) >= len(left.Array) {
		return Value{}, NewRuntimeError(p, fmt.Sprintf("index out of bounds: %d", int(index.Number)))
	}
	return left.Array[int(index.Number)], nil
}

type DeclarationStatement struct {
	Identifier *Identifier
	Value      Expression
	LetToken   lexer.Token
}

func (d *DeclarationStatement) GetToken() lexer.Token {
	return d.LetToken
}

func (d *DeclarationStatement) String() string {
	return fmt.Sprintf("let %s = %s", d.Identifier.String(), d.Value.String())
}

func (d *DeclarationStatement) Execute(env *Environment) error {
	if _, ok := env.Get(d.Identifier.String()); ok {
		return NewRuntimeError(d, fmt.Sprintf("variable already declared: %s", d.Identifier.String()))
	}
	value, err := d.Value.Eval(env)
	if err != nil {
		return err
	}

	env.Define(d.Identifier.String(), value)
	return nil
}

type AssignmentStatement struct {
	Identifier  *Identifier
	Value       Expression
	AssignToken lexer.Token
}

func (a *AssignmentStatement) GetToken() lexer.Token {
	return a.AssignToken
}

func (a *AssignmentStatement) String() string {
	return fmt.Sprintf("%s = %s", a.Identifier.String(), a.Value.String())
}

func (a *AssignmentStatement) Execute(env *Environment) error {
	if _, ok := env.Get(a.Identifier.String()); !ok {
		return NewRuntimeError(a, fmt.Sprintf("undefined variable: %s", a.Identifier.String()))
	}
	value, err := a.Value.Eval(env)
	if err != nil {
		return err
	}
	env.Set(a.Identifier.String(), value)
	return nil
}

type IndexAssignmentStatement struct {
	Left        *Identifier
	Index       Expression
	Value       Expression
	AssignToken lexer.Token
}

func (i *IndexAssignmentStatement) GetToken() lexer.Token {
	return i.AssignToken
}

func (i *IndexAssignmentStatement) String() string {
	return fmt.Sprintf("%s[%s] = %s", i.Left.String(), i.Index.String(), i.Value.String())
}

func (i *IndexAssignmentStatement) Execute(env *Environment) error {
	arr, ok := env.Get(i.Left.String())
	if !ok {
		return NewRuntimeError(i, fmt.Sprintf("undefined variable: %s", i.Left.String()))
	} else if arr.Type != Array {
		return NewRuntimeError(i, fmt.Sprintf("left side of index assignment must be an array, got %s", arr.Type))
	}
	index, err := i.Index.Eval(env)
	if err != nil {
		return err
	}
	value, err := i.Value.Eval(env)
	if err != nil {
		return err
	}
	if index.Type != Number {
		return NewRuntimeError(i, fmt.Sprintf("index must be a number, got %s", index.Type))
	}
	if int(index.Number) < 0 || int(index.Number) >= len(arr.Array) {
		return NewRuntimeError(i, fmt.Sprintf("index out of bounds: %d", int(index.Number)))
	}
	arr.Array[int(index.Number)] = value
	return nil
}

type IfStatement struct {
	Condition Expression
	Then      []Statement
	Else      []Statement
	IfToken   lexer.Token
}

func (i *IfStatement) GetToken() lexer.Token {
	return i.IfToken
}

func (i *IfStatement) String() string {
	var thenBody strings.Builder
	for _, stmt := range i.Then {
		thenBody.WriteString(stmt.String() + "\n")
	}
	var elseBody strings.Builder
	for _, stmt := range i.Else {
		elseBody.WriteString(stmt.String() + "\n")
	}
	return fmt.Sprintf("if %s {\n%s} else {\n%s}", i.Condition.String(), thenBody.String(), elseBody.String())
}

func (i *IfStatement) Execute(env *Environment) error {
	val, err := i.Condition.Eval(env)
	if err != nil {
		return err
	}

	childEnv := NewEnvironment(env)

	if val.AsBoolean() {
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
	Condition  Expression
	Body       []Statement
	WhileToken lexer.Token
}

func (w *WhileStatement) GetToken() lexer.Token {
	return w.WhileToken
}

func (w *WhileStatement) String() string {
	var body strings.Builder
	for _, stmt := range w.Body {
		body.WriteString(stmt.String() + "\n")
	}
	return fmt.Sprintf("while %s {\n%s}", w.Condition.String(), body.String())
}

func (w *WhileStatement) Execute(env *Environment) error {
	val, err := w.Condition.Eval(env)
	if err != nil {
		return err
	}

	for val.AsBoolean() {
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
	}
	return nil
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var result strings.Builder
	for _, stmt := range p.Statements {
		result.WriteString(stmt.String() + "\n")
	}
	return result.String()
}

func (p *Program) Execute(env *Environment) error {
	for _, stmt := range p.Statements {
		if err := stmt.Execute(env); err != nil {
			return err
		}
	}
	return nil
}

type ExpressionStatement struct {
	Expr Expression
}

func (e ExpressionStatement) GetToken() lexer.Token {
	return e.Expr.GetToken()
}

func (e ExpressionStatement) Execute(env *Environment) error {
	_, err := e.Expr.Eval(env)
	return err
}

func (e ExpressionStatement) String() string {
	return e.Expr.String()
}

type FunctionStatement struct {
	Name      *Identifier
	Args      []*Identifier
	Body      []Statement
	FuncToken lexer.Token
}

func (f FunctionStatement) GetToken() lexer.Token {
	return f.FuncToken
}

func (f FunctionStatement) Execute(env *Environment) error {
	var argNames []string
	for _, arg := range f.Args {
		argNames = append(argNames, arg.String())
	}
	funcVal := Func{ArgNames: argNames, Body: f.Body}
	env.Set(f.Name.String(), Value{Type: Function, Function: funcVal})
	return nil
}

func (f FunctionStatement) String() string {
	var str strings.Builder
	fmt.Fprintf(&str, "func %s(", f.Name)
	str.WriteString(") {\n")
	for _, stmt := range f.Body {
		str.WriteString("  " + stmt.String() + "\n")
	}
	str.WriteString("}")
	return str.String()
}

type FunctionCallExpression struct {
	Name      *Identifier
	Args      []Expression
	LeftParen lexer.Token
}

func (f FunctionCallExpression) GetToken() lexer.Token {
	return f.LeftParen
}

func (f FunctionCallExpression) Eval(env *Environment) (Value, error) {
	resolved, ok := env.Get(f.Name.String())
	if !ok {
		return Value{}, NewRuntimeError(f, fmt.Sprintf("undefined function: %s", f.Name))
	}

	if resolved.Type == NativeFunction {
		var args []Value
		for _, arg := range f.Args {
			v, err := arg.Eval(env)
			if err != nil {
				return Value{}, err
			}
			args = append(args, v)
		}
		nativeFn := resolved.NativeFunction
		return nativeFn(args)
	}

	if resolved.Type != Function {
		return Value{}, NewRuntimeError(f, fmt.Sprintf("function call to non-function type: %s", f.Name))
	}

	funcVal := resolved.Function
	if len(f.Args) > len(funcVal.ArgNames) {
		return Value{}, NewRuntimeError(f, fmt.Sprintf("too many arguments for function %s", f.Name))
	} else if len(f.Args) < len(funcVal.ArgNames) {
		return Value{}, NewRuntimeError(f, fmt.Sprintf("too few arguments for function %s", f.Name))
	}

	funcEnv := NewEnvironment(env)
	for i := 0; i < len(f.Args); i++ {
		val, err := f.Args[i].Eval(env)
		if err != nil {
			return Value{}, err
		}
		funcEnv.Define(funcVal.ArgNames[i], val)
	}
	for _, s := range funcVal.Body {
		err := s.Execute(funcEnv)
		var ret *ReturnSignal
		if errors.As(err, &ret) {
			return ret.Value, nil
		} else if err != nil {
			return Value{}, err
		}
	}
	return Value{}, nil
}

func (f FunctionCallExpression) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s(", f.Name)
	if len(f.Args) > 0 {
		fmt.Fprint(&b, f.Args[0].String())
	}
	for _, arg := range f.Args[1:] {
		fmt.Fprintf(&b, "%s, ", arg.String())
	}
	b.WriteByte(')')
	return b.String()
}

type ReturnStatement struct {
	Return      Expression
	ReturnToken lexer.Token
}

func (r *ReturnStatement) GetToken() lexer.Token {
	return r.ReturnToken
}

func (r ReturnStatement) Execute(env *Environment) error {
	if r.Return == nil {
		return &ReturnSignal{}
	}
	val, err := r.Return.Eval(env)
	if err != nil {
		return err
	}
	return &ReturnSignal{Value: val}
}

func (r ReturnStatement) String() string {
	return fmt.Sprintf("return %s", r.Return)
}

type RuntimeError struct {
	Msg string
	lexer.Token
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("[Line %d:%d]: %s", e.Token.Line, e.Token.Column, e.Msg)
}

func (e *RuntimeError) Format(fileName, source string) string {
	lines := strings.Split(source, "\n")
	if e.Line > len(lines) || e.Line < 1 {
		return fmt.Sprintf("[%s:%d:%d]: %s", fileName, e.Token.Line, e.Token.Column, e.Msg)
	}

	errorLine := lines[e.Line-1]

	var caretPadding strings.Builder
	for i := 2; i < e.Token.Column; i++ {
		if i < len(errorLine) && errorLine[i-1] == '\t' {
			caretPadding.WriteString("\t")
		} else {
			caretPadding.WriteString(" ")
		}
	}

	return fmt.Sprintf(
		"[%s:%d:%d]: %s\n    %s\n    %s^",
		fileName, e.Line, e.Column, e.Msg, errorLine, caretPadding.String(),
	)
}

func NewRuntimeError(n Node, msg string) error {
	return &RuntimeError{Msg: msg, Token: n.GetToken()}
}
