package parser

import (
	"fmt"

	"github.com/vpaulo/seda/ast"
	"github.com/vpaulo/seda/lexer"
)

// Operator precedence constants
const (
	_ int = iota
	LOWEST
	ASSIGNMENT  // =
	EQUALS      // ==
	LESSGREATER // > or <
	RANGE_PREC  // .. or ...
	SUM         // +
	PRODUCT     // *
	POWER       // ^
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	DOT         // obj.property
)

// precedences maps token types to their precedence
var precedences = map[lexer.TokenType]int{
	lexer.ASSIGN:          ASSIGNMENT,
	lexer.RANGE:           RANGE_PREC,
	lexer.RANGE_INCLUSIVE: RANGE_PREC,
	lexer.EQ:              EQUALS,
	lexer.NOT_EQ:          EQUALS,
	lexer.LT:              LESSGREATER,
	lexer.GT:              LESSGREATER,
	lexer.LTE:             LESSGREATER,
	lexer.GTE:             LESSGREATER,
	lexer.PLUS:            SUM,
	lexer.MINUS:           SUM,
	lexer.DIVIDE:          PRODUCT,
	lexer.MULTIPLY:        PRODUCT,
	lexer.MODULO:          PRODUCT,
	lexer.POWER:           POWER,
	lexer.LPAREN:          CALL,
	lexer.LBRACKET:        INDEX,
	lexer.DOT:             DOT,
	lexer.AND:             EQUALS,
	lexer.OR:              EQUALS,
}

type (
	prefix_parse_fn func() ast.Expression
	infix_parse_fn  func(ast.Expression) ast.Expression
)

// Parser represents the parser
type Parser struct {
	lexer *lexer.Lexer

	current_token lexer.Token
	peek_token    lexer.Token

	errors []string

	prefix_parse_fns map[lexer.TokenType]prefix_parse_fn
	infix_parse_fns  map[lexer.TokenType]infix_parse_fn
}

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  l,
		errors: []string{},
	}

	// Initialize prefix parse functions
	parser.prefix_parse_fns = make(map[lexer.TokenType]prefix_parse_fn)
	parser.register_prefix(lexer.IDENT, parser.parse_identifier)
	parser.register_prefix(lexer.NUMBER, parser.parse_number_literal)
	parser.register_prefix(lexer.STRING, parser.parse_string_literal)
	parser.register_prefix(lexer.TRUE, parser.parse_boolean_literal)
	parser.register_prefix(lexer.FALSE, parser.parse_boolean_literal)
	parser.register_prefix(lexer.MINUS, parser.parse_prefix_expression)
	parser.register_prefix(lexer.NOT, parser.parse_prefix_expression)
	parser.register_prefix(lexer.LPAREN, parser.parse_grouped_expression)
	parser.register_prefix(lexer.LBRACKET, parser.parse_array_literal)
	parser.register_prefix(lexer.LBRACE, parser.parse_map_literal)
	parser.register_prefix(lexer.SELF, parser.parse_self_expression)
	parser.register_prefix(lexer.CASE, parser.parse_case_expression)
	parser.register_prefix(lexer.FN, parser.parse_anonymous_function)

	// Initialize infix parse functions
	parser.infix_parse_fns = make(map[lexer.TokenType]infix_parse_fn)
	parser.register_infix(lexer.PLUS, parser.parse_infix_expression)
	parser.register_infix(lexer.MINUS, parser.parse_infix_expression)
	parser.register_infix(lexer.DIVIDE, parser.parse_infix_expression)
	parser.register_infix(lexer.MULTIPLY, parser.parse_infix_expression)
	parser.register_infix(lexer.MODULO, parser.parse_infix_expression)
	parser.register_infix(lexer.POWER, parser.parse_infix_expression)
	parser.register_infix(lexer.EQ, parser.parse_infix_expression)
	parser.register_infix(lexer.NOT_EQ, parser.parse_infix_expression)
	parser.register_infix(lexer.LT, parser.parse_infix_expression)
	parser.register_infix(lexer.GT, parser.parse_infix_expression)
	parser.register_infix(lexer.LTE, parser.parse_infix_expression)
	parser.register_infix(lexer.GTE, parser.parse_infix_expression)
	parser.register_infix(lexer.AND, parser.parse_infix_expression)
	parser.register_infix(lexer.OR, parser.parse_infix_expression)
	parser.register_infix(lexer.LPAREN, parser.parse_call_expression)
	parser.register_infix(lexer.LBRACKET, parser.parse_index_expression)
	parser.register_infix(lexer.DOT, parser.parse_dot_expression)
	parser.register_infix(lexer.ASSIGN, parser.parse_assignment_expression)
	parser.register_infix(lexer.RANGE, parser.parse_range_expression)
	parser.register_infix(lexer.RANGE_INCLUSIVE, parser.parse_range_expression)

	// Read two tokens, so current_token and peek_token are both set
	parser.next_token()
	parser.next_token()

	// Skip any initial comments
	for parser.current_token.Type == lexer.COMMENT {
		parser.next_token()
	}

	return parser
}

