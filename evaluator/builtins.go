package evaluator

import (
	"nocap/object"
)

var builtins = map[string]*object.Builtin{
	"count": &object.Builtin{Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}

		switch arg := args[0].(type) {
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		case *object.Hash:
			return &object.Integer{Value: int64(len(arg.Pairs))}
		default:
			return newError("argument to `count` not supported, got %s",
				args[0].Type())
		}
	},
		Name: "count",
	},
	"caughtIn4K": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			logs := &object.Array{}

			logs.Elements = append(logs.Elements, args...)

			return logs
		},
		Name: "caughtIn4K",
	},
	"slide": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `slide` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
		Name: "slide",
	},
	"spread": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 2 || len(args) < 1 {
				return newError("wrong number of arguments. got=%d, want=2 or 1",
					len(args))
			}

			if len(args) == 1 {
				if args[0].Type() != object.STRING_OBJ {
					return newError("argument to `spread` must be STRING or INTEGER, got %s",
						args[0].Type())
				}

				// Split the string into characters
				str := args[0].(*object.String).Value
				elements := make([]object.Object, len(str))
				for i, char := range str {
					elements[i] = &object.String{Value: string(char)}
				}

				return &object.Array{Elements: elements}
			}

			if args[0].Type() != object.INTEGER_OBJ || args[1].Type() != object.INTEGER_OBJ {
				return newError("arguments to `spread` must be INTEGER, got %s and %s",
					args[0].Type(), args[1].Type())
			}

			start := args[0].(*object.Integer).Value
			end := args[1].(*object.Integer).Value

			if start > end {
				return newError("start value %d is greater than end value %d", start, end)
			}

			length := end - start + 1
			elements := make([]object.Object, length)
			for i := int64(0); i < length; i++ {
				elements[i] = &object.Integer{Value: start + i}
			}

			return &object.Array{Elements: elements}
		},
		Name: "spread",
	},
}
