package evaluator

import (
	"nocap/object"
)

var builtins = map[string]*object.Builtin{
	"count": &object.Builtin{Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("count needs 1 argument but you gave it %d 🥲", len(args))
		}

		switch arg := args[0].(type) {
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		case *object.Hash:
			return &object.Integer{Value: int64(len(arg.Pairs))}
		default:
			return newError("count can only be used with arrays, strings, or hashes, not %s 🙄", arg.Type())
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
				return newError("slide needs 2 arguments but you gave it %d 🥲", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("slide needs an array to work with, not %s - can't slide on that! 🛝", args[0].Type())
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
				return newError("spread needs 1 or 2 arguments but you gave it %d 🥲", len(args))
			}

			if len(args) == 1 {
				if args[0].Type() != object.STRING_OBJ {
					return newError("spread expected a string but got %s - can't spread that 🧈", args[0].Type())
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
				return newError("spread needs two whole numbers, not %s and %s - those two don't make a range", args[0].Type(), args[1].Type())
			}

			start := args[0].(*object.Integer).Value
			end := args[1].(*object.Integer).Value

			if start > end {
				return newError("spread(%d, %d)? That's backwards - start cannot be greater than the end", start, end)
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