// ParseProgram parses the entire program
func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.current_token.Type != lexer.EOF {
		// Skip semicolons
		if parser.current_token.Type == lexer.SEMICOLON {
			parser.next_token()
			continue
		}

		stmt := parser.parse_statement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		parser.next_token()

		// Skip optional semicolon after statement
		if parser.current_token.Type == lexer.SEMICOLON {
			parser.next_token()
		}
	}

	return program
}

// Errors returns parsing errors
func (parser *Parser) Errors() []string {
	return parser.errors
}

// register_prefix registers a prefix parse function
func (parser *Parser) register_prefix(token_type lexer.TokenType, fn prefix_parse_fn) {
	parser.prefix_parse_fns[token_type] = fn
}

// register_infix registers an infix parse function
func (parser *Parser) register_infix(token_type lexer.TokenType, fn infix_parse_fn) {
	parser.infix_parse_fns[token_type] = fn
}

// next_token advances to the next token
func (parser *Parser) next_token() {
	parser.current_token = parser.peek_token
	parser.peek_token = parser.lexer.NextToken()

	// Skip comment tokens
	for parser.peek_token.Type == lexer.COMMENT {
		parser.peek_token = parser.lexer.NextToken()
	}
}

// parse_statement parses a statement
func (parser *Parser) parse_statement() ast.Statement {
	switch parser.current_token.Type {
	case lexer.VAR:
		return parser.parse_var_statement(false)
	case lexer.CONST:
		return parser.parse_var_statement(true)
	case lexer.FN:
		// Try to parse as function statement; if it fails (anonymous function), parse as expression
		stmt := parser.parse_fn_statement()
		if stmt != nil {
			return stmt
		}
		return parser.parse_expression_statement()
	case lexer.STRUCT:
		return parser.parse_struct_statement()
	case lexer.TYPE:
		return parser.parse_type_statement()
	case lexer.MODULE:
		return parser.parse_module_statement()
	case lexer.USING:
		return parser.parse_using_statement()
	case lexer.IF:
		return parser.parse_if_statement()
	case lexer.CASE:
		return parser.parse_case_statement()
	case lexer.FOR:
		return parser.parse_for_statement()
	case lexer.RETURN:
		return parser.parse_return_statement()
	case lexer.BREAK:
		return parser.parse_break_statement()
	case lexer.CHECK:
		return parser.parse_check_statement()
	case lexer.COMMENT:
		// Skip comments
		return nil
	case lexer.WHERE, lexer.ELSE, lexer.END:
		// These are not statements but block terminators
		return nil
	default:
		return parser.parse_expression_statement()
	}
}

// parse_var_statement parses variable declarations
func (parser *Parser) parse_var_statement(is_constant bool) *ast.VarStatement {
	stmt := &ast.VarStatement{IsConstant: is_constant}

	if !parser.expect_peek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Value: parser.current_token.Literal}

	// Optional type annotation
	if parser.peek_token.Type == lexer.COLON {
		parser.next_token()
		parser.next_token() // Move to the type token
		stmt.Type = parser.parse_type_annotation()
	}

	if !parser.expect_peek(lexer.ASSIGN) {
		return nil
	}

	parser.next_token()
	stmt.Value = parser.parse_expression(LOWEST)

	return stmt
}

