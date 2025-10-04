package lexer

// Lexer represents the lexical analyser
type Lexer struct {
	input         string
	position      int  // current position in input (points to current char)
	read_position int  // current reading position in input (after current char)
	char          byte // current char under examination
	line          int  // current line number
	column        int  // current column number
}

// New creates a new lexer instance
func New(input string) *Lexer {
	lexer := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	lexer.read_char()
	return lexer
}

// new_token creates a new token
func new_token(token_type TokenType, char string, line, column int) Token {
	return Token{Type: token_type, Literal: char, Line: line, Column: column}
}

// is_letter checks if a character is a letter or underscore
func is_letter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_' || char > 127
}

// is_digit checks if a character is a digit
func is_digit(char byte) bool {
	return '0' <= char && char <= '9'
}

// read_char reads the next character and advances position
func (lexer *Lexer) read_char() {
	if lexer.read_position >= len(lexer.input) {
		lexer.char = 0 // ASCII NUL character represents "EOF"
	} else {
		lexer.char = lexer.input[lexer.read_position]
	}
	lexer.position = lexer.read_position
	lexer.read_position++

	if lexer.char == '\n' {
		lexer.line++
		lexer.column = 0
	} else {
		lexer.column++
	}
}

// peek_char returns the next character without advancing position
func (lexer *Lexer) peek_char() byte {
	if lexer.read_position >= len(lexer.input) {
		return 0
	}
	return lexer.input[lexer.read_position]
}

// read_identifier reads an identifier (variable name, function name, etc.)
func (lexer *Lexer) read_identifier() string {
	position := lexer.position
	for is_letter(lexer.char) || is_digit(lexer.char) {
		lexer.read_char()
	}
	return lexer.input[position:lexer.position]
}

// read_number reads a number (integer or float)
func (lexer *Lexer) read_number() string {
	position := lexer.position
	for is_digit(lexer.char) {
		lexer.read_char()
	}

	// Check for decimal point (but not range operators .. or ...)
	if lexer.char == '.' && is_digit(lexer.peek_char()) {
		lexer.read_char() // consume '.'
		for is_digit(lexer.char) {
			lexer.read_char()
		}
	}

	return lexer.input[position:lexer.position]
}

// read_string reads a string literal
// Returns empty string if unclosed (check lexer.char == 0 after calling)
func (lexer *Lexer) read_string() string {
	var result []byte
	lexer.read_char() // skip opening quote

	for lexer.char != '"' && lexer.char != 0 {
		if lexer.char == '\\' {
			lexer.read_char() // consume the backslash

			// Handle escape sequences
			switch lexer.char {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '\\':
				result = append(result, '\\')
			case '"':
				result = append(result, '"')
			case '0':
				result = append(result, '\000')
			default:
				// For unknown escape sequences, include both backslash and character
				result = append(result, '\\')
				result = append(result, lexer.char)
			}
		} else {
			result = append(result, lexer.char)
		}
		lexer.read_char()
	}

	return string(result)
}

// read_comment reads a single-line comment
func (lexer *Lexer) read_comment() string {
	position := lexer.position
	for lexer.char != '\n' && lexer.char != 0 {
		lexer.read_char()
	}
	return lexer.input[position:lexer.position]
}

// read_block_comment reads a multiline comment #| ... |#
// Handles nested #| |# pairs by tracking depth
func (lexer *Lexer) read_block_comment() string {
	start_position := lexer.position
	lexer.read_char() // skip '#'
	lexer.read_char() // skip '|'

	depth := 1 // We've entered one comment block

	for depth > 0 {
		if lexer.char == 0 {
			// EOF reached - unclosed comment
			return lexer.input[start_position:lexer.position]
		}

		// Check for nested opening #|
		if lexer.char == '#' && lexer.peek_char() == '|' {
			depth++
			lexer.read_char() // skip '#'
			lexer.read_char() // skip '|'
			continue
		}

		// Check for closing |#
		if lexer.char == '|' && lexer.peek_char() == '#' {
			depth--
			lexer.read_char() // skip '|'
			lexer.read_char() // skip '#'
			if depth == 0 {
				return lexer.input[start_position:lexer.position]
			}
			continue
		}

		lexer.read_char()
	}

	return lexer.input[start_position:lexer.position]
}

// skip_whitespace skips whitespace characters (except newlines in some contexts)
func (lexer *Lexer) skip_whitespace() {
	for lexer.char == ' ' || lexer.char == '\t' || lexer.char == '\n' || lexer.char == '\r' {
		lexer.read_char()
	}
}

