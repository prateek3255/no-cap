package parser

import (
	"fmt"
	"nocap/ast"
	"nocap/lexer"
	"strconv"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"fr x = 5;", "x", 5},
		{"fr y = noCap;", "y", true},
		{"fr foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"yeet 5;", 5},
		{"yeet noCap;", true},
		{"yeet foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "yeet" {
			t.Fatalf("returnStmt.TokenLiteral not 'yeet', got %q",
				returnStmt.TokenLiteral())
		}
		if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestFunctionStatement(t *testing.T) {
	input := `cook sup(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.statements does not contain %d statments. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "sup" {
		t.Fatalf("stmt.Name.Value not '%s'. got=%s", "sup", stmt.Name.Value)
	}

	if len(stmt.Parameters) != 2 {
		t.Fatalf("stmt.Parameters does not contain 2 parameters. got=%d\n",
			len(stmt.Parameters))
	}

	if stmt.Parameters[0].Value != "x" {
		t.Fatalf("stmt.Parameters[0].Value not '%s'. got=%s", "x", stmt.Parameters[0].Value)
	}

	if stmt.Parameters[1].Value != "y" {
		t.Fatalf("stmt.Parameters[1].Value not '%s'. got=%s", "y", stmt.Parameters[1].Value)
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("stmt.Body.Statements does not contain 1 statements. got=%d\n",
			len(stmt.Body.Statements))
	}

	bodyStmt, ok := stmt.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Body.Statements[0] is not ast.ExpressionStatement. got=%T",
			stmt.Body.Statements[0])
	}

	if !testInfixExpression(t, bodyStmt.Expression, "x", "+", "y") {
		return
	}
}

func TestForStatement(t *testing.T) {
	input := `stalk (x in y) { 
		x;
	}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.statements does not contain %d statments. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForStatment. got=%T",
			program.Statements[0])
	}

	if stmt.Key.Value != "x" {
		t.Fatalf("stmt.Key.Value not '%s'. got=%s", "x", stmt.Key.Value)
	}

	if !testIdentifier(t, stmt.Items, "y") {
		return
	}

	if len(stmt.Body.Statements) != 1 {
		t.Errorf("body is not 1 statements. got=%d\n",
			len(stmt.Body.Statements))
	}

	body, ok := stmt.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			stmt.Body.Statements[0])
	}

	if !testIdentifier(t, body.Expression, "x") {
		return
	}
}