// parse_return_statement parses return statements
func (parser *Parser) parse_return_statement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{}

	// Check if there's a return value
	if parser.peek_token.Type != lexer.SEMICOLON && parser.peek_token.Type != lexer.END && parser.peek_token.Type != lexer.EOF {
		parser.next_token()
		stmt.Value = parser.parse_expression(LOWEST)
	}

	return stmt
}

// parse_break_statement parses break statements
func (parser *Parser) parse_break_statement() *ast.BreakStatement {
	return &ast.BreakStatement{}
}

// parse_fn_statement parses function declarations
func (parser *Parser) parse_fn_statement() *ast.FnStatement {
	stmt := &ast.FnStatement{}

	// Check if next token is LPAREN - if so, this is an anonymous function expression, not a statement
	if parser.peek_token.Type == lexer.LPAREN {
		return nil
	}

	if !parser.expect_peek(lexer.IDENT) {
		return nil
	}

	// Check if this is a method (Type.method syntax)
	if parser.peek_token.Type == lexer.DOT {
		// This is a method
		stmt.Receiver = &ast.TypeAnnotation{Name: parser.current_token.Literal}
		parser.next_token() // consume DOT
		if !parser.expect_peek(lexer.IDENT) {
			return nil
		}
	}

	stmt.Name = &ast.Identifier{Value: parser.current_token.Literal}

	if !parser.expect_peek(lexer.LPAREN) {
		return nil
	}

	stmt.Parameters = parser.parse_function_parameters()

	// Optional return type
	if parser.peek_token.Type == lexer.COLON {
		parser.next_token()
		parser.next_token()
		stmt.ReturnType = parser.parse_type_annotation()
	}

	if !parser.expect_peek(lexer.DOUBLE_COLON) {
		return nil
	}

	stmt.Body = parser.parse_block_statement()

	// Check for where block
	if parser.current_token.Type == lexer.WHERE {
		stmt.WhereBlock = parser.parse_where_block()
		// parse_where_block leaves us positioned at END token
	}
	// parse_block_statement and parse_where_block both leave us at END token
	// No need to call expect_peek since we should already be on END

	return stmt
}

// parse_struct_statement parses struct declarations
func (parser *Parser) parse_struct_statement() *ast.StructStatement {
	stmt := &ast.StructStatement{}

	if !parser.expect_peek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Value: parser.current_token.Literal}

	if !parser.expect_peek(lexer.DOUBLE_COLON) {
		return nil
	}

	stmt.Fields = []*ast.StructField{}

	parser.next_token()
	for parser.current_token.Type != lexer.END && parser.current_token.Type != lexer.EOF {
		if parser.current_token.Type == lexer.IDENT {
			field := &ast.StructField{}
			field.Name = &ast.Identifier{Value: parser.current_token.Literal}

			if !parser.expect_peek(lexer.COLON) {
				return nil
			}

			parser.next_token()
			field.Type = parser.parse_type_annotation()

			stmt.Fields = append(stmt.Fields, field)

			if parser.peek_token.Type == lexer.COMMA {
				parser.next_token()
			}
		}
		parser.next_token()
	}

	return stmt
}

// parse_type_statement parses type alias declarations
func (parser *Parser) parse_type_statement() *ast.TypeStatement {
	stmt := &ast.TypeStatement{}

	if !parser.expect_peek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Value: parser.current_token.Literal}

	if !parser.expect_peek(lexer.ASSIGN) {
		return nil
	}

	parser.next_token()
	stmt.Type = parser.parse_type_annotation()

	return stmt
}

// parse_module_statement parses module statements
func (parser *Parser) parse_module_statement() *ast.ModuleStatement {
	stmt := &ast.ModuleStatement{}

	if !parser.expect_peek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Value: parser.current_token.Literal}

	if !parser.expect_peek(lexer.DOUBLE_COLON) {
		return nil
	}

	stmt.Body = parser.parse_block_statement()

	return stmt
}

