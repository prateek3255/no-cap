package evaluator

import (
	"nocap/lexer"
	"nocap/object"
	"nocap/parser"
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

func TestEvalIntegerAndFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5 + 3.5", 8.5},
		{"1.5 * 2.5", 3.75},
		{"10 - 2.5", 7.5},
		{"2 * 3.0", 6.0},
		{"8 / 2.0", 4.0},
		{"-5 + 3.5", -1.5},
		{"-10 - 2.5", -12.5},
		{"2 * -3.0", -6.0},
		{"2 * (3.0 + 2)", 10.0},
		{"(5 + 10.0 * 2 + 15 / 3) * 2 + -10", 50.0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
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

func TestGhostedExpression(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"ghosted"},
		{"ghosted is ghosted"},
		{"fr x = ghosted; x"},
		{"vibe (ghosted) { 10 }"},
		{"vibe (nah ghosted) { 10 }"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		// For conditional expressions that would return nothing, we expect null
		if tt.input == "vibe (ghosted) { 10 }" {
			testNullObject(t, evaluated)
		} else if tt.input == "vibe (nah ghosted) { 10 }" {
			// nah ghosted should be true, so this should return 10
			testIntegerObject(t, evaluated, 10)
		} else if tt.input == "ghosted is ghosted" {
			// ghosted is ghosted should be true
			testBooleanObject(t, evaluated, true)
		} else {
			// All other cases should return null
			testNullObject(t, evaluated)
		}
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

func TestIfElseIfExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Basic else-if: first condition true
		{"vibe (noCap) { 10 } unless (cap) { 20 }", 10},

		// Basic else-if: first condition false, second condition true
		{"vibe (cap) { 10 } unless (noCap) { 20 }", 20},

		// Multiple else-ifs: first condition true
		{"vibe (noCap) { 10 } unless (cap) { 20 } unless (noCap) { 30 }", 10},

		// Multiple else-ifs: second condition true
		{"vibe (cap) { 10 } unless (noCap) { 20 } unless (cap) { 30 }", 20},

		// Multiple else-ifs: third condition true
		{"vibe (cap) { 10 } unless (cap) { 20 } unless (noCap) { 30 }", 30},

		// No conditions true, no else clause
		{"vibe (cap) { 10 } unless (cap) { 20 } unless (cap) { 30 }", nil},

		// No conditions true, with else clause
		{"vibe (cap) { 10 } unless (cap) { 20 } unless (cap) { 30 } nvm { 40 }", 40},

		// Complex conditions with arithmetic
		{"vibe (1 > 2) { 10 } unless (2 > 1) { 20 } unless (3 > 2) { 30 }", 20},

		// Nested arithmetic in conditions
		{"vibe (1 + 1 is 3) { 10 } unless (2 * 2 is 4) { 20 } unless (5 - 2 is 3) { 30 }", 20},

		// String comparisons
		{"vibe (\"hello\" is \"world\") { 10 } unless (\"foo\" is \"foo\") { 20 } nvm { 30 }", 20},

		// Variable evaluations
		{`
		fr x = 5;
		vibe (x < 3) { 10 } unless (x > 3) { 20 } unless (x is 5) { 30 } nvm { 40 }
		`, 20},

		// Complex nested example
		{`
		fr a = 2;
		fr b = 3;
		vibe (a > b) { 
			100 
		} unless (a + b is 5) { 
			200 
		} unless (a * b is 6) { 
			300 
		} nvm { 
			400 
		}
		`, 200},
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
			"what the hell is + supposed to do between a integer and a boolean 🐘🐧",
		},
		{
			"5 + noCap; 5;",
			"what the hell is + supposed to do between a integer and a boolean 🐘🐧",
		},
		{
			"-noCap",
			"idk how to: -boolean 😬",
		},
		{
			"noCap + cap;",
			"idk how to + a boolean with a boolean 😬",
		},
		{
			"noCap + cap + noCap + cap;",
			"idk how to + a boolean with a boolean 😬",
		},
		{
			"5; noCap + cap; 5",
			"idk how to + a boolean with a boolean 😬",
		},
		{
			`"Hello" - "World"`,
			"idk how to - a string with a string 😬",
		},
		{
			"vibe (10 > 1) { noCap + cap; }",
			"idk how to + a boolean with a boolean 😬",
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
			"idk how to + a boolean with a boolean 😬",
		},
		{
			"foobar",
			"foobar? never heard of them 🤷‍♀️",
		},
		{
			`{"name": "Monkey"}[cook(x) { x }];`,
			"function cannot be used as a hash key - try something more primitive 🔑",
		},
		{
			`999[1]`,
			"you can't use [] with integer 🤷‍♂️",
		},
		{
			`cook(x) { x }(1, 2)`,
			"expected 1 arguments, but got 2 - you sure you know what you're doing? 🤔",
		},
		{
			`[1, 2, 3][4]`,
			"this array only goes from 1-3, but you tried to grab 4 - that's way off! 📏",
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
		{"cook add(x, y) { x + y;}; add(5, 5);", 10},
		{"cook identity(x) { yeet x; }; identity(5);", 5},
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
		{`count("")`, 0},
		{`count("four")`, 4},
		{`count("hello world")`, 11},
		{`count(1)`, "count can only be used with arrays, strings, or hashes, not integer 🙄"},
		{`count("one", "two")`, "count needs 1 argument but you gave it 2 🥲"},
		{`count([1, 2, 3])`, 3},
		{`count([])`, 0},
		{`caughtIn4K("hello", "world!")`, nil},
		{`slide([], 1)`, []int{1}},
		{`slide(1, 1)`, "slide needs an array to work with, not integer - can't slide on that! 🛝"},
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
			"[1, 2, 3][1]",
			1,
		},
		{
			"[1, 2, 3][2]",
			2,
		},
		{
			"[1, 2, 3][3]",
			3,
		},
		{
			"fr i = 1; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			2,
		},
		{
			"fr myArray = [1, 2, 3]; myArray[3];",
			3,
		},
		{
			"fr myArray = [1, 2, 3]; myArray[1] + myArray[2] + myArray[3];",
			6,
		},
		{
			"fr myArray = [1, 2, 3]; fr i = myArray[1]; myArray[i]",
			1,
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

	expectedMessage := "what the hell is + supposed to do between a builtin function and a integer 🐘🐧"
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

func TestForStatementWithHash(t *testing.T) {
	input := `
		fr h = {
			"one": 1,
			"two": 2
		};
		fr count = 0;
		stalk (i in h) {
			count = count + h[i];
		}
		count;
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("Eval didn't return Integer. got=%T (%+v)", evaluated, evaluated)
	}

	testIntegerObject(t, result, 3)
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
			`fr count = 0; stalk (i in spread(1, 5)) { count = count + i; } count;`,
			15, // 1 + 2 + 3 + 4 + 5 = 15
		},
		// Single element range
		{
			`fr sum = 0; stalk (i in spread(5, 5)) { sum = sum + i; } sum;`,
			5,
		},
		// Range starting from 0
		{
			`fr sum = 0; stalk (i in spread(0, 3)) { sum = sum + i; } sum;`,
			6, // 0 + 1 + 2 + 3 = 6
		},
		// Using range with break
		{
			`fr sum = 0; stalk (i in spread(1, 10)) { vibe (i > 3) { bounce; } sum = sum + i; } sum;`,
			6, // 1 + 2 + 3 = 6
		},
		// Using range with continue
		{
			`fr sum = 0; stalk (i in spread(1, 5)) { vibe (i is 3) { pass; } sum = sum + i; } sum;`,
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
			expected: "hey! you can't just bounce outside of a loop 🫠",
		},
		{
			input:    `pass;`,
			expected: "hey! you can't just pass outside of a loop 🫠",
		},
		{
			input:    `vibe (noCap) { bounce; }`,
			expected: "hey! you can't just bounce outside of a loop 🫠",
		},
		{
			input:    `vibe (noCap) { pass; }`,
			expected: "hey! you can't just pass outside of a loop 🫠",
		},
		{
			input:    `fr f = cook() { bounce; }; f();`,
			expected: "hey! you can't just bounce outside of a loop 🫠",
		},
		{
			input:    `fr f = cook() { pass; }; f();`,
			expected: "hey! you can't just pass outside of a loop 🫠",
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
			`fr result = ""; stalk (char in spread("hello")) { result = result + char; } result;`,
			"hello",
		},
		// String range with single character
		{
			`fr result = ""; stalk (char in spread("a")) { result = result + char; } result;`,
			"a",
		},
		// Empty string range
		{
			`fr result = ""; stalk (char in spread("")) { result = result + char; } result;`,
			"",
		},
		// String range with special characters
		{
			`fr result = ""; stalk (char in spread("a!@")) { result = result + char; } result;`,
			"a!@",
		},
		// Count characters in string using range
		{
			`fr count = 0; stalk (char in spread("test")) { count = count + 1; } count;`,
			4,
		},
		// String range with break
		{
			`fr result = ""; stalk (char in spread("hello")) { vibe (char is "l") { bounce; } result = result + char; } result;`,
			"he",
		},
		// String range with continue
		{
			`fr result = ""; stalk (char in spread("hello")) { vibe (char is "l") { pass; } result = result + char; } result;`,
			"heo",
		},
		// Using string range to find specific character
		{
			`fr found = 0; stalk (char in spread("monkey")) { vibe (char is "k") { found = 1; bounce; } } found;`,
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

func TestIndexExpressionAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Basic array assignment
		{
			`fr arr = [1, 2, 3]; arr[1] = 10; arr[1];`,
			10,
		},
		{
			`fr arr = [1, 2, 3]; arr[2] = 20; arr[2];`,
			20,
		},
		{
			`fr arr = [1, 2, 3]; arr[3] = 30; arr[3];`,
			30,
		},
		// Array assignment with expressions
		{
			`fr arr = [1, 2, 3]; arr[1] = arr[2] + arr[3]; arr[1];`,
			5,
		},
		{
			`fr arr = [1, 2, 3]; fr idx = 2; arr[idx] = 100; arr[2];`,
			100,
		},
		// Basic hash assignment
		{
			`fr hash = {"a": 1, "b": 2}; hash["a"] = 10; hash["a"];`,
			10,
		},
		{
			`fr hash = {"a": 1, "b": 2}; hash["b"] = 20; hash["b"];`,
			20,
		},
		// Hash assignment with new keys
		{
			`fr hash = {"a": 1}; hash["c"] = 30; hash["c"];`,
			30,
		},
		// Hash assignment with different key types
		{
			`fr hash = {1: "one", 2: "two"}; hash[1] = "ONE"; hash[1];`,
			"ONE",
		},
		{
			`fr hash = {noCap: 1, cap: 2}; hash[noCap] = 100; hash[noCap];`,
			100,
		},
		// Multiple assignments
		{
			`fr arr = [1, 2, 3]; arr[1] = 10; arr[2] = 20; arr[3] = 30; arr[1] + arr[2] + arr[3];`,
			60,
		},
		{
			`fr hash = {}; hash["x"] = 1; hash["y"] = 2; hash["z"] = 3; hash["x"] + hash["y"] + hash["z"];`,
			6,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			str, ok := tt.expected.(string)
			if ok {
				testStringObject(t, evaluated, str)
			}
		}
	}
}

func TestComplexIndexExpressionAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Array inside hash
		{
			`fr data = {"nums": [1, 2, 3]}; data["nums"][1] = 100; data["nums"][1];`,
			100,
		},
		{
			`fr data = {"nums": [1, 2, 3]}; data["nums"][2] = data["nums"][1] + data["nums"][3]; data["nums"][2];`,
			4,
		},
		// Hash inside array
		{
			`fr data = [{"a": 1}, {"b": 2}]; data[1]["a"] = 100; data[1]["a"];`,
			100,
		},
		{
			`fr data = [{"a": 1}, {"b": 2}]; data[2]["c"] = 300; data[2]["c"];`,
			300,
		},
		// Nested arrays
		{
			`fr matrix = [[1, 2], [3, 4]]; matrix[1][2] = 20; matrix[1][2];`,
			20,
		},
		{
			`fr matrix = [[1, 2], [3, 4]]; matrix[2][1] = matrix[1][1] + matrix[1][2]; matrix[2][1];`,
			3,
		},
		// Nested hashes
		{
			`fr nested = {"outer": {"inner": 1}}; nested["outer"]["inner"] = 100; nested["outer"]["inner"];`,
			100,
		},
		{
			`fr nested = {"a": {"b": {"c": 1}}}; nested["a"]["b"]["c"] = 999; nested["a"]["b"]["c"];`,
			999,
		},
		// Mixed complex nesting
		{
			`fr complex = {"data": [{"values": [1, 2, 3]}]}; complex["data"][1]["values"][2] = 200; complex["data"][1]["values"][2];`,
			200,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		}
	}
}

func TestIndexExpressionAssignmentErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		// Array index out of bounds
		{
			`fr arr = [1, 2, 3]; arr[5] = 10;`,
			"this array only goes from 1-3, but you tried to grab 5 - that's way off! 📏",
		},
		{
			`fr arr = [1, 2, 3]; arr[-1] = 10;`,
			"this array only goes from 1-3, but you tried to grab -1 - that's way off! 📏",
		},
		// Invalid types for index assignment
		{
			`fr num = 42; num[0] = 10;`,
			"seriously what are you trying to do here? [] can't be used with items of type integer 🙄",
		},
		{
			`fr str = "hello"; str[0] = "H";`,
			"seriously what are you trying to do here? [] can't be used with items of type string 🙄",
		},
		{
			`fr fn = cook(x) { x }; fn[0] = 10;`,
			"seriously what are you trying to do here? [] can't be used with items of type function 🙄",
		},
		// Invalid array index type
		{
			`fr arr = [1, 2, 3]; arr["invalid"] = 10;`,
			"hey you can only use [] with whole numbers, string aint it",
		},
		// Invalid hash key type
		{
			`fr hash = {}; hash[cook(x) { x }] = 10;`,
			"function cannot be used as a hash key - try something more primitive 🔑",
		},
		{
			`fr hash = {}; hash[[1, 2, 3]] = 10;`,
			"array cannot be used as a hash key - try something more primitive 🔑",
		},
		// Assignment to non-existent nested structure
		{
			`fr arr = [1, 2, 3]; arr[1]["key"] = 10;`,
			"seriously what are you trying to do here? [] can't be used with items of type integer 🙄",
		},
		{
			`fr hash = {"a": 1}; hash["a"][1] = 10;`,
			"seriously what are you trying to do here? [] can't be used with items of type integer 🙄",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned for input %q. got=%T(%+v)",
				tt.input, evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message for input %q. expected=%q, got=%q",
				tt.input, tt.expectedMessage, errObj.Message)
		}
	}
}

func TestIndexExpressionAssignmentWithDifferentTypes(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Assigning different types to arrays
		{
			`fr arr = [1, 2, 3]; arr[1] = "string"; arr[1];`,
			"string",
		},
		{
			`fr arr = [1, 2, 3]; arr[2] = noCap; arr[2];`,
			true,
		},
		{
			`fr arr = [1, 2, 3]; arr[3] = [4, 5, 6]; count(arr[3]);`,
			3,
		},
		{
			`fr arr = [1, 2, 3]; arr[1] = {"key": "value"}; arr[1]["key"];`,
			"value",
		},
		// Assigning different types to hashes
		{
			`fr hash = {}; hash["number"] = 42; hash["number"];`,
			42,
		},
		{
			`fr hash = {}; hash["string"] = "hello"; hash["string"];`,
			"hello",
		},
		{
			`fr hash = {}; hash["bool"] = cap; hash["bool"];`,
			false,
		},
		{
			`fr hash = {}; hash["array"] = [1, 2, 3]; count(hash["array"]);`,
			3,
		},
		{
			`fr hash = {}; hash["nested"] = {"inner": 1}; hash["nested"]["inner"];`,
			1,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testStringObject(t, evaluated, expected)
		case bool:
			testBooleanObject(t, evaluated, expected)
		}
	}
}

func TestIndexExpressionAssignmentInLoops(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		// Modifying array in for loop
		{
			`
			fr arr = [1, 2, 3, 4, 5];
			fr indices = [1, 3, 5];
			stalk (i in indices) {
				arr[i] = arr[i] * 10;
			}
			arr[1] + arr[3] + arr[5];
			`,
			90, // 10 + 30 + 50 = 90
		},
		// Building hash in loop
		{
			`
			fr hash = {};
			fr nums = [1, 2, 3];
			stalk (num in nums) {
				hash[num] = num * num;
			}
			hash[1] + hash[2] + hash[3];
			`,
			14, // 1 + 4 + 9
		},
		// Modifying nested structure in loop
		{
			`
			fr matrix = [[1, 2], [3, 4], [5, 6]];
			fr count = 1;
			stalk (row in matrix) {
				matrix[count][1] = matrix[count][1] + 10;
				count = count + 1;
			}
			matrix[1][1] + matrix[2][1] + matrix[3][1];
			`,
			39, // 11 + 13 + 15
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestCaughtIn4KLogging(t *testing.T) {
	tests := []struct {
		input        string
		expectedLogs []string
		description  string
	}{
		{
			input:        `caughtIn4K("Hello", "World");`,
			expectedLogs: []string{"Hello", "World"},
			description:  "Basic logging with multiple string arguments",
		},
		{
			input:        `caughtIn4K(42, noCap, "test");`,
			expectedLogs: []string{"42", "true", "test"},
			description:  "Logging with mixed types (integer, boolean, string)",
		},
		{
			input: `
				fr x = 10;
				fr y = "debug";
				caughtIn4K("Variable x:", x, "Variable y:", y);
			`,
			expectedLogs: []string{"Variable x:", "10", "Variable y:", "debug"},
			description:  "Logging with variables",
		},
		{
			input: `
				fr arr = [1, 2, 3];
				caughtIn4K("Array:", arr);
			`,
			expectedLogs: []string{"Array:", "[1, 2, 3]"},
			description:  "Logging with array",
		},
		{
			input: `
				fr logger = cook(msg, value) {
					caughtIn4K("Function log:", msg, value);
				};
				logger("test message", 123);
			`,
			expectedLogs: []string{"Function log:", "test message", "123"},
			description:  "Logging from within a function",
		},
		{
			input: `
				fr x = 1;
				caughtIn4K("First log:", x);
				x = x + 1;
				caughtIn4K("Second log:", x);
				fr nested = cook() {
					caughtIn4K("Nested function log:", x * 2);
				};
				nested();
			`,
			expectedLogs: []string{"First log:", "1", "Second log:", "2", "Nested function log:", "4"},
			description:  "Multiple logging calls including from nested function",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Create a new environment and evaluate the input to capture logs
			env := object.NewEnvironment()
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			Eval(program, env)

			if len(env.Logs) != len(tt.expectedLogs) {
				t.Errorf("Expected %d logs, got %d. Logs: %v",
					len(tt.expectedLogs), len(env.Logs), env.Logs)
				return
			}

			for i, expectedLog := range tt.expectedLogs {
				if env.Logs[i] != expectedLog {
					t.Errorf("Log %d: expected %q, got %q", i, expectedLog, env.Logs[i])
				}
			}
		})
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

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f",
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

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%q, want=%q",
			result.Value, expected)
		return false
	}
	return true
}
