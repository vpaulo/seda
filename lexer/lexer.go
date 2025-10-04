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

// readChar reads the next character and advances position
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