// parse_using_statement parses using statements for external module imports
func (parser *Parser) parse_using_statement() *ast.UsingStatement {
	stmt := &ast.UsingStatement{}

	if !parser.expect_peek(lexer.STRING) {
		return nil
	}

	stmt.Path = &ast.StringLiteral{Value: parser.current_token.Literal}

	// Check for optional "as" alias
	if parser.peek_token.Type == lexer.AS {
		parser.next_token() // consume "as"
		if !parser.expect_peek(lexer.IDENT) {
			return nil
		}
		stmt.Alias = &ast.Identifier{Value: parser.current_token.Literal}
	}

	return stmt
}

// parse_if_statement parses if statements
func (parser *Parser) parse_if_statement() *ast.IfStatement {
	stmt := &ast.IfStatement{}

	parser.next_token()
	stmt.Condition = parser.parse_expression(LOWEST)

	if !parser.expect_peek(lexer.DOUBLE_COLON) {
		return nil
	}

	stmt.ThenBlock = parser.parse_block_statement()

	// Handle else if and else clauses
	for parser.current_token.Type == lexer.ELSE {
		if parser.peek_token.Type == lexer.IF {
			// else if
			parser.next_token() // move to IF
			parser.next_token() // move to condition
			else_if := &ast.ElseIfClause{}
			else_if.Condition = parser.parse_expression(LOWEST)
			if !parser.expect_peek(lexer.DOUBLE_COLON) {
				return nil
			}
			else_if.Block = parser.parse_block_statement()
			stmt.ElseIfs = append(stmt.ElseIfs, else_if)
		} else {
			// else
			if !parser.expect_peek(lexer.DOUBLE_COLON) {
				return nil
			}
			stmt.ElseBlock = parser.parse_block_statement()
			break
		}
	}

	// parseBlockStatement leaves us at END token - no need to expectPeek

	return stmt
}

// parse_case_statement parses case statements
func (parser *Parser) parse_case_statement() *ast.CaseStatement {
	stmt := &ast.CaseStatement{}

	parser.next_token()
	stmt.Expression = parser.parse_expression(LOWEST)

	if !parser.expect_peek(lexer.DOUBLE_COLON) {
		return nil
	}

	stmt.Branches = []*ast.CaseBranch{}

	parser.next_token()
	for parser.current_token.Type != lexer.END && parser.current_token.Type != lexer.EOF {
		branch := &ast.CaseBranch{}
		branch.Pattern = parser.parse_expression(LOWEST)

		if !parser.expect_peek(lexer.ARROW) {
			return nil
		}

		parser.next_token()
		branch.Result = parser.parse_expression(LOWEST)

		stmt.Branches = append(stmt.Branches, branch)
		parser.next_token()
	}

	return stmt
}

// parse_case_expression parses case expressions (returns a value)
func (parser *Parser) parse_case_expression() ast.Expression {
	expr := &ast.CaseExpression{}

	parser.next_token()
	expr.Expression = parser.parse_expression(LOWEST)

	if !parser.expect_peek(lexer.DOUBLE_COLON) {
		return nil
	}

	expr.Branches = []*ast.CaseBranch{}

	parser.next_token()
	for parser.current_token.Type != lexer.END && parser.current_token.Type != lexer.EOF {
		branch := &ast.CaseBranch{}
		branch.Pattern = parser.parse_expression(LOWEST)

		if !parser.expect_peek(lexer.ARROW) {
			return nil
		}

		parser.next_token()
		branch.Result = parser.parse_expression(LOWEST)

		expr.Branches = append(expr.Branches, branch)
		parser.next_token()
	}

	return expr
}

// parse_for_statement parses for loops
func (parser *Parser) parse_for_statement() *ast.ForStatement {
	stmt := &ast.ForStatement{}

	if !parser.expect_peek(lexer.IDENT) {
		return nil
	}

	// Check if this is "for index, value" or just "for value"
	if parser.peek_token.Type == lexer.COMMA {
		stmt.Index = &ast.Identifier{Value: parser.current_token.Literal}
		parser.next_token()
		if !parser.expect_peek(lexer.IDENT) {
			return nil
		}
	}

	stmt.Variable = &ast.Identifier{Value: parser.current_token.Literal}

	if !parser.expect_peek(lexer.IN) {
		return nil
	}

	parser.next_token()
	stmt.Iterable = parser.parse_expression(LOWEST)

	if !parser.expect_peek(lexer.DOUBLE_COLON) {
		return nil
	}

	stmt.Body = parser.parse_block_statement()

	// parse_block_statement leaves us at END token - no need to expectPeek

	return stmt
}

