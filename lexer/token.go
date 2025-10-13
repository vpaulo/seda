package lexer

import "fmt"

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Literals
	IDENT   // variable names, function names
	NUMBER  // 123, 45.67
	STRING  // "hello world"
	BOOLEAN // true, false

	// Operators
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /
	MODULO   // %
	POWER    // ^

	// Comparison
	EQ     // ==
	NOT_EQ // !=
	LT     // <
	GT     // >
	LTE    // <=
	GTE    // >=

	// Logical
	AND // and, &&
	OR  // or, ||
	NOT // not, !

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	DOT       // .

	// Brackets
	LPAREN   // (
	RPAREN   // )
	LBRACKET // [
	RBRACKET // ]
	LBRACE   // {
	RBRACE   // }

	// Block delimiters
	DOUBLE_COLON // ::

	// Keywords
	VAR      // var
	CONST    // const
	FN       // fn
	STRUCT   // struct
	TYPE     // type
	MODULE   // module
	USING    // using
	AS       // as
	IF       // if
	ELSE     // else
	CASE     // case
	FOR      // for
	IN       // in
	CHECK    // check
	WHERE    // where
	END      // end
	IS       // is
	ISA      // isA
	ISNOT    // isNot
	CONTAINS // contains
	ISGREATER // isGreater
	ISLESS    // isLess
	ISTRUE    // isTrue
	ISFALSE   // isFalse
	ISEMPTY   // isEmpty
	STARTSWITH // startsWith
	ENDSWITH   // endsWith
	RAISES     // raises
	SELF     // self
	RETURN   // return
	BREAK    // break
	TRUE     // true
	FALSE    // false
	NIL      // nil

	// Type keywords
	NUMBER_TYPE  // number
	STRING_TYPE  // string
	BOOLEAN_TYPE // boolean

	// Special operators
	ARROW           // =>
	TYPE_ARROW      // ->
	RANGE           // ..
	RANGE_INCLUSIVE // ...

	// Comments
	COMMENT // # comment
)

// Token represents a single token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// String returns a string representation of the token type
func (t TokenType) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case BOOLEAN:
		return "BOOLEAN"
	case ASSIGN:
		return "="
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case MULTIPLY:
		return "*"
	case DIVIDE:
		return "/"
	case MODULO:
		return "%"
	case POWER:
		return "^"
	case EQ:
		return "=="
	case NOT_EQ:
		return "!="
	case LT:
		return "<"
	case GT:
		return ">"
	case LTE:
		return "<="
	case GTE:
		return ">="
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case COMMA:
		return ","
	case SEMICOLON:
		return ";"
	case COLON:
		return ":"
	case DOT:
		return "."
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case LBRACKET:
		return "["
	case RBRACKET:
		return "]"
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case DOUBLE_COLON:
		return "::"
	case VAR:
		return "var"
	case CONST:
		return "const"
	case FN:
		return "fn"
	case STRUCT:
		return "struct"
	case TYPE:
		return "type"
	case MODULE:
		return "module"
	case USING:
		return "using"
	case AS:
		return "as"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case CASE:
		return "case"
	case FOR:
		return "for"
	case IN:
		return "in"
	case CHECK:
		return "check"
	case WHERE:
		return "where"
	case END:
		return "end"
	case IS:
		return "is"
	case ISA:
		return "isA"
	case ISNOT:
		return "isNot"
	case CONTAINS:
		return "contains"
	case ISGREATER:
		return "isGreater"
	case ISLESS:
		return "isLess"
	case ISTRUE:
		return "isTrue"
	case ISFALSE:
		return "isFalse"
	case ISEMPTY:
		return "isEmpty"
	case STARTSWITH:
		return "startsWith"
	case ENDSWITH:
		return "endsWith"
	case RAISES:
		return "raises"
	case SELF:
		return "self"
	case RETURN:
		return "return"
	case BREAK:
		return "break"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case NIL:
		return "nil"
	case NUMBER_TYPE:
		return "number"
	case STRING_TYPE:
		return "string"
	case BOOLEAN_TYPE:
		return "boolean"
	case ARROW:
		return "=>"
	case TYPE_ARROW:
		return "->"
	case RANGE:
		return ".."
	case RANGE_INCLUSIVE:
		return "..."
	case COMMENT:
		return "COMMENT"
	default:
		return "UNKNOWN"
	}
}

// String returns a string representation of the token
func (t Token) String() string {
	return fmt.Sprintf("{Type: %s, Literal: %q, Line: %d, Column: %d}",
		t.Type, t.Literal, t.Line, t.Column)
}

// keywords maps string literals to their token types
var keywords = map[string]TokenType{
	"var":      VAR,
	"const":    CONST,
	"fn":       FN,
	"struct":   STRUCT, // TODO: don't think i need this keyword
	"type":     TYPE,
	"module":   MODULE,
	"using":    USING,
	"as":       AS,
	"if":       IF,
	"else":     ELSE,
	"case":     CASE,
	"for":      FOR,
	"in":       IN,
	"check":    CHECK,
	"where":    WHERE,
	"end":        END,
	"is":         IS,
	"isA":        ISA,
	"isNot":      ISNOT,
	"contains":   CONTAINS,
	"isGreater":  ISGREATER,
	"isLess":     ISLESS,
	"isTrue":     ISTRUE,
	"isFalse":    ISFALSE,
	"isEmpty":    ISEMPTY,
	"startsWith": STARTSWITH,
	"endsWith":   ENDSWITH,
	"raises":     RAISES,
	"self":       SELF,
	"return":   RETURN,
	"break":    BREAK,
	"true":     TRUE,
	"false":    FALSE,
	"nil":      NIL,
	"number":   NUMBER_TYPE,
	"string":   STRING_TYPE,
	"boolean":  BOOLEAN_TYPE,
	"and":      AND,
	"or":       OR,
	"not":      NOT,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