// NextToken scans the input and returns the next token
func (lexer *Lexer) NextToken() Token {
	var tok Token

	lexer.skip_whitespace()

	switch lexer.char {
	case '=':
		if lexer.peek_char() == '=' {
			char := lexer.char
			lexer.read_char()
			tok = new_token(EQ, string(char)+string(lexer.char), lexer.line, lexer.column-1)
		} else if lexer.peek_char() == '>' {
			char := lexer.char
			lexer.read_char()
			tok = new_token(ARROW, string(char)+string(lexer.char), lexer.line, lexer.column-1)
		} else {
			tok = new_token(ASSIGN, string(lexer.char), lexer.line, lexer.column)
		}
	case '+':
		tok = new_token(PLUS, string(lexer.char), lexer.line, lexer.column)
	case '-':
		// TODO: I don't i'm going to use this, but i will keep it for now
		if lexer.peek_char() == '>' {
			char := lexer.char
			lexer.read_char()
			tok = new_token(TYPE_ARROW, string(char)+string(lexer.char), lexer.line, lexer.column-1)
		} else {
			tok = new_token(MINUS, string(lexer.char), lexer.line, lexer.column)
		}
	case '*':
		tok = new_token(MULTIPLY, string(lexer.char), lexer.line, lexer.column)
	case '/':
		tok = new_token(DIVIDE, string(lexer.char), lexer.line, lexer.column)
	case '%':
		tok = new_token(MODULO, string(lexer.char), lexer.line, lexer.column)
	case '^':
		tok = new_token(POWER, string(lexer.char), lexer.line, lexer.column)
	case '!':
		if lexer.peek_char() == '=' {
			char := lexer.char
			lexer.read_char()
			tok = new_token(NOT_EQ, string(char)+string(lexer.char), lexer.line, lexer.column-1)
		} else {
			tok = new_token(NOT, string(lexer.char), lexer.line, lexer.column)
		}
	case '<':
		if lexer.peek_char() == '=' {
			char := lexer.char
			lexer.read_char()
			tok = new_token(LTE, string(char)+string(lexer.char), lexer.line, lexer.column-1)
		} else {
			tok = new_token(LT, string(lexer.char), lexer.line, lexer.column)
		}
	case '>':
		if lexer.peek_char() == '=' {
			char := lexer.char
			lexer.read_char()
			tok = new_token(GTE, string(char)+string(lexer.char), lexer.line, lexer.column-1)
		} else {
			tok = new_token(GT, string(lexer.char), lexer.line, lexer.column)
		}
	case '&':
		if lexer.peek_char() == '&' {
			char := lexer.char
			lexer.read_char()
			tok = new_token(AND, string(char)+string(lexer.char), lexer.line, lexer.column-1)
		} else {
			tok = new_token(ILLEGAL, string(lexer.char), lexer.line, lexer.column)
		}
	case '|':
		if lexer.peek_char() == '|' {
			char := lexer.char
			lexer.read_char()
			tok = new_token(OR, string(char)+string(lexer.char), lexer.line, lexer.column-1)
		} else {
			tok = new_token(ILLEGAL, string(lexer.char), lexer.line, lexer.column)
		}
	case ',':
		tok = new_token(COMMA, string(lexer.char), lexer.line, lexer.column)
	case ';':
		tok = new_token(SEMICOLON, string(lexer.char), lexer.line, lexer.column)
	case ':':
		if lexer.peek_char() == ':' {
			char := lexer.char
			lexer.read_char()
			tok = new_token(DOUBLE_COLON, string(char)+string(lexer.char), lexer.line, lexer.column-1)
		} else {
			tok = new_token(COLON, string(lexer.char), lexer.line, lexer.column)
		}
	case '.':
		if lexer.peek_char() == '.' {
			char := lexer.char
			start_column := lexer.column
			lexer.read_char() // consume second '.'
			if lexer.peek_char() == '.' {
				lexer.read_char() // consume third '.'
				tok = new_token(RANGE_INCLUSIVE, string(char)+"..", lexer.line, start_column)
			} else {
				tok = new_token(RANGE, string(char)+string(lexer.char), lexer.line, start_column)
			}
		} else {
			tok = new_token(DOT, string(lexer.char), lexer.line, lexer.column)
		}
	case '(':
		tok = new_token(LPAREN, string(lexer.char), lexer.line, lexer.column)
	case ')':
		tok = new_token(RPAREN, string(lexer.char), lexer.line, lexer.column)
	case '[':
		tok = new_token(LBRACKET, string(lexer.char), lexer.line, lexer.column)
	case ']':
		tok = new_token(RBRACKET, string(lexer.char), lexer.line, lexer.column)
	case '{':
		tok = new_token(LBRACE, string(lexer.char), lexer.line, lexer.column)
	case '}':
		tok = new_token(RBRACE, string(lexer.char), lexer.line, lexer.column)
	case '"':
		start_line := lexer.line
		start_column := lexer.column
		if lexer.char == 0 && len(lexer.input) > 0 && lexer.input[len(lexer.input)-1] != '"' {
			// Unclosed string
			tok = new_token(ILLEGAL, lexer.read_string(), start_line, start_column)
		} else {
			tok = new_token(STRING, lexer.read_string(), start_line, start_column)
		}
	case '#':
		if lexer.peek_char() == '|' {
			// Multiline comment
			tok = new_token(COMMENT, lexer.read_block_comment(), lexer.line, lexer.column)
		} else {
			// Single line comment
			tok = new_token(COMMENT, lexer.read_comment(), lexer.line, lexer.column)
		}
	case 0:
		tok = new_token(EOF, "", lexer.line, lexer.column)
	default:
		if is_letter(lexer.char) {
			literal := lexer.read_identifier()
			tok = new_token(LookupIdent(literal), literal, lexer.line, lexer.column)
			return tok // early return to avoid read_char() call
		} else if is_digit(lexer.char) {
			tok = new_token(NUMBER, lexer.read_number(), lexer.line, lexer.column)
			return tok // early return to avoid read_char() call
		} else {
			tok = new_token(ILLEGAL, string(lexer.char), lexer.line, lexer.column)
		}
	}

	lexer.read_char()
	return tok
}

// GetAllTokens returns all tokens from the input (useful for testing)
func (lexer *Lexer) GetAllTokens() []Token {
	var tokens []Token
	for {
		tok := lexer.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return tokens
}