// parse_check_statement parses check blocks
func (parser *Parser) parse_check_statement() *ast.CheckStatement {
	stmt := &ast.CheckStatement{}

	// Check for optional label
	if parser.peek_token.Type == lexer.STRING {
		parser.next_token()
		stmt.Label = parser.current_token.Literal
	}

	if !parser.expect_peek(lexer.DOUBLE_COLON) {
		return nil
	}

	stmt.Assertions = []*ast.Assertion{}

	parser.next_token()
	for parser.current_token.Type != lexer.END && parser.current_token.Type != lexer.EOF {
		assertion := parser.parse_assertion()
		if assertion != nil {
			stmt.Assertions = append(stmt.Assertions, assertion)
		}
		parser.next_token()
	}

	return stmt
}

// parse_expression_statement parses expression statements
func (parser *Parser) parse_expression_statement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{}
	stmt.Expression = parser.parse_expression(LOWEST)
	return stmt
}

// parse_bock_statement parses a block of statements
func (parser *Parser) parse_block_statement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{}

	parser.next_token()

	for parser.current_token.Type != lexer.END && parser.current_token.Type != lexer.EOF &&
		parser.current_token.Type != lexer.WHERE && parser.current_token.Type != lexer.ELSE {
		if parser.current_token.Type == lexer.COMMENT {
			// Skip comments
			parser.next_token()
			continue
		}
		stmt := parser.parse_statement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		parser.next_token()
	}

	// Check if we hit EOF without finding END
	if parser.current_token.Type == lexer.EOF {
		msg := fmt.Sprintf("line %d:%d: expected 'end' keyword, got EOF",
			parser.current_token.Line, parser.current_token.Column)
		parser.errors = append(parser.errors, msg)
	}

	return block
}

// parse_function_parameters parses function parameter lists
func (parser *Parser) parse_function_parameters() []*ast.Parameter {
	params := []*ast.Parameter{}

	if parser.peek_token.Type == lexer.RPAREN {
		parser.next_token()
		return params
	}

	parser.next_token()

	param := &ast.Parameter{}
	param.Name = &ast.Identifier{Value: parser.current_token.Literal}

	if parser.peek_token.Type == lexer.COLON {
		parser.next_token()
		parser.next_token()
		param.Type = parser.parse_type_annotation()
	}

	params = append(params, param)

	for parser.peek_token.Type == lexer.COMMA {
		parser.next_token()
		parser.next_token()

		param := &ast.Parameter{}
		param.Name = &ast.Identifier{Value: parser.current_token.Literal}

		if parser.peek_token.Type == lexer.COLON {
			parser.next_token()
			parser.next_token()
			param.Type = parser.parse_type_annotation()
		}

		params = append(params, param)
	}

	if !parser.expect_peek(lexer.RPAREN) {
		return nil
	}

	return params
}

// parse_type_annotation parses type annotations
func (parser *Parser) parse_type_annotation() *ast.TypeAnnotation {
	var type_name string
	switch parser.current_token.Type {
	case lexer.IDENT:
		type_name = parser.current_token.Literal
	case lexer.NUMBER_TYPE:
		type_name = "number"
	case lexer.STRING_TYPE:
		type_name = "string"
	case lexer.BOOLEAN_TYPE:
		type_name = "boolean"
	default:
		type_name = parser.current_token.Literal
	}
	ta := &ast.TypeAnnotation{Name: type_name}

	// Handle generic types like Array[String]
	if parser.peek_token.Type == lexer.LBRACKET {
		parser.next_token()
		ta.Parameters = []*ast.TypeAnnotation{}

		parser.next_token()
		ta.Parameters = append(ta.Parameters, parser.parse_type_annotation())

		for parser.peek_token.Type == lexer.COMMA {
			parser.next_token()
			parser.next_token()
			ta.Parameters = append(ta.Parameters, parser.parse_type_annotation())
		}

		if !parser.expect_peek(lexer.RBRACKET) {
			return nil
		}
	}

	return ta
}

