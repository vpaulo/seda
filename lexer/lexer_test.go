package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `var x = 5
const name = "hello"
fn add(a, b) :: a + b end`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{VAR, "var"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{NUMBER, "5"},
		{CONST, "const"},
		{IDENT, "name"},
		{ASSIGN, "="},
		{STRING, "hello"},
		{FN, "fn"},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "a"},
		{COMMA, ","},
		{IDENT, "b"},
		{RPAREN, ")"},
		{DOUBLE_COLON, "::"},
		{IDENT, "a"},
		{PLUS, "+"},
		{IDENT, "b"},
		{END, "end"},
		{EOF, ""},
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

func TestKeywords(t *testing.T) {
	input := `var const fn struct type if else case for in check where end is self true false`

	expectedKeywords := []TokenType{
		VAR, CONST, FN, STRUCT, TYPE, IF, ELSE, CASE, FOR, IN, CHECK, WHERE, END, IS, SELF, TRUE, FALSE,
	}

	l := New(input)

	for i, expectedType := range expectedKeywords {
		tok := l.NextToken()
		if tok.Type != expectedType {
			t.Fatalf("keyword[%d] wrong. expected=%q, got=%q", i, expectedType, tok.Type)
		}
	}
}

func TestOperators(t *testing.T) {
	input := `= + - * / % == != < > <= >= :: => -> && || !`

	expectedOperators := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{ASSIGN, "="},
		{PLUS, "+"},
		{MINUS, "-"},
		{MULTIPLY, "*"},
		{DIVIDE, "/"},
		{MODULO, "%"},
		{EQ, "=="},
		{NOT_EQ, "!="},
		{LT, "<"},
		{GT, ">"},
		{LTE, "<="},
		{GTE, ">="},
		{DOUBLE_COLON, "::"},
		{ARROW, "=>"},
		{TYPE_ARROW, "->"},
		{AND, "&&"},
		{OR, "||"},
		{NOT, "!"},
	}

	l := New(input)

	for i, expected := range expectedOperators {
		tok := l.NextToken()
		if tok.Type != expected.expectedType {
			t.Fatalf("operator[%d] type wrong. expected=%q, got=%q",
				i, expected.expectedType, tok.Type)
		}
		if tok.Literal != expected.expectedLiteral {
			t.Fatalf("operator[%d] literal wrong. expected=%q, got=%q",
				i, expected.expectedLiteral, tok.Literal)
		}
	}
}

func TestDelimiters(t *testing.T) {
	input := `( ) [ ] { } , . :`

	expectedTokens := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LPAREN, "("},
		{RPAREN, ")"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{LBRACE, "{"},
		{RBRACE, "}"},
		{COMMA, ","},
		{DOT, "."},
		{COLON, ":"},
	}

	l := New(input)

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != expected.expectedType {
			t.Fatalf("delimiter[%d] type wrong. expected=%q, got=%q",
				i, expected.expectedType, tok.Type)
		}
		if tok.Literal != expected.expectedLiteral {
			t.Fatalf("delimiter[%d] literal wrong. expected=%q, got=%q",
				i, expected.expectedLiteral, tok.Literal)
		}
	}
}

func TestStringEscapeSequences(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello\nworld"`, "hello\nworld"},
		{`"tab\there"`, "tab\there"},
		{`"quote\"test"`, "quote\"test"},
		{`"backslash\\test"`, "backslash\\test"},
		{`"carriage\rreturn"`, "carriage\rreturn"},
		{`"null\0char"`, "null\x00char"},
		{`"mixed\n\t\"\\"`, "mixed\n\t\"\\"},
	}

	for i, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != STRING {
			t.Fatalf("test[%d] - token type wrong. expected=STRING, got=%q", i, tok.Type)
		}

		if tok.Literal != tt.expected {
			t.Errorf("test[%d] - escape sequence wrong. expected=%q, got=%q",
				i, tt.expected, tok.Literal)
		}
	}
}

func TestNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"42", "42"},
		{"0", "0"},
		{"123456789", "123456789"},
		{"999", "999"},
	}

	for i, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != NUMBER {
			t.Fatalf("test[%d] - token type wrong. expected=NUMBER, got=%q", i, tok.Type)
		}

		if tok.Literal != tt.expected {
			t.Errorf("test[%d] - number literal wrong. expected=%q, got=%q",
				i, tt.expected, tok.Literal)
		}
	}
}

func TestStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world"`, "hello world"},
		{`"test"`, "test"},
		{`""`, ""},
		{`"with spaces   "`, "with spaces   "},
		{`"123 numbers"`, "123 numbers"},
	}

	for i, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != STRING {
			t.Fatalf("test[%d] - token type wrong. expected=STRING, got=%q", i, tok.Type)
		}

		if tok.Literal != tt.expected {
			t.Errorf("test[%d] - string literal wrong. expected=%q, got=%q",
				i, tt.expected, tok.Literal)
		}
	}
}

func TestIdentifiers(t *testing.T) {
	tests := []string{
		"x", "y", "foo", "bar", "myVariable", "test123", "_underscore", "CamelCase",
	}

	for i, expected := range tests {
		l := New(expected)
		tok := l.NextToken()

		if tok.Type != IDENT {
			t.Fatalf("test[%d] - token type wrong. expected=IDENT, got=%q", i, tok.Type)
		}

		if tok.Literal != expected {
			t.Errorf("test[%d] - identifier wrong. expected=%q, got=%q",
				i, expected, tok.Literal)
		}
	}
}

func TestAllKeywords(t *testing.T) {
	input := `var const fn struct type if else case for in check where end is isA contains self true false return break module using as`

	expectedTokens := []TokenType{
		VAR, CONST, FN, STRUCT, TYPE, IF, ELSE, CASE, FOR, IN,
		CHECK, WHERE, END, IS, ISA, CONTAINS, SELF, TRUE, FALSE, RETURN, BREAK,
		MODULE, USING, AS,
	}

	l := New(input)

	for i, expectedType := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != expectedType {
			t.Fatalf("keyword[%d] wrong. expected=%q, got=%q (literal: %q)",
				i, expectedType, tok.Type, tok.Literal)
		}
	}
}

func TestComplexExpression(t *testing.T) {
	input := `var result = (a + b) * c - d / e`

	expectedTokens := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{VAR, "var"},
		{IDENT, "result"},
		{ASSIGN, "="},
		{LPAREN, "("},
		{IDENT, "a"},
		{PLUS, "+"},
		{IDENT, "b"},
		{RPAREN, ")"},
		{MULTIPLY, "*"},
		{IDENT, "c"},
		{MINUS, "-"},
		{IDENT, "d"},
		{DIVIDE, "/"},
		{IDENT, "e"},
		{EOF, ""},
	}

	l := New(input)

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != expected.expectedType {
			t.Fatalf("token[%d] type wrong. expected=%q, got=%q",
				i, expected.expectedType, tok.Type)
		}
		if tok.Literal != expected.expectedLiteral {
			t.Fatalf("token[%d] literal wrong. expected=%q, got=%q",
				i, expected.expectedLiteral, tok.Literal)
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	input := `[1, 2, 3, "hello", true]`

	expectedTokens := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LBRACKET, "["},
		{NUMBER, "1"},
		{COMMA, ","},
		{NUMBER, "2"},
		{COMMA, ","},
		{NUMBER, "3"},
		{COMMA, ","},
		{STRING, "hello"},
		{COMMA, ","},
		{TRUE, "true"},
		{RBRACKET, "]"},
		{EOF, ""},
	}

	l := New(input)

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != expected.expectedType {
			t.Fatalf("token[%d] type wrong. expected=%q, got=%q",
				i, expected.expectedType, tok.Type)
		}
		if tok.Literal != expected.expectedLiteral {
			t.Fatalf("token[%d] literal wrong. expected=%q, got=%q",
				i, expected.expectedLiteral, tok.Literal)
		}
	}
}

