package object

import "fmt"

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

type Environment struct {
	store map[string]Object
	outer *Environment
	Logs  []string
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Update(name string, val Object) Object {
	_, ok := e.store[name]
	if ok {
		e.store[name] = val
		return val
	}

	if e.outer != nil {
		return e.outer.Update(name, val)
	}

	return &Error{Message: fmt.Sprintf("bruh, you can't just arbitrarily assign to: %q without defining it first 🙄", name)}
}

func (e *Environment) AddLogs(log string) {
	if e.outer != nil {
		e.outer.AddLogs(log)
	} else {
		e.Logs = append(e.Logs, log)
	}
}
