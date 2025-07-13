package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"noCap", true},
		{"cap", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 is 1", true},
		{"1 aint 1", false},
		{"1 is 2", false},
		{"1 aint 2", true},
		{"noCap is noCap", true},
		{"cap is cap", true},
		{"noCap is cap", false},
		{"noCap aint cap", true},
		{"cap aint noCap", true},
		{"(1 < 2) is noCap", true},
		{"(1 < 2) is cap", false},
		{"(1 > 2) is noCap", false},
		{"(1 > 2) is cap", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"nah noCap", false},
		{"nah cap", true},
		{"nah 5", false},
		{"nah nah noCap", true},
		{"nah nah cap", false},
		{"nah nah 5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"vibe (noCap) { 10 }", 10},
		{"vibe (cap) { 10 }", nil},
		{"vibe (1) { 10 }", 10},
		{"vibe (1 < 2) { 10 }", 10},
		{"vibe (1 > 2) { 10 }", nil},
		{"vibe (1 > 2) { 10 } nvm { 20 }", 20},
		{"vibe (1 < 2) { 10 } nvm { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"yeet 10;", 10},
		{"yeet 10; 9;", 10},
		{"yeet 2 * 5; 9;", 10},
		{"9; yeet 2 * 5; 9;", 10},
		{"vibe (10 > 1) { yeet 10; }", 10},
		{
			`
vibe (10 > 1) {
  vibe (10 > 1) {
    yeet 10;
  }

  yeet 1;
}
`,
			10,
		},
		{
			`
fr f = cook(x) {
  yeet x;
  x + 10;
};
f(10);`,
			10,
		},
		{
			`
fr f = cook(x) {
   fr result = x + 10;
   yeet result;
   yeet 10;
};
f(10);`,
			20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + noCap;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + noCap; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-noCap",
			"unknown operator: -BOOLEAN",
		},
		{
			"noCap + cap;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"noCap + cap + noCap + cap;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; noCap + cap; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			"vibe (10 > 1) { noCap + cap; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
vibe (10 > 1) {
  vibe (10 > 1) {
    yeet noCap + cap;
  }

  yeet 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`{"name": "Monkey"}[cook(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
		{
			`999[1]`,
			"index operator not supported: INTEGER",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"fr a = 5; a;", 5},
		{"fr a = 5 * 5; a;", 25},
		{"fr a = 5; fr b = a; b;", 5},
		{"fr a = 5; fr b = a; fr c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "cook(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"fr identity = cook(x) { x; }; identity(5);", 5},
		{"fr identity = cook(x) { yeet x; }; identity(5);", 5},
		{"fr double = cook(x) { x * 2; }; double(5);", 10},
		{"fr add = cook(x, y) { x + y; }; add(5, 5);", 10},
		{"fr add = cook(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"cook(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
fr first = 10;
fr second = 10;
fr third = 10;

fr ourFunction = cook(first) {
  fr second = 20;

  first + second + third;
};

ourFunction(20) + first + second;`

	testIntegerObject(t, testEval(input), 70)
}

func TestClosures(t *testing.T) {
	input := `
fr newAdder = cook(x) {
  cook(y) { x + y };
};

fr addTwo = newAdder(2);
addTwo(2);`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`puts("hello", "world!")`, nil},
		{`first([1, 2, 3])`, 1},
		{`first([])`, nil},
		{`first(1)`, "argument to `first` must be ARRAY, got INTEGER"},
		{`last([1, 2, 3])`, 3},
		{`last([])`, nil},
		{`last(1)`, "argument to `last` must be ARRAY, got INTEGER"},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, nil},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`, "argument to `push` must be ARRAY, got INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedElem))
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"fr i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"fr myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"fr myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"fr myArray = [1, 2, 3]; fr i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `fr two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		noCap: 5,
		cap: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`fr key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{noCap: 5}[noCap]`,
			5,
		},
		{
			`{cap: 5}[cap]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestAssignmentStatement(t *testing.T) {
	input := `
		fr count = 0;
		count = count + 1;
		count;
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("Eval didn't return Integer. got=%T (%+v)", evaluated, evaluated)
	}

	testIntegerObject(t, result, 1)
}

func TestAssigmentError(t *testing.T) {
	input := `
		count = count + 1;
	`

	evaluated := testEval(input)
	err, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("Eval didn't return Error. got=%T (%+v)", evaluated, evaluated)
	}

	expectedMessage := "identifier not found: count"
	if err.Message != expectedMessage {
		t.Fatalf("wrong error message. expected=%q, got=%q", expectedMessage, err.Message)
	}
}

func TestForStatement(t *testing.T) {
	input := `
		fr items = [1, 2, 3, 4];
		fr count = 0;
		stalk (i in items) {
			count = count + i;
		}
		count;
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("Eval didn't return Integer. got=%T (%+v)", evaluated, evaluated)
	}

	testIntegerObject(t, result, 10)
}

func TestForStatementWithBreak(t *testing.T) {
	input := `
		fr items = [1, 2, 3, 4];
		fr count = 0;
		stalk (i in items) {
			vibe (i is 3) { bounce; }
			count = count + i;
		}
		count;
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("Eval didn't return Integer. got=%T (%+v)", evaluated, evaluated)
	}

	testIntegerObject(t, result, 3) // 1 + 2, break before 3 is added
}

