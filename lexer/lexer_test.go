package lexer

import (
	"testing"

	"nocap/token"
)

func TestNextToken(t *testing.T) {
	input := `fr five = 5;
fr ten = 10;

// Single line comment before function
fr add = cook(x, y) {
  x + y;
};

fr result = add(five, ten); // Single line comment at end of line
nah-*/5;
5 < 10 > 5;
5 <= 10 >= 5;
10 % 5;

/* Multiline comment
That spans across multiple
lines */
vibe (5 < 10) {
	yeet noCap;
} nvm {
	yeet cap;
}

10 is 10;
10 aint 9;
noCap and cap;
cap or noCap;
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
fr x =.04
x * 5.05
// Comment at the end of file
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "fr"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "fr"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "fr"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "cook"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "fr"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "nah"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LTE, "<="},
		{token.INT, "10"},
		{token.GTE, ">="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.MODULO, "%"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "vibe"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "yeet"},
		{token.TRUE, "noCap"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "nvm"},
		{token.LBRACE, "{"},
		{token.RETURN, "yeet"},
		{token.FALSE, "cap"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "is"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "aint"},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.TRUE, "noCap"},
		{token.AND, "and"},
		{token.FALSE, "cap"},
		{token.SEMICOLON, ";"},
		{token.FALSE, "cap"},
		{token.OR, "or"},
		{token.TRUE, "noCap"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.LET, "fr"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.FLOAT, ".04"},
		{token.IDENT, "x"},
		{token.ASTERISK, "*"},
		{token.FLOAT, "5.05"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIdentifiersWithDigits(t *testing.T) {
	input := `caughtIn4K var2 test123 _test4 test_5_more a1b2c3 count`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "caughtIn4K"},
		{token.IDENT, "var2"},
		{token.IDENT, "test123"},
		{token.IDENT, "_test4"},
		{token.IDENT, "test_5_more"},
		{token.IDENT, "a1b2c3"},
		{token.IDENT, "count"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNumbersStillWorkCorrectly(t *testing.T) {
	input := `123 456.789 .123 0.5 999`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "123"},
		{token.FLOAT, "456.789"},
		{token.FLOAT, ".123"},
		{token.FLOAT, "0.5"},
		{token.INT, "999"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIdentifierNumberBoundary(t *testing.T) {
	// Test that identifiers and numbers are properly separated
	input := `var1 123 test2 456.789 a3b 0.5`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "var1"},
		{token.INT, "123"},
		{token.IDENT, "test2"},
		{token.FLOAT, "456.789"},
		{token.IDENT, "a3b"},
		{token.FLOAT, "0.5"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestEdgeCasesForIdentifiersWithDigits(t *testing.T) {
	// Test edge cases to ensure we don't break existing functionality
	input := `_1 a_ _a 1a 1.5a`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "_1"},
		{token.IDENT, "a_"},
		{token.IDENT, "_a"},
		{token.INT, "1"},
		{token.IDENT, "a"},
		{token.FLOAT, "1.5"},
		{token.IDENT, "a"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