// parse_where_block parses where test blocks
func (parser *Parser) parse_where_block() *ast.WhereBlock {
	wb := &ast.WhereBlock{}

	if !parser.expect_peek(lexer.DOUBLE_COLON) {
		return nil
	}

	wb.Assertions = []*ast.Assertion{}

	parser.next_token()
	for parser.current_token.Type != lexer.END && parser.current_token.Type != lexer.EOF {
		if parser.current_token.Type == lexer.COMMENT {
			// Skip comments
			parser.next_token()
			continue
		}
		assertion := parser.parse_assertion()
		if assertion != nil {
			wb.Assertions = append(wb.Assertions, assertion)
		}
		parser.next_token()
	}

	return wb
}

// parse_assertion parses test assertions
func (parser *Parser) parse_assertion() *ast.Assertion {
	assertion := &ast.Assertion{}
	assertion.Left = parser.parse_expression(LOWEST)

	// Check for assertion operators
	if !parser.is_assertion_operator(parser.peek_token.Type) {
		parser.peek_error(lexer.IS)
		return nil
	}

	parser.next_token()
	assertion.Operator = parser.current_token.Literal

	parser.next_token()
	assertion.Right = parser.parse_expression(LOWEST)

	return assertion
}

func (parser *Parser) is_assertion_operator(token_type lexer.TokenType) bool {
	return token_type == lexer.IS || token_type == lexer.ISA || token_type == lexer.CONTAINS
}

// parse_expression parses expressions using Pratt parsing
func (parser *Parser) parse_expression(precedence int) ast.Expression {
	prefix := parser.prefix_parse_fns[parser.current_token.Type]
	if prefix == nil {
		parser.no_prefix_parse_fn_error(parser.current_token.Type)
		return nil
	}
	left_exp := prefix()

	for parser.peek_token.Type != lexer.EOF && precedence < parser.peek_precedence() {
		infix := parser.infix_parse_fns[parser.peek_token.Type]
		if infix == nil {
			return left_exp
		}

		parser.next_token()
		left_exp = infix(left_exp)
	}

	return left_exp
}

// Prefix parsing functions
func (parser *Parser) parse_identifier() ast.Expression {
	return &ast.Identifier{Value: parser.current_token.Literal}
}

func (parser *Parser) parse_number_literal() ast.Expression {
	return &ast.NumberLiteral{Value: parser.current_token.Literal}
}

func (parser *Parser) parse_string_literal() ast.Expression {
	str_value := parser.current_token.Literal

	// Check if string contains interpolation patterns #{...}
	if !contains_interpolation(str_value) {
		return &ast.StringLiteral{Value: str_value}
	}

	// Parse interpolated string
	return parser.parse_interpolated_string(str_value)
}

// containsInterpolation checks if a string contains #{...} patterns
func contains_interpolation(s string) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '#' && s[i+1] == '{' {
			return true
		}
	}
	return false
}

// parse_interpolated_string parses a string with #{...} interpolations
func (parser *Parser) parse_interpolated_string(str_value string) ast.Expression {
	parts := []ast.Expression{}
	current := ""
	i := 0

	for i < len(str_value) {
		// Check for interpolation start
		if i < len(str_value)-1 && str_value[i] == '#' && str_value[i+1] == '{' {
			// Add current string part if not empty
			if current != "" {
				parts = append(parts, &ast.StringLiteral{Value: current})
				current = ""
			}

			// Find matching closing brace
			i += 2 // skip #{
			brace_depth := 1
			expr_start := i

			for i < len(str_value) && brace_depth > 0 {
				if str_value[i] == '{' {
					brace_depth++
				} else if str_value[i] == '}' {
					brace_depth--
				}
				if brace_depth > 0 {
					i++
				}
			}

			if brace_depth != 0 {
				// Unclosed interpolation - treat as regular string
				return &ast.StringLiteral{Value: str_value}
			}

			// Parse the expression inside #{}
			expr_str := str_value[expr_start:i]
			expr_lexer := lexer.New(expr_str)
			expr_parser := New(expr_lexer)
			expr := expr_parser.parse_expression(LOWEST)

			if expr != nil {
				parts = append(parts, expr)
			}

			i++ // skip closing }
		} else {
			current += string(str_value[i])
			i++
		}
	}

	// Add remaining string part
	if current != "" {
		parts = append(parts, &ast.StringLiteral{Value: current})
	}

	// If only one part and it's a string literal, return it directly
	if len(parts) == 1 {
		if str_lit, ok := parts[0].(*ast.StringLiteral); ok {
			return str_lit
		}
	}

	return &ast.InterpolatedString{Parts: parts}
}

