package parser

import (
	"fmt"
)

type Environment struct {
	variables map[string]Value
	parent    *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		variables: make(map[string]Value),
		parent:    parent,
	}
}

var defaultVars map[string]Value = map[string]Value{
	"print": {
		Type: NativeFunction,
		NativeFunction: func(vs []Value) (Value, error) {
			if len(vs) < 1 {
				return Value{}, fmt.Errorf("print expects at least 1 value")
			}

			fmt.Print(vs[0])
			for _, v := range vs[1:] {
				fmt.Printf(" %v", v)
			}
			fmt.Println()
			return Value{}, nil
		},
	},
}

func NewDefaultEnvironment() *Environment {
	return &Environment{
		variables: defaultVars,
	}
}

func (env *Environment) Set(name string, value Value) {
	if _, ok := env.variables[name]; ok {
		env.variables[name] = value
	} else {
		if env.parent != nil {
			env.parent.Set(name, value)
		} else {
			env.variables[name] = value
		}
	}
}

func (env *Environment) Define(name string, value Value) {
	env.variables[name] = value
}

func (env *Environment) Get(name string) (Value, bool) {
	value, ok := env.variables[name]
	if !ok && env.parent != nil {
		return env.parent.Get(name)
	}
	return value, ok
}

type ValueType int

const (
	Void ValueType = iota
	Number
	String
	Boolean
	Array
	Function
	NativeFunction
)

func (v ValueType) String() string {
	switch v {
	case Void:
		return "Void"
	case Number:
		return "Number"
	case String:
		return "String"
	case Boolean:
		return "Boolean"
	case Array:
		return "Array"
	case Function:
		return "Function"
	case NativeFunction:
		return "NativeFn"

	default:
		return "unknown"
	}
}

type Value struct {
	Type           ValueType
	Number         float64
	Str            string
	Boolean        bool
	Array          []Value
	Function       Func
	NativeFunction func([]Value) (Value, error)
}

func (v Value) String() string {
	switch v.Type {
	case Void:
		return "void"
	case Number:
		return fmt.Sprintf("%f", v.Number)
	case String:
		return v.Str
	case Boolean:
		return fmt.Sprintf("%t", v.Boolean)
	case Array:
		return fmt.Sprintf("%v", v.Array)
	case Function:
		return "fn"
	case NativeFunction:
		return "nativeFn"
	default:
		return "Unknown value type"
	}
}

func (v Value) AsBoolean() bool {
	switch v.Type {
	case Void:
		return false
	case Number:
		return v.Number != 0
	case Boolean:
		return v.Boolean
	case Array:
		return len(v.Array) > 0
	case String:
		return len(v.Str) > 0
	default:
		return true
	}
}

type Func struct {
	ArgNames []string
	Body     []Statement
}
