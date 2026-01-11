package parser

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
	Unknown ValueType = iota
	Number
	String
	Boolean
	Array
	Function
)

type Value struct {
	Type     ValueType
	Number   float64
	String   string
	Boolean  bool
	Array    []Value
	Function Func
}

func (v ValueType) String() string {
	switch v {
	case Number:
		return "Number"
	case String:
		return "String"
	case Boolean:
		return "Boolean"
	default:
		return "unknown"
	}
}

type Func struct {
	ArgNames []string
	Body     []Statement
}
