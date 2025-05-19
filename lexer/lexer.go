package lexer

import "github.com/Jitesh117/brainrotLang-interpreter/token"

type Lexer struct {
	input        string // the input string being analyzed
	position     int    // current position in input (points to current char)
	readPosition int    // current reading position in input (after current char)
	ch           byte   // current char under examination
}

// New initializes a new Lexer isntance with the provided input string.
// It reads the first character to prime the lexer for tokenization.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // Initialise the lexer by reading the first characcter
	return l
}

// NextToken examines the current characters and returns the next token.
// It handles different types of tokens including operators, delimiters, identifers, and literals

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace() // Skip any whitespace to find the next meaningful character

	switch l.ch {
	case '=':
		// Handle '==' for equality comparison
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			// Handle '=' for assignment
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		// Handle '!=' for inequality comparison
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			// Handle '!' for negation
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '"':
		// Handle string literals
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case 0:
		// Handle end of input
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// Handle identifiers and numeric literals
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			// Handle unknown or illegal characters
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar() // Move to the next character for further analysis
	return tok
}

// skipWhitespace advances the lexer's position past an ywhitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readChar reads the next character from the input and advances the lexer's positions.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // End of input; set current character to NUL (0)
	} else {
		l.ch = l.input[l.readPosition] // Read the next character from input
	}
	l.position = l.readPosition // Update current position
	l.readPosition += 1         // Advance reading position
}

// peekChar looks ahead at the next character without advancing the lexer's position.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0 // End of input; return NUL (0)
	} else {
		return l.input[l.readPosition] // Return the next character
	}
}

// readIdentifier reads an identifier or keyword from the input.
// Identifiers start with a letter or underscore and may contain additional letters, digits, or underscores.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a numeric literal (integer) from the input.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readString reads a string literal from the input, including the quotes.
func (l *Lexer) readString() string {
	position := l.position + 1 // Skip the opening quote
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 { // End of string or input
			break
		}
	}
	return l.input[position:l.position] // Return the string contents, excluding quotes
}

// isLetter checks if a character is a letter or underscore, used to identify the start of identifiers.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit checks if a character is a digit, used to identify numeric literals.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// newToken creates a new token with the given type and literal value.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
