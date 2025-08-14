# noCap

noCap is a programming language with with syntax inspired by GenZ slangs. This is a fun side project built upon the foundations I learned from the book [Writing an Interpreter in Go](https://interpreterbook.com/) by Thorsten Ball.

Try it out yourself at: https://nocap.prateeksurana.me


## Local Development
The project has two targets hence two cmds:
1. **wasm**: This compiles to wasm that is used by the web based editor. Use the command to generate the wasm file:

    ```sh
    GOOS=js GOARCH=wasm go build -o nocap.wasm cmd/wasm/main.go
    ```

2. **cli**: This is compiled to  native binaries and is used by the [nocap CLI](https://www.npmjs.com/package/nocap-cli). Use `go run cmd/cli/main.go` to run the CLI locally.