func TestWhileStatement(t *testing.T) {
	input := `onRepeat (x < y) { 
		x;
	}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.statements does not contain %d statments. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatment. got=%T",
			program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", "<", "y") {
		return
	}

	if len(stmt.Body.Statements) != 1 {
		t.Errorf("body is not 1 statements. got=%d\n",
			len(stmt.Body.Statements))
	}

	body, ok := stmt.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			stmt.Body.Statements[0])
	}

	if !testIdentifier(t, body.Expression, "x") {
		return
	}
}

func TestContinueStatement(t *testing.T) {
	input := "pass;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ContinueStatement. got=%T", program.Statements[0])
	}
	if stmt.TokenLiteral() != "pass" {
		t.Errorf("stmt.TokenLiteral not 'pass'. got=%q", stmt.TokenLiteral())
	}
}

func TestBreakStatement(t *testing.T) {
	input := "bounce;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("stmt not *ast.BreakStatement. got=%T", program.Statements[0])
	}
	if stmt.TokenLiteral() != "bounce" {
		t.Errorf("stmt.TokenLiteral not 'bounce'. got=%q", stmt.TokenLiteral())
	}
}

func TestAssignmentStatement(t *testing.T) {
	input := "count = count + 1;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}
	if stmt.Name.Value != "count" {
		t.Errorf("stmt.Name.Value not 'count'. got=%s", stmt.Name.Value)
	}
	if stmt.TokenLiteral() != "=" {
		t.Errorf("stmt.TokenLiteral not '='. got=%q", stmt.TokenLiteral())
	}
	if !testInfixExpression(t, stmt.Value, "count", "+", 1) {
		return
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func TestFloatLiteralExpression(t *testing.T) {
	input := "3.14;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("exp not *ast.FloatLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 3.14 {
		t.Errorf("literal.Value not %f. got=%f", 3.14, literal.Value)
	}
	if literal.TokenLiteral() != "3.14" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "3.14",
			literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"nah 5;", "nah", 5},
		{"-15;", "-", 15},
		{"-3.14;", "-", 3.14},
		{"nah foobar;", "nah", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"nah noCap;", "nah", true},
		{"nah cap;", "nah", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 % 5;", 5, "%", 5},
		{"5.0 > 5;", 5.0, ">", 5},
		{"5 is 5;", 5, "is", 5},
		{"5 aint 5;", 5, "aint", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar is barfoo;", "foobar", "is", "barfoo"},
		{"foobar aint barfoo;", "foobar", "aint", "barfoo"},
		{"noCap is noCap", true, "is", true},
		{"noCap aint cap", true, "aint", false},
		{"cap is cap", false, "is", false},
		{"noCap and cap", true, "and", false},
		{"cap or noCap", false, "or", true},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"nah -a",
			"(nah (-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 is 3 < 4",
			"((5 > 4) is (3 < 4))",
		},
		{
			"5 < 4 aint 3 > 4",
			"((5 < 4) aint (3 > 4))",
		},
		{
			"3 + 4 * 5 is 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) is ((3 * 1) + (4 * 5)))",
		},
		{
			"noCap",
			"noCap",
		},
		{
			"cap",
			"cap",
		},
		{
			"3 > 5 is cap",
			"((3 > 5) is cap)",
		},
		{
			"3 < 5 is noCap",
			"((3 < 5) is noCap)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"-(4.15 + 5.15)",
			"(-(4.15 + 5.15))",
		},
		{
			"nah (noCap is noCap)",
			"(nah (noCap is noCap))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
		{
			"noCap and cap or noCap",
			"((noCap and cap) or noCap)",
		},
		{
			"cap or noCap and cap",
			"(cap or (noCap and cap))",
		},
		{
			"5 > 3 and 2 < 4",
			"((5 > 3) and (2 < 4))",
		},
		{
			"5 < 3 or 2 > 1 and 3 is 3",
			"((5 < 3) or ((2 > 1) and (3 is 3)))",
		},
		{
			"nah x and y or z",
			"(((nah x) and y) or z)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"noCap;", true},
		{"cap;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

func TestNullExpression(t *testing.T) {
	input := "ghosted;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	null, ok := stmt.Expression.(*ast.Null)
	if !ok {
		t.Fatalf("exp not *ast.Null. got=%T", stmt.Expression)
	}
	if null.TokenLiteral() != "ghosted" {
		t.Errorf("null.TokenLiteral not 'ghosted'. got=%q", null.TokenLiteral())
	}
}

func TestIfExpression(t *testing.T) {
	input := `vibe (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `vibe (x < y) { x } nvm { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `cook(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "cook() {};", expectedParams: []string{}},
		{input: "cook(x) {};", expectedParams: []string{"x"}},
		{input: "cook(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestParsingEmptyArrayLiterals(t *testing.T) {
	input := "[]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 0 {
		t.Errorf("len(array.Elements) not 0. got=%d", len(array.Elements))
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsBooleanKeys(t *testing.T) {
	input := `{noCap: 1, cap: 2}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"noCap": 1,
		"cap":   2,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		boolean, ok := key.(*ast.Boolean)
		if !ok {
			t.Errorf("key is not ast.BooleanLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[boolean.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		integer, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key is not ast.IntegerLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[integer.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(value)
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "fr" {
		t.Errorf("s.TokenLiteral not 'fr'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case float64:
		return testFloatLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}

	return true
}

func testFloatLiteral(t *testing.T, fl ast.Expression, value float64) bool {
	fmt.Printf("testFloatLiteral: %v\n", value)
	floatLit, ok := fl.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("fl not *ast.FloatLiteral. got=%T", fl)
		return false
	}

	if floatLit.Value != value {
		t.Errorf("floatLit.Value not %f. got=%f", value, floatLit.Value)
		return false
	}

	// Parse the expected value from the token literal instead of formatting the float
	expectedTokenLiteral := floatLit.TokenLiteral()
	parsedValue, err := strconv.ParseFloat(expectedTokenLiteral, 64)
	if err != nil {
		t.Errorf("failed to parse token literal as float: %s", expectedTokenLiteral)
		return false
	}

	if parsedValue != value {
		t.Errorf("floatLit.TokenLiteral parsed value not %f. got=%f from %s",
			value, parsedValue, expectedTokenLiteral)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	expectedLiteral := "cap"
	if value {
		expectedLiteral = "noCap"
	}
	if bo.TokenLiteral() != expectedLiteral {
		t.Errorf("bo.TokenLiteral not %s. got=%s",
			expectedLiteral, bo.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	str, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("exp not *ast.StringLiteral. got=%T", exp)
		return false
	}

	if str.Value != value {
		t.Errorf("str.Value not %s. got=%s", value, str.Value)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestIndexExpressionAssignmentStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedIndex      interface{}
		expectedValue      interface{}
	}{
		{"x[5] = 10;", "x", 5, 10},
		{"arr[0] = noCap;", "arr", 0, true},
		{"myArray[key] = value;", "myArray", "key", "value"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.IndexExpressionAssignmentStatement)
		if !ok {
			t.Fatalf("stmt not *ast.IndexExpressionAssignmentStatement. got=%T", program.Statements[0])
		}

		if !testIdentifier(t, stmt.Left.Left, tt.expectedIdentifier) {
			return
		}

		if !testLiteralExpression(t, stmt.Left.Index, tt.expectedIndex) {
			return
		}

		if !testLiteralExpression(t, stmt.Value, tt.expectedValue) {
			return
		}
	}
}

func TestComplexIndexExpressionAssignmentStatement(t *testing.T) {
	tests := []struct {
		input       string
		description string
	}{
		{"x[5 + 6] = 5;", "arithmetic expression in index"},
		{"x[y[8] - z[9]] = 5;", "nested index expressions with arithmetic"},
		{"x[5][\"h\"] = 5;", "chained index expressions"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.IndexExpressionAssignmentStatement)
			if !ok {
				t.Fatalf("stmt not *ast.IndexExpressionAssignmentStatement. got=%T", program.Statements[0])
			}

			// Verify that we have a valid index expression on the left
			if stmt.Left == nil {
				t.Fatalf("stmt.Left is nil")
			}

			// Verify that we have a value on the right
			if stmt.Value == nil {
				t.Fatalf("stmt.Value is nil")
			}

			// For the specific test cases, verify the structure
			switch tt.input {
			case "x[5 + 6] = 5;":
				// Left should be x
				if !testIdentifier(t, stmt.Left.Left, "x") {
					return
				}
				// Index should be 5 + 6
				if !testInfixExpression(t, stmt.Left.Index, 5, "+", 6) {
					return
				}
				// Value should be 5
				if !testLiteralExpression(t, stmt.Value, 5) {
					return
				}

			case "x[y[8] - z[9]] = 5;":
				// Left should be x
				if !testIdentifier(t, stmt.Left.Left, "x") {
					return
				}
				// Index should be y[8] - z[9]
				indexInfix, ok := stmt.Left.Index.(*ast.InfixExpression)
				if !ok {
					t.Fatalf("stmt.Left.Index not *ast.InfixExpression. got=%T", stmt.Left.Index)
				}
				if indexInfix.Operator != "-" {
					t.Errorf("indexInfix.Operator not '-'. got=%s", indexInfix.Operator)
				}

				// Left side should be y[8]
				leftIndex, ok := indexInfix.Left.(*ast.IndexExpression)
				if !ok {
					t.Fatalf("indexInfix.Left not *ast.IndexExpression. got=%T", indexInfix.Left)
				}
				if !testIdentifier(t, leftIndex.Left, "y") {
					return
				}
				if !testLiteralExpression(t, leftIndex.Index, 8) {
					return
				}

				// Right side should be z[9]
				rightIndex, ok := indexInfix.Right.(*ast.IndexExpression)
				if !ok {
					t.Fatalf("indexInfix.Right not *ast.IndexExpression. got=%T", indexInfix.Right)
				}
				if !testIdentifier(t, rightIndex.Left, "z") {
					return
				}
				if !testLiteralExpression(t, rightIndex.Index, 9) {
					return
				}

				// Value should be 5
				if !testLiteralExpression(t, stmt.Value, 5) {
					return
				}

			case "x[5][\"h\"] = 5;":
				// This should parse as (x[5])["h"] = 5
				// So stmt.Left should be a chained index expression
				innerIndex, ok := stmt.Left.Left.(*ast.IndexExpression)
				if !ok {
					t.Fatalf("stmt.Left.Left not *ast.IndexExpression. got=%T", stmt.Left.Left)
				}
				if !testIdentifier(t, innerIndex.Left, "x") {
					return
				}
				if !testLiteralExpression(t, innerIndex.Index, 5) {
					return
				}
				if !testStringLiteral(t, stmt.Left.Index, "h") {
					return
				}
				if !testLiteralExpression(t, stmt.Value, 5) {
					return
				}
			}
		})
	}
}

func TestIndexExpressionStillWorksAsExpression(t *testing.T) {
	tests := []struct {
		input       string
		description string
	}{
		{"x[4] + y[9];", "index expressions in arithmetic"},
		{"x[4];", "standalone index expression"},
		{"fn()[5];", "index on function call result"},
		{"arr[i][j];", "chained index expressions"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("stmt not *ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			// Verify that the expression is parsed correctly based on the test case
			switch tt.input {
			case "x[4] + y[9];":
				// Should be an infix expression
				infixExp, ok := stmt.Expression.(*ast.InfixExpression)
				if !ok {
					t.Fatalf("stmt.Expression not *ast.InfixExpression. got=%T", stmt.Expression)
				}
				if infixExp.Operator != "+" {
					t.Errorf("infixExp.Operator not '+'. got=%s", infixExp.Operator)
				}

				// Left should be x[4]
				leftIndex, ok := infixExp.Left.(*ast.IndexExpression)
				if !ok {
					t.Fatalf("infixExp.Left not *ast.IndexExpression. got=%T", infixExp.Left)
				}
				if !testIdentifier(t, leftIndex.Left, "x") {
					return
				}
				if !testLiteralExpression(t, leftIndex.Index, 4) {
					return
				}

				// Right should be y[9]
				rightIndex, ok := infixExp.Right.(*ast.IndexExpression)
				if !ok {
					t.Fatalf("infixExp.Right not *ast.IndexExpression. got=%T", infixExp.Right)
				}
				if !testIdentifier(t, rightIndex.Left, "y") {
					return
				}
				if !testLiteralExpression(t, rightIndex.Index, 9) {
					return
				}

			case "x[4];":
				// Should be an index expression
				indexExp, ok := stmt.Expression.(*ast.IndexExpression)
				if !ok {
					t.Fatalf("stmt.Expression not *ast.IndexExpression. got=%T", stmt.Expression)
				}
				if !testIdentifier(t, indexExp.Left, "x") {
					return
				}
				if !testLiteralExpression(t, indexExp.Index, 4) {
					return
				}

			case "fn()[5];":
				// Should be an index expression where left is a call expression
				indexExp, ok := stmt.Expression.(*ast.IndexExpression)
				if !ok {
					t.Fatalf("stmt.Expression not *ast.IndexExpression. got=%T", stmt.Expression)
				}

				callExp, ok := indexExp.Left.(*ast.CallExpression)
				if !ok {
					t.Fatalf("indexExp.Left not *ast.CallExpression. got=%T", indexExp.Left)
				}
				if !testIdentifier(t, callExp.Function, "fn") {
					return
				}
				if !testLiteralExpression(t, indexExp.Index, 5) {
					return
				}

			case "arr[i][j];":
				// Should be a chained index expression
				indexExp, ok := stmt.Expression.(*ast.IndexExpression)
				if !ok {
					t.Fatalf("stmt.Expression not *ast.IndexExpression. got=%T", stmt.Expression)
				}

				// Left should be arr[i]
				leftIndex, ok := indexExp.Left.(*ast.IndexExpression)
				if !ok {
					t.Fatalf("indexExp.Left not *ast.IndexExpression. got=%T", indexExp.Left)
				}
				if !testIdentifier(t, leftIndex.Left, "arr") {
					return
				}
				if !testIdentifier(t, leftIndex.Index, "i") {
					return
				}

				// Index should be j
				if !testIdentifier(t, indexExp.Index, "j") {
					return
				}
			}
		})
	}
}

func TestIfElseIfExpression(t *testing.T) {
	input := `vibe (x < y) { x } unless (x > z) { z } nvm { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	// Test main condition
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	// Test consequence
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// Test else if
	if len(exp.ElseIfs) != 1 {
		t.Fatalf("exp.ElseIfs does not contain 1 else if. got=%d", len(exp.ElseIfs))
	}

	elseIf := exp.ElseIfs[0]
	if !testInfixExpression(t, elseIf.Condition, "x", ">", "z") {
		return
	}

	if len(elseIf.Consequence.Statements) != 1 {
		t.Errorf("elseIf consequence is not 1 statements. got=%d\n",
			len(elseIf.Consequence.Statements))
	}

	elseIfConsequence, ok := elseIf.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("elseIf Statements[0] is not ast.ExpressionStatement. got=%T",
			elseIf.Consequence.Statements[0])
	}

	if !testIdentifier(t, elseIfConsequence.Expression, "z") {
		return
	}

	// Test alternative (else)
	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestMultipleElseIfExpressions(t *testing.T) {
	input := `vibe (x < 1) { "first" } unless (x < 2) { "second" } unless (x < 3) { "third" } nvm { "default" }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	// Test main condition
	if !testInfixExpression(t, exp.Condition, "x", "<", 1) {
		return
	}

	// Test main consequence
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testStringLiteral(t, consequence.Expression, "first") {
		return
	}

	// Test multiple else ifs
	if len(exp.ElseIfs) != 2 {
		t.Fatalf("exp.ElseIfs does not contain 2 else ifs. got=%d", len(exp.ElseIfs))
	}

	// Test first else if
	firstElseIf := exp.ElseIfs[0]
	if !testInfixExpression(t, firstElseIf.Condition, "x", "<", 2) {
		return
	}

	firstElseIfConsequence, ok := firstElseIf.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("First elseIf Statements[0] is not ast.ExpressionStatement. got=%T",
			firstElseIf.Consequence.Statements[0])
	}

	if !testStringLiteral(t, firstElseIfConsequence.Expression, "second") {
		return
	}

	// Test second else if
	secondElseIf := exp.ElseIfs[1]
	if !testInfixExpression(t, secondElseIf.Condition, "x", "<", 3) {
		return
	}

	secondElseIfConsequence, ok := secondElseIf.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Second elseIf Statements[0] is not ast.ExpressionStatement. got=%T",
			secondElseIf.Consequence.Statements[0])
	}

	if !testStringLiteral(t, secondElseIfConsequence.Expression, "third") {
		return
	}

	// Test alternative (else)
	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testStringLiteral(t, alternative.Expression, "default") {
		return
	}
}

func TestIfElseIfWithoutElse(t *testing.T) {
	input := `vibe (x < y) { x } unless (x > z) { z }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	// Test main condition
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	// Test consequence
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// Test else if
	if len(exp.ElseIfs) != 1 {
		t.Fatalf("exp.ElseIfs does not contain 1 else if. got=%d", len(exp.ElseIfs))
	}

	elseIf := exp.ElseIfs[0]
	if !testInfixExpression(t, elseIf.Condition, "x", ">", "z") {
		return
	}

	elseIfConsequence, ok := elseIf.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("elseIf Statements[0] is not ast.ExpressionStatement. got=%T",
			elseIf.Consequence.Statements[0])
	}

	if !testIdentifier(t, elseIfConsequence.Expression, "z") {
		return
	}

	// Test that there's no alternative (else)
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative should be nil, got=%+v", exp.Alternative)
	}
}

func TestComplexElseIfConditions(t *testing.T) {
	input := `vibe (x + y > 10) { "big" } unless (x * y is 0) { "zero product" } unless (x is y) { "equal" } nvm { "other" }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	// Test main condition (x + y > 10)
	mainCondition, ok := exp.Condition.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp.Condition is not ast.InfixExpression. got=%T", exp.Condition)
	}
	if mainCondition.Operator != ">" {
		t.Errorf("mainCondition.Operator not '>'. got=%s", mainCondition.Operator)
		return
	}
	if !testInfixExpression(t, mainCondition.Left, "x", "+", "y") {
		return
	}
	if !testLiteralExpression(t, mainCondition.Right, 10) {
		return
	}

	// Test main consequence
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testStringLiteral(t, consequence.Expression, "big") {
		return
	}

	// Test else ifs
	if len(exp.ElseIfs) != 2 {
		t.Fatalf("exp.ElseIfs does not contain 2 else ifs. got=%d", len(exp.ElseIfs))
	}

	// Test first else if (x * y is 0)
	firstElseIf := exp.ElseIfs[0]
	firstElseIfCondition, ok := firstElseIf.Condition.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("firstElseIf.Condition is not ast.InfixExpression. got=%T", firstElseIf.Condition)
	}
	if firstElseIfCondition.Operator != "is" {
		t.Errorf("firstElseIfCondition.Operator not 'is'. got=%s", firstElseIfCondition.Operator)
		return
	}
	if !testInfixExpression(t, firstElseIfCondition.Left, "x", "*", "y") {
		return
	}
	if !testLiteralExpression(t, firstElseIfCondition.Right, 0) {
		return
	}

	firstElseIfConsequence, ok := firstElseIf.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("First elseIf Statements[0] is not ast.ExpressionStatement. got=%T",
			firstElseIf.Consequence.Statements[0])
	}

	if !testStringLiteral(t, firstElseIfConsequence.Expression, "zero product") {
		return
	}

	// Test second else if (x is y)
	secondElseIf := exp.ElseIfs[1]
	if !testInfixExpression(t, secondElseIf.Condition, "x", "is", "y") {
		return
	}

	secondElseIfConsequence, ok := secondElseIf.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Second elseIf Statements[0] is not ast.ExpressionStatement. got=%T",
			secondElseIf.Consequence.Statements[0])
	}

	if !testStringLiteral(t, secondElseIfConsequence.Expression, "equal") {
		return
	}

	// Test alternative (else)
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testStringLiteral(t, alternative.Expression, "other") {
		return
	}
}