func (parser *Parser) parse_boolean_literal() ast.Expression {
	return &ast.BooleanLiteral{Value: parser.current_token.Type == lexer.TRUE}
}

func (parser *Parser) parse_self_expression() ast.Expression {
	return &ast.Identifier{Value: "self"}
}

func (parser *Parser) parse_prefix_expression() ast.Expression {
	expression := &ast.PrefixExpression{
		Operator: parser.current_token.Literal,
	}

	parser.next_token()
	expression.Right = parser.parse_expression(PREFIX)

	return expression
}

func (parser *Parser) parse_grouped_expression() ast.Expression {
	parser.next_token()
	exp := parser.parse_expression(LOWEST)
	if !parser.expect_peek(lexer.RPAREN) {
		return nil
	}
	return exp
}

func (parser *Parser) parse_array_literal() ast.Expression {
	array := &ast.ArrayLiteral{}
	array.Elements = parser.parse_expression_list(lexer.RBRACKET)
	return array
}

func (parser *Parser) parse_map_literal() ast.Expression {
	map_lit := &ast.MapLiteral{}
	map_lit.Pairs = []ast.MapPair{}

	if parser.peek_token.Type == lexer.RBRACE {
		parser.next_token()
		return map_lit
	}

	parser.next_token()

	key := parser.parse_expression(LOWEST)
	if !parser.expect_peek(lexer.COLON) {
		return nil
	}

	parser.next_token()
	value := parser.parse_expression(LOWEST)
	map_lit.Pairs = append(map_lit.Pairs, ast.MapPair{Key: key, Value: value})

	for parser.peek_token.Type == lexer.COMMA {
		parser.next_token()
		parser.next_token()

		key := parser.parse_expression(LOWEST)
		if !parser.expect_peek(lexer.COLON) {
			return nil
		}

		parser.next_token()
		value := parser.parse_expression(LOWEST)
		map_lit.Pairs = append(map_lit.Pairs, ast.MapPair{Key: key, Value: value})
	}

	if !parser.expect_peek(lexer.RBRACE) {
		return nil
	}

	return map_lit
}

// Infix parsing functions
func (parser *Parser) parse_infix_expression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Left:     left,
		Operator: parser.current_token.Literal,
	}

	precedence := parser.current_precedence()
	parser.next_token()
	expression.Right = parser.parse_expression(precedence)

	return expression
}

func (parser *Parser) parse_call_expression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Function: fn}
	exp.Arguments = parser.parse_expression_list(lexer.RPAREN)
	return exp
}

func (parser *Parser) parse_index_expression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Left: left}

	parser.next_token()
	exp.Index = parser.parse_expression(LOWEST)

	if !parser.expect_peek(lexer.RBRACKET) {
		return nil
	}

	return exp
}

func (parser *Parser) parse_dot_expression(left ast.Expression) ast.Expression {
	exp := &ast.DotExpression{Left: left}

	// Allow both identifiers and keywords as property names
	parser.next_token()
	if !parser.is_identifier_like(parser.current_token.Type) {
		msg := fmt.Sprintf("line %d:%d: expected property name, got %s instead",
			parser.current_token.Line, parser.current_token.Column, parser.current_token.Type)
		parser.errors = append(parser.errors, msg)
		return nil
	}

	exp.Property = &ast.Identifier{Value: parser.current_token.Literal}
	return exp
}

