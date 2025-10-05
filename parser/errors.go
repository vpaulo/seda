package parser

import (
	"fmt"
	"slices"

	"github.com/vpaulo/seda/ast"
	"github.com/vpaulo/seda/lexer"
)

// ParseError represents a parsing error with context
type ParseError struct {
	Message  string
	Line     int
	Column   int
	Token    lexer.Token
	Expected []lexer.TokenType
	Actual   lexer.TokenType
	Context  string
}

func (pe ParseError) Error() string {
	if len(pe.Expected) == 0 {
		return fmt.Sprintf("Parse error at line %d, column %d: %s (got %s)",
			pe.Line, pe.Column, pe.Message, pe.Actual)
	}

	expected_str := ""
	for i, expected := range pe.Expected {
		if i > 0 {
			if i == len(pe.Expected)-1 {
				expected_str += " or "
			} else {
				expected_str += ", "
			}
		}
		expected_str += expected.String()
	}

	return fmt.Sprintf("Parse error at line %d, column %d: expected %s, got %s",
		pe.Line, pe.Column, expected_str, pe.Actual)
}

// Enhanced error reporting methods
func (parser *Parser) record_error(message string) {
	err := ParseError{
		Message: message,
		Line:    parser.current_token.Line,
		Column:  parser.current_token.Column,
		Token:   parser.current_token,
		Actual:  parser.current_token.Type,
	}
	parser.errors = append(parser.errors, err.Error())
}

func (parser *Parser) record_expected_error(expected ...lexer.TokenType) {
	err := ParseError{
		Message:  "unexpected token",
		Line:     parser.current_token.Line,
		Column:   parser.current_token.Column,
		Token:    parser.current_token,
		Expected: expected,
		Actual:   parser.current_token.Type,
	}
	parser.errors = append(parser.errors, err.Error())
}

func (parser *Parser) record_peek_error(expected lexer.TokenType) {
	err := ParseError{
		Message:  "unexpected token",
		Line:     parser.peek_token.Line,
		Column:   parser.peek_token.Column,
		Token:    parser.peek_token,
		Expected: []lexer.TokenType{expected},
		Actual:   parser.peek_token.Type,
	}
	parser.errors = append(parser.errors, err.Error())
}

// Error recovery methods
func (parser *Parser) synchronize() {
	parser.next_token()

	for parser.current_token.Type != lexer.EOF {
		if parser.current_token.Type == lexer.END {
			return
		}

		switch parser.peek_token.Type {
		case lexer.VAR, lexer.CONST, lexer.FN, lexer.STRUCT, lexer.TYPE,
			lexer.IF, lexer.FOR, lexer.CHECK:
			return
		}

		parser.next_token()
	}
}

func (parser *Parser) skip_to_end() {
	for parser.current_token.Type != lexer.END && parser.current_token.Type != lexer.EOF {
		parser.next_token()
	}
}

func (parser *Parser) skip_to_next_statement() {
	for parser.current_token.Type != lexer.EOF {
		switch parser.current_token.Type {
		case lexer.VAR, lexer.CONST, lexer.FN, lexer.STRUCT, lexer.TYPE,
			lexer.IF, lexer.FOR, lexer.CHECK, lexer.END:
			return
		}
		parser.next_token()
	}
}

// Enhanced expect_peek with better error messages
func (parser *Parser) expect_peek_with_recovery(t lexer.TokenType) bool {
	if parser.peek_token.Type == t {
		parser.next_token()
		return true
	}

	parser.record_peek_error(t)

	// Try to recover by finding the expected token nearby
	lookahead := 0
	for lookahead < 5 && parser.peek_token.Type != lexer.EOF {
		if parser.peek_token.Type == t {
			// Found the expected token, skip to it
			for lookahead > 0 {
				parser.next_token()
				lookahead--
			}
			parser.next_token()
			return true
		}
		parser.next_token()
		lookahead++
	}

	return false
}

// Improved parsing methods with error recovery
func (parser *Parser) parse_statement_with_recovery() ast.Statement {
	defer func() {
		if r := recover(); r != nil {
			parser.record_error(fmt.Sprintf("panic during parsing: %v", r))
			parser.synchronize()
		}
	}()

	return parser.parse_statement()
}

func (parser *Parser) parse_block_statement_with_recovery() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{}

	parser.next_token()

	for parser.current_token.Type != lexer.END && parser.current_token.Type != lexer.EOF {
		stmt := parser.parse_statement_with_recovery()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		// If we encounter an error, try to recover
		if len(parser.errors) > 0 {
			parser.skip_to_next_statement()
		}

		parser.next_token()
	}

	return block
}

// Context-aware error messages
func (parser *Parser) get_parsing_context() string {
	// This could be enhanced to track what we're currently parsing
	// (function, struct, if statement, etc.)
	return "unknown context"
}

func (parser *Parser) record_contextual_error(message, context string) {
	full_message := fmt.Sprintf("%s (in %s)", message, context)
	parser.record_error(full_message)
}

// Validation helpers
func (parser *Parser) validate_identifier(ident *ast.Identifier, context string) bool {
	if ident == nil || ident.Value == "" {
		parser.record_contextual_error("invalid identifier", context)
		return false
	}

	// Check for reserved words in inappropriate contexts
	if is_reserved_word(ident.Value) {
		parser.record_contextual_error(fmt.Sprintf("'%s' is a reserved word", ident.Value), context)
		return false
	}

	return true
}

func is_reserved_word(word string) bool {
	reserved := []string{
		"var", "const", "fn", "type", "module", "using", "as", "struct", "if", "else", "case",
		"for", "in", "check", "where", "end", "is", "isA", "contains", "self", "return", "break", "true", "false",
		"number", "string", "boolean", "and", "or", "not",
	}

	return slices.Contains(reserved, word)

	// for _, r := range reserved {
	// 	if word == r {
	// 		return true
	// 	}
	// }
	// return false
}

// Pretty error formatting
func (parser *Parser) FormatErrors() []string {
	if len(parser.errors) == 0 {
		return []string{}
	}

	formatted := make([]string, len(parser.errors))
	for i, err := range parser.errors {
		formatted[i] = fmt.Sprintf("  %d. %s", i+1, err)
	}

	return formatted
}

func (parser *Parser) HasErrors() bool {
	return len(parser.errors) > 0
}

func (parser *Parser) ClearErrors() {
	parser.errors = []string{}
}
