package parser

import (
	"errors"
	"fmt"
	"maps"
	"strconv"
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
		NativeFunction: func(n Node, vs []Value) (Value, error) {
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
	"len": {
		Type: NativeFunction,
		NativeFunction: func(n Node, vs []Value) (Value, error) {
			if len(vs) < 1 {
				return Value{}, fmt.Errorf("len expects at least 1 value")
			} else if len(vs) > 1 {
				return Value{}, fmt.Errorf("len expects at most 1 value")
			}

			v := vs[0]
			switch v.Type {
			case Array:
				return Value{Type: Number, Number: float64(len(v.Array))}, nil
			case String:
				return Value{Type: Number, Number: float64(len(v.Str))}, nil
			default:
				return Value{}, fmt.Errorf("value of type %s has no len", v.Type)
			}
		},
	},
	"push": {
		Type: NativeFunction,
		NativeFunction: func(n Node, vs []Value) (Value, error) {
			if len(vs) < 2 {
				return Value{}, fmt.Errorf("push expects at least 2 values")
			}

			arr := vs[0]
			vs = vs[1:]

			if arr.Type != Array {
				return Value{}, fmt.Errorf("expected array, got: %s", arr.Type)
			}

			pushed := make([]Value, len(arr.Array)+len(vs))
			copy(pushed, arr.Array)
			copy(pushed[len(arr.Array):], vs)

			return Value{Type: Array, Array: pushed}, nil
		},
	},
	"map": {
		Type: NativeFunction,
		NativeFunction: func(n Node, vs []Value) (Value, error) {
			if len(vs) < 2 {
				return Value{}, fmt.Errorf("map expects at least 2 values")
			}

			arr := vs[0]
			fn := vs[1]

			if arr.Type != Array {
				return Value{}, fmt.Errorf("expected array, got: %s", arr.Type)
			}

			res := make([]Value, 0, len(arr.Array))
			for _, v := range arr.Array {
				mapped, err := callFunction(n, fn, []Value{v})
				if err != nil {
					return Value{}, err
				}
				res = append(res, mapped)
			}

			return Value{Type: Array, Array: res}, nil
		},
	},
}

func callFunction(n Node, fn Value, args []Value) (Value, error) {
	switch fn.Type {
	case NativeFunction:
		result, err := fn.NativeFunction(n, args)
		if err != nil {
			return Value{}, NewRuntimeError(n, err.Error())
		}
		return result, nil
	case Function:
		funcVal := fn.Function
		if len(args) != len(funcVal.ArgNames) {
			return Value{}, NewRuntimeError(n, fmt.Sprintf("expected %d arguments, got %d", len(funcVal.ArgNames), len(args)))
		}
		callEnv := NewEnvironment(funcVal.Env)
		for i, name := range funcVal.ArgNames {
			callEnv.Define(name, args[i])
		}
		for _, s := range funcVal.Body {
			err := s.Execute(callEnv)
			var ret *ReturnSignal
			if errors.As(err, &ret) {
				return ret.Value, nil
			} else if err != nil {
				return Value{}, err
			}
		}
		return Value{}, nil
	default:
		return Value{}, NewRuntimeError(n, fmt.Sprintf("cannot call value of type %s", fn.Type))
	}
}

func NewDefaultEnvironment() *Environment {
	variables := make(map[string]Value, len(defaultVars))
	maps.Copy(variables, defaultVars)
	return &Environment{
		variables: variables,
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
	NativeFunction func(Node, []Value) (Value, error)
}

func (v Value) String() string {
	switch v.Type {
	case Void:
		return "void"
	case Number:
		return strconv.FormatFloat(v.Number, 'f', -1, 64)
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
	Env      *Environment
}