// is_identifier_like returns true if the token can be used as an identifier/property name
func (parser *Parser) is_identifier_like(t lexer.TokenType) bool {
	// Allow IDENT and most keywords to be used as property names
	return t == lexer.IDENT ||
		t == lexer.FOR ||
		t == lexer.IF ||
		t == lexer.ELSE ||
		t == lexer.CASE ||
		t == lexer.RETURN ||
		t == lexer.VAR ||
		t == lexer.CONST ||
		t == lexer.FN ||
		t == lexer.STRUCT ||
		t == lexer.TYPE ||
		t == lexer.MODULE ||
		t == lexer.USING ||
		t == lexer.AS ||
		t == lexer.IN ||
		t == lexer.CHECK ||
		t == lexer.WHERE ||
		t == lexer.END ||
		t == lexer.IS ||
		t == lexer.ISA ||
		t == lexer.CONTAINS ||
		t == lexer.SELF ||
		t == lexer.TRUE ||
		t == lexer.FALSE ||
		t == lexer.NUMBER_TYPE ||
		t == lexer.STRING_TYPE ||
		t == lexer.BOOLEAN_TYPE
}

func (parser *Parser) parse_range_expression(left ast.Expression) ast.Expression {
	exp := &ast.RangeExpression{
		Start:     left,
		Inclusive: parser.current_token.Type == lexer.RANGE_INCLUSIVE,
	}

	precedence := parser.current_precedence()
	parser.next_token()
	exp.End = parser.parse_expression(precedence)

	return exp
}

func (parser *Parser) parse_assignment_expression(left ast.Expression) ast.Expression {
	exp := &ast.AssignmentExpression{Left: left}

	parser.next_token()
	exp.Value = parser.parse_expression(LOWEST)

	return exp
}

// Utility methods
func (parser *Parser) parse_expression_list(end lexer.TokenType) []ast.Expression {
	args := []ast.Expression{}

	if parser.peek_token.Type == end {
		parser.next_token()
		return args
	}

	parser.next_token()
	args = append(args, parser.parse_expression(LOWEST))

	for parser.peek_token.Type == lexer.COMMA {
		parser.next_token()
		parser.next_token()
		args = append(args, parser.parse_expression(LOWEST))
	}

	if !parser.expect_peek(end) {
		return nil
	}

	return args
}

func (parser *Parser) current_precedence() int {
	if p, ok := precedences[parser.current_token.Type]; ok {
		return p
	}
	return LOWEST
}

func (parser *Parser) peek_precedence() int {
	if p, ok := precedences[parser.peek_token.Type]; ok {
		return p
	}
	return LOWEST
}

func (parser *Parser) expect_peek(t lexer.TokenType) bool {
	if parser.peek_token.Type == t {
		parser.next_token()
		return true
	} else {
		parser.peek_error(t)
		return false
	}
}

// Error handling
func (parser *Parser) peek_error(t lexer.TokenType) {
	msg := fmt.Sprintf("line %d:%d: expected next token to be %s, got %s instead",
		parser.peek_token.Line, parser.peek_token.Column, t, parser.peek_token.Type)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) no_prefix_parse_fn_error(t lexer.TokenType) {
	msg := fmt.Sprintf("line %d:%d: no prefix parse function for %s found",
		parser.current_token.Line, parser.current_token.Column, t)
	parser.errors = append(parser.errors, msg)
}

// parse_function_literal parses function literals: (params) :: body end
func (parser *Parser) parse_function_literal() ast.Expression {
	lit := &ast.FunctionLiteral{}

	// We're already at LPAREN
	lit.Parameters = parser.parse_function_parameters()

	// Check if we have DOUBLE_COLON - if not, this isn't a function literal
	if parser.peek_token.Type != lexer.DOUBLE_COLON {
		return nil // This will cause fallback to grouped expression
	}

	parser.next_token() // consume DOUBLE_COLON
	lit.Body = parser.parse_block_statement()

	return lit
}

// parse_anonymous_function parses anonymous function expressions: fn() :: body end
func (parser *Parser) parse_anonymous_function() ast.Expression {
	lit := &ast.FunctionLiteral{}

	// We're at FN token
	if !parser.expect_peek(lexer.LPAREN) {
		return nil
	}

	lit.Parameters = parser.parse_function_parameters()

	if !parser.expect_peek(lexer.DOUBLE_COLON) {
		return nil
	}

	lit.Body = parser.parse_block_statement()

	return lit
}
