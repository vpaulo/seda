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
	input := `var const fn type module using as if else case for in check where end is isA contains self return break true false and or not`

	expectedKeywords := []TokenType{
		VAR, CONST, FN, TYPE, MODULE, USING, AS, IF, ELSE, CASE, FOR, IN, CHECK, WHERE, END, IS, ISA, CONTAINS, SELF, RETURN, BREAK, TRUE, FALSE, AND, OR, NOT,
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

func TestMultilineComments(t *testing.T) {
	input := `#| This is a multiline comment
that spans multiple lines
and should be ignored |#
var x = 5`

	l := New(input)

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{COMMENT, "#| This is a multiline comment\nthat spans multiple lines\nand should be ignored |#"},
		{VAR, "var"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{NUMBER, "5"},
		{EOF, ""},
	}

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

func TestMultilineCommentInline(t *testing.T) {
	input := `var x = 5 #| inline comment |# var y = 10`

	l := New(input)

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{VAR, "var"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{NUMBER, "5"},
		{COMMENT, "#| inline comment |#"},
		{VAR, "var"},
		{IDENT, "y"},
		{ASSIGN, "="},
		{NUMBER, "10"},
		{EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
	}
}

func TestMultilineCommentUnclosed(t *testing.T) {
	input := `var x = 5
#| This comment is not closed
var y = 10`

	l := New(input)

	// Should still tokenize without crashing
	for {
		tok := l.NextToken()
		if tok.Type == EOF {
			break
		}
	}
}

func TestMultipleMultilineComments(t *testing.T) {
	input := `#| First comment |#
var x = 5
#| Second comment |#
var y = 10`

	l := New(input)

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{COMMENT, "#| First comment |#"},
		{VAR, "var"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{NUMBER, "5"},
		{COMMENT, "#| Second comment |#"},
		{VAR, "var"},
		{IDENT, "y"},
		{ASSIGN, "="},
		{NUMBER, "10"},
		{EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
	}
}

func TestMixedComments(t *testing.T) {
	input := `# Single line comment
var x = 5
#| Multiline
   comment |#
var y = 10 # Another single line`

	l := New(input)

	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}

	// Should have both types of comments
	hasComment := false
	for _, tok := range tokens {
		if tok.Type == COMMENT {
			hasComment = true
		}
	}

	if !hasComment {
		t.Fatalf("Expected to find COMMENT tokens")
	}
}