func TestMapLiteral(t *testing.T) {
	input := `{"name": "Alice", "age": 25}`

	expectedTokens := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LBRACE, "{"},
		{STRING, "name"},
		{COLON, ":"},
		{STRING, "Alice"},
		{COMMA, ","},
		{STRING, "age"},
		{COLON, ":"},
		{NUMBER, "25"},
		{RBRACE, "}"},
		{EOF, ""},
	}

	l := New(input)

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != expected.expectedType {
			t.Fatalf("token[%d] type wrong. expected=%q, got=%q",
				i, expected.expectedType, tok.Type)
		}
		if tok.Literal != expected.expectedLiteral {
			t.Fatalf("token[%d] literal wrong. expected=%q, got=%q",
				i, expected.expectedLiteral, tok.Literal)
		}
	}
}

func TestDotAccess(t *testing.T) {
	input := `obj.property.method()`

	expectedTokens := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{IDENT, "obj"},
		{DOT, "."},
		{IDENT, "property"},
		{DOT, "."},
		{IDENT, "method"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{EOF, ""},
	}

	l := New(input)

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != expected.expectedType {
			t.Fatalf("token[%d] type wrong. expected=%q, got=%q",
				i, expected.expectedType, tok.Type)
		}
		if tok.Literal != expected.expectedLiteral {
			t.Fatalf("token[%d] literal wrong. expected=%q, got=%q",
				i, expected.expectedLiteral, tok.Literal)
		}
	}
}

func TestWhitespaceHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected []TokenType
	}{
		{"   x   ", []TokenType{IDENT, EOF}},
		{"\tx\t", []TokenType{IDENT, EOF}},
		{"\n\nx\n\n", []TokenType{IDENT, EOF}},
		{"  \t\n  x  \t\n  ", []TokenType{IDENT, EOF}},
	}

	for i, tt := range tests {
		l := New(tt.input)
		for j, expectedType := range tt.expected {
			tok := l.NextToken()
			if tok.Type != expectedType {
				t.Fatalf("test[%d] token[%d] type wrong. expected=%q, got=%q",
					i, j, expectedType, tok.Type)
			}
		}
	}
}

func TestLineTracking(t *testing.T) {
	input := `var x = 5; const y = 10`

	expectedTokens := []struct {
		expectedType TokenType
		expectedLine int
	}{
		{VAR, 1},
		{IDENT, 1},
		{ASSIGN, 1},
		{NUMBER, 1},
		{SEMICOLON, 1},
		{CONST, 1},
		{IDENT, 1},
		{ASSIGN, 1},
		{NUMBER, 1},
	}

	l := New(input)

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != expected.expectedType {
			t.Fatalf("token[%d] type wrong. expected=%q, got=%q",
				i, expected.expectedType, tok.Type)
		}
		if tok.Line != expected.expectedLine {
			t.Errorf("token[%d] line wrong. expected=%d, got=%d",
				i, expected.expectedLine, tok.Line)
		}
	}
}

func TestMultilineString(t *testing.T) {
	input := `"line1
line2
line3"`

	l := New(input)
	tok := l.NextToken()

	if tok.Type != STRING {
		t.Fatalf("token type wrong. expected=STRING, got=%q", tok.Type)
	}

	expected := "line1\nline2\nline3"
	if tok.Literal != expected {
		t.Errorf("multiline string wrong. expected=%q, got=%q", expected, tok.Literal)
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{"empty input", "", []TokenType{EOF}},
		{"only whitespace", "   \t\n   ", []TokenType{EOF}},
		{"single char", "x", []TokenType{IDENT, EOF}},
		{"single number", "5", []TokenType{NUMBER, EOF}},
		{"empty string", `""`, []TokenType{STRING, EOF}},
		{"empty array", "[]", []TokenType{LBRACKET, RBRACKET, EOF}},
		{"empty map", "{}", []TokenType{LBRACE, RBRACE, EOF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			for i, expectedType := range tt.expected {
				tok := l.NextToken()
				if tok.Type != expectedType {
					t.Fatalf("token[%d] type wrong. expected=%q, got=%q",
						i, expectedType, tok.Type)
				}
			}
		})
	}
}