func TestForStatementWithContinue(t *testing.T) {
	input := `
		fr items = [1, 2, 3, 4];
		fr count = 0;
		stalk (i in items) {
			vibe (i is 2) { pass; }
			count = count + i;
		}
		count;
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("Eval didn't return Integer. got=%T (%+v)", evaluated, evaluated)
	}

	testIntegerObject(t, result, 8) // 1 + 3 + 4, skip 2
}

func TestForStatementWithRange(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		// Basic range test
		{
			`fr count = 0; stalk (i in range(1, 5)) { count = count + i; } count;`,
			15, // 1 + 2 + 3 + 4 + 5 = 15
		},
		// Single element range
		{
			`fr sum = 0; stalk (i in range(5, 5)) { sum = sum + i; } sum;`,
			5,
		},
		// Range starting from 0
		{
			`fr sum = 0; stalk (i in range(0, 3)) { sum = sum + i; } sum;`,
			6, // 0 + 1 + 2 + 3 = 6
		},
		// Using range with break
		{
			`fr sum = 0; stalk (i in range(1, 10)) { vibe (i > 3) { bounce; } sum = sum + i; } sum;`,
			6, // 1 + 2 + 3 = 6
		},
		// Using range with continue
		{
			`fr sum = 0; stalk (i in range(1, 5)) { vibe (i is 3) { pass; } sum = sum + i; } sum;`,
			12, // 1 + 2 + 4 + 5 = 12
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		result, ok := evaluated.(*object.Integer)
		if !ok {
			t.Fatalf("Eval didn't return Integer for input %q. got=%T (%+v)", tt.input, evaluated, evaluated)
		}
		testIntegerObject(t, result, tt.expected)
	}
}

func TestWhileStatement(t *testing.T) {
	input := `
		fr i = 0;
		fr count = 0;
 		onRepeat (i < 5) {
			count = count + i;
			i = i + 1;
		}
		count;
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("Eval didn't return Integer. got=%T (%+v)", evaluated, evaluated)
	}

	testIntegerObject(t, result, 10)
}

func TestWhileStatementWithBreak(t *testing.T) {
	input := `
		fr i = 0;
		fr count = 0;
 		onRepeat (i < 5) {
			vibe (i is 3) { bounce; }
			count = count + i;
			i = i + 1;
		}
		count;
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("Eval didn't return Integer. got=%T (%+v)", evaluated, evaluated)
	}

	testIntegerObject(t, result, 3) // 0 + 1 + 2, break before adding 3
}

func TestWhileStatementWithContinue(t *testing.T) {
	input := `
		fr i = 0;
		fr count = 0;
 		onRepeat (i < 5) {
			vibe (i is 2) { 
				i = i + 1; 
				pass; 
			}
			count = count + i;
			i = i + 1;
		}
		count;
	`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("Eval didn't return Integer. got=%T (%+v)", evaluated, evaluated)
	}

	testIntegerObject(t, result, 8) // 0 + 1 + 3 + 4, skip 2
}

func TestBreakAndContinueOutsideLoopScenarios(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    `bounce;`,
			expected: "break statement cannot be used outside of loop",
		},
		{
			input:    `pass;`,
			expected: "continue statement cannot be used outside of loop",
		},
		{
			input:    `vibe (noCap) { bounce; }`,
			expected: "break statement cannot be used outside of loop",
		},
		{
			input:    `vibe (noCap) { pass; }`,
			expected: "continue statement cannot be used outside of loop",
		},
		{
			input:    `fr f = cook() { bounce; }; f();`,
			expected: "break statement cannot be used outside of loop",
		},
		{
			input:    `fr f = cook() { pass; }; f();`,
			expected: "continue statement cannot be used outside of loop",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		err, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("Eval didn't return Error. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if err.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, err.Message)
		}
	}
}

func TestForStatementWithStringRange(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Basic string range - concatenate all characters
		{
			`fr result = ""; stalk (char in range("hello")) { result = result + char; } result;`,
			"hello",
		},
		// String range with single character
		{
			`fr result = ""; stalk (char in range("a")) { result = result + char; } result;`,
			"a",
		},
		// Empty string range
		{
			`fr result = ""; stalk (char in range("")) { result = result + char; } result;`,
			"",
		},
		// String range with special characters
		{
			`fr result = ""; stalk (char in range("a!@")) { result = result + char; } result;`,
			"a!@",
		},
		// Count characters in string using range
		{
			`fr count = 0; stalk (char in range("test")) { count = count + 1; } count;`,
			4,
		},
		// String range with break
		{
			`fr result = ""; stalk (char in range("hello")) { vibe (char is "l") { bounce; } result = result + char; } result;`,
			"he",
		},
		// String range with continue
		{
			`fr result = ""; stalk (char in range("hello")) { vibe (char is "l") { pass; } result = result + char; } result;`,
			"heo",
		},
		// Using string range to find specific character
		{
			`fr found = 0; stalk (char in range("monkey")) { vibe (char is "k") { found = 1; bounce; } } found;`,
			1,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Fatalf("Eval didn't return String for input %q. got=%T (%+v)", tt.input, evaluated, evaluated)
			}
			if str.Value != expected {
				t.Errorf("String has wrong value for input %q. got=%q, want=%q", tt.input, str.Value, expected)
			}
		case int:
			result, ok := evaluated.(*object.Integer)
			if !ok {
				t.Fatalf("Eval didn't return Integer for input %q. got=%T (%+v)", tt.input, evaluated, evaluated)
			}
			testIntegerObject(t, result, int64(expected))
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
