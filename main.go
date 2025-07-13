package main

import (
	"fmt"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"syscall/js"
)

func main() {
	js.Global().Set("executeNoCap", ExecuteNoCap())
	<-make(chan struct{}) // Block forever
}

func ExecuteNoCap() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid number of arguments. Expected 1 argument."
		}

		input := args[0].String()
		l := lexer.New(input)
		p := parser.New(l)
		env := object.NewEnvironment()

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			errors := make([]string, 0, len(p.Errors()))
			for _, msg := range p.Errors() {
				errors = append(errors, msg)
			}
			return errors
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			return fmt.Sprintf("Evaluated result: %s", evaluated.Inspect())
		} else {
			return "Evaluation returned nil."
		}
	})
}
