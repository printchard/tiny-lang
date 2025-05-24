package parser

type Environment struct {
	variables map[string]float64
	parent    *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		variables: make(map[string]float64),
		parent:    parent,
	}
}

func (env *Environment) Set(name string, value float64) {
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

func (env *Environment) Define(name string, value float64) {
	env.variables[name] = value
}

func (env *Environment) Get(name string) (float64, bool) {
	value, ok := env.variables[name]
	if !ok && env.parent != nil {
		return env.parent.Get(name)
	}
	return value, ok
}
