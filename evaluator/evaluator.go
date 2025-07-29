package evaluator

import (
	"fmt"
	"nocap/ast"
	"nocap/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.AssignmentStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		obj := env.Update(node.Name.Value, val)
		if isError(obj) {
			return obj
		}

	case *ast.IndexExpressionAssignmentStatement:
		item := Eval(node.Left.Left, env)
		if isError(item) {
			return item
		}

		if item.Type() != object.ARRAY_OBJ && item.Type() != object.HASH_OBJ {
			return newError("seriously what are you trying to do here? [] can't be used with items of type %s 🙄", item.Type())
		}

		index := Eval(node.Left.Index, env)
		if isError(index) {
			return index
		}

		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		return evalIndexExpressionAssignmentStatement(item, index, val)

	case *ast.ForStatement:
		return evalForStatement(node, env)

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.BreakStatement:
		return &object.Break{}

	case *ast.ContinueStatement:
		return &object.Continue{}

	case *ast.FunctionStatement:
		params := node.Parameters
		body := node.Body
		fn := &object.Function{Parameters: params, Env: env, Body: body}

		env.Set(node.Name.Value, fn)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.Null:
		return NULL

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args, env)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Break:
			return newError("hey! you can't just bounce outside of a loop 🫠")
		case *object.Continue:
			return newError("hey! you can't just pass outside of a loop 🫠")
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(
	block *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ || rt == object.BREAK_OBJ || rt == object.CONTINUE_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "nah":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("what the hell is this? %s%s 🐘🐧", operator, right.Type())
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case (left.Type() == object.INTEGER_OBJ || left.Type() == object.FLOAT_OBJ) && (right.Type() == object.INTEGER_OBJ || right.Type() == object.FLOAT_OBJ):
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "is":
		return nativeBoolToBooleanObject(left == right)
	case operator == "aint":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("what the hell is %s supposed to do between a %s and a %s 😬",
			operator, left.Type(), right.Type())
	default:
		return newError("idk how to %s a %s with a %s 😬",
			operator, left.Type(), right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch right.Type() {
	case object.INTEGER_OBJ:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	case object.FLOAT_OBJ:
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	default:
		return newError("idk how to: -%s 😬", right.Type())
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("my math teacher said no dividing by zero! 😤")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "is":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "aint":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("idk how to %s a %s with a %s 😬",
			operator, left.Type(), right.Type())
	}
}

func evalFloatInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	var leftVal, rightVal float64

	if left.Type() == object.INTEGER_OBJ {
		leftVal = float64(left.(*object.Integer).Value)
	} else {
		leftVal = left.(*object.Float).Value
	}

	if right.Type() == object.INTEGER_OBJ {
		rightVal = float64(right.(*object.Integer).Value)
	} else {
		rightVal = right.(*object.Float).Value
	}

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("my math teacher said no dividing by zero! 😤")
		}

		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "is":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "aint":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("idk how to %s a %s with a %s 😬",
			operator, left.Type(), right.Type())
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "is":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "aint":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("idk how to %s a %s with a %s 😬",
			operator, left.Type(), right.Type())
	}
}

func evalIfExpression(
	ie *ast.IfExpression,
	env *object.Environment,
) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	}

	// Check else-if conditions
	for _, elseIf := range ie.ElseIfs {
		condition := Eval(elseIf.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return Eval(elseIf.Consequence, env)
		}
	}

	// Fall back to alternative or NULL
	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return NULL
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("%s? Never heard of them 🤷‍♀️", node.Value)
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object, env *object.Environment) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return applyBuiltIn(fn, args, env)

	default:
		return newError("%s can't be cooked! 😭", fn.Type())
	}
}

func applyBuiltIn(fn *object.Builtin, args []object.Object, env *object.Environment) object.Object {
	if fn.Name == "caughtIn4K" {
		l := fn.Fn(args...)

		logs, ok := l.(*object.Array)
		if !ok {
			return l
		}

		for _, log := range logs.Elements {
			env.AddLogs(log.Inspect())
		}

		return NULL
	}

	return fn.Fn(args...)
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		if left.Type() == object.ARRAY_OBJ {
			return newError("hey you can only use [] with whole numbers, %s aint it", index.Type())
		}
		return newError("you can't use [] with %s 🤷‍♂️", left.Type())
	}
}

func evalIndexExpressionAssignmentStatement(item, index, value object.Object) object.Object {
	switch {
	case item.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		arr := item.(*object.Array)
		idx := index.(*object.Integer).Value
		max := int64(len(arr.Elements))

		if idx < 1 || idx > max {
			return newError("this array only goes from 1-%d, but you tried to grab %d - that's way off! 📏", max, idx)
		}

		arr.Elements[idx-1] = value
		return NULL

	case item.Type() == object.ARRAY_OBJ && index.Type() != object.INTEGER_OBJ:
		return newError("hey you can only use [] with whole numbers, %s aint it", index.Type())

	case item.Type() == object.HASH_OBJ:
		hashObject := item.(*object.Hash)
		key, ok := index.(object.Hashable)
		if !ok {
			return newError("%s cannot be used as a hash key - try something more primitive 🔑", index.Type())
		}

		hashObject.Pairs[key.HashKey()] = object.HashPair{Key: index, Value: value}
		return NULL

	default:
		return newError("you can't use [] with %s 🤷‍♂️", item.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements))

	if idx < 1 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx-1]
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("%s cannot be used as a hash key - try something more primitive 🔑", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("%s cannot be used as a hash key - try something more primitive 🔑", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalForStatement(node *ast.ForStatement, env *object.Environment) object.Object {
	items := Eval(node.Items, env)
	if isError(items) {
		return items
	}

	if items.Type() != object.ARRAY_OBJ && items.Type() != object.HASH_OBJ {
		return newError("%s can't be looped over - try something iterable like an array, string or a hash 🌀", items.Type())
	}

	var elements []object.Object

	if items.Type() == object.ARRAY_OBJ {
		elements = items.(*object.Array).Elements
	} else {
		hashObj := items.(*object.Hash)

		for _, pair := range hashObj.Pairs {
			elements = append(elements, pair.Key)
		}
	}

	var result object.Object = NULL
	for _, element := range elements {
		extendedEnv := extendForEnv(element, node.Key, env)
		stmtResult := evalBlockStatement(node.Body, extendedEnv)
		if stmtResult != nil {
			switch stmtResult := stmtResult.(type) {
			case *object.Error:
				return stmtResult
			case *object.ReturnValue:
				return stmtResult
			case *object.Break:
				return NULL // break: exit the loop
			case *object.Continue:
				continue // continue: skip to next iteration
			default:
				result = stmtResult
			}
		}
	}

	return result
}

func evalWhileStatement(node *ast.WhileStatement, env *object.Environment) object.Object {
	var result object.Object = NULL
	for isTruthy(Eval(node.Condition, env)) {
		stmtResult := evalBlockStatement(node.Body, env)
		if stmtResult != nil {
			switch stmtResult := stmtResult.(type) {
			case *object.Error:
				return stmtResult
			case *object.ReturnValue:
				return stmtResult
			case *object.Break:
				return NULL // break: exit the loop
			case *object.Continue:
				continue // continue: skip to next iteration
			default:
				result = stmtResult
			}
		}
	}

	return result
}

func extendForEnv(item object.Object, key *ast.Identifier, e *object.Environment) *object.Environment {
	env := object.NewEnclosedEnvironment(e)

	env.Set(key.Value, item)

	return env
}
