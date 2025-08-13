package main

import (
	"encoding/json"
	"nocap/evaluator"
	"nocap/lexer"
	"nocap/object"
	"nocap/parser"
	"syscall/js"
)

func main() {
	js.Global().Set("executeNoCap", ExecuteNoCap())
	<-make(chan struct{}) // Block forever
}

func ExecuteNoCap() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) (o any) {
		output := func(result *string, errors []string, logs []string) string {
			var resultValue any
			if result != nil {
				resultValue = *result
			}

			x := map[string]any{
				"result": resultValue,
				"errors": errors,
				"logs":   logs,
			}

			jsonString, err := json.Marshal(x)
			if err != nil {
				return "Something went wrong"
			}

			return string(jsonString)
		}

		// Single panic recovery for the entire function
		defer func() {
			if r := recover(); r != nil {
				// Return error through output function
				o = output(nil, []string{"This is awkward... something went very wrong and it's not your fault 😬!"}, []string{})
			}
		}()

		if len(args) != 1 {
			return output(nil, []string{"Invalid number of arguments. Expected 1 argument."}, []string{})
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

			return output(nil, errors, []string{})
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			switch res := evaluated.(type) {
			case *object.Error:
				return output(nil, []string{res.Inspect()}, env.Logs)
			case *object.Null:
				return output(nil, []string{}, env.Logs)
			default:
				result := evaluated.Inspect()
				return output(&result, []string{}, env.Logs)
			}
		} else {
			return output(nil, []string{}, env.Logs)
		}
	})
}
