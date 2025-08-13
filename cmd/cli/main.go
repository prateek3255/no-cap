package main

import (
	"errors"
	"fmt"
	"os"

	"nocap/evaluator"
	"nocap/lexer"
	"nocap/object"
	"nocap/parser"

	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "run <file>",
	Short: "Run a NoCap script",
	Example: "Run a NoCap script from a file:\n" +
		"$  nocap run script.nocap\n",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("you need to provide exactly one file to run")
		}

		if args[0] == "" {
			return errors.New("the file name cannot be empty")
		}

		// Check if the file exists
		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			return errors.New("the specified file does not exist")
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) (er error) {
		content, err := os.ReadFile(args[0])
		if err != nil {
			return errors.New("failed to read the file: " + err.Error())
		}

		defer func() {
			if r := recover(); r != nil {
				er = errors.New("this is awkward... something went very wrong and it's not your fault 😬")
			}
		}()

		input := string(content)
		l := lexer.New(input)
		p := parser.New(l)
		env := object.NewEnvironment()

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			fmt.Printf("\033[31m\nError: %s\n\033[0m", p.Errors()[0])
		}

		evaluated := evaluator.Eval(program, env)

		if len(env.Logs) > 0 {
			fmt.Println()
			for _, log := range env.Logs {
				fmt.Printf("\033[34m%s\033[0m\n", log)
			}
		}

		if evaluated != nil {
			switch res := evaluated.(type) {
			case *object.Error:
				fmt.Printf("\033[31m\nError: %s\n\033[0m", res.Inspect())
			case *object.Null:
				return nil
			default:
				result := evaluated.Inspect()
				fmt.Println("\n\033[32m" + result + "\033[0m")
			}
		}

		return nil
	},
}

var rootCmd = &cobra.Command{
	Use:   "nocap",
	Short: "A programming language for GenZ",
}

func main() {
	rootCmd.AddCommand(executeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
