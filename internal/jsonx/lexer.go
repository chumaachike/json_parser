package jsonx

import (
	"unicode"
)

type Lexer struct {
	input string
	pos   int
	line  int
	col   int
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) advance() byte {
	ch := l.peek()
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}
	return ch
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		switch l.peek() {
		case ' ', '\t', '\n', '\r':
			l.advance()
		default:
			return
		}
	}
}

func (l *Lexer) readString() Token {
	l.advance() // consume the opening "

	var value []rune
	for {
		ch := l.peek()

		if ch == 0 { // EOF
			return l.makeToken(ILLEGAL, "")
		}

		if ch == '"' { // closing quote
			l.advance() // consume it
			break
		}

		if ch == '\n' { // strings can’t span lines in JSON
			return l.makeToken(ILLEGAL, "")
		}

		value = append(value, rune(ch))
		l.advance()
	}

	return l.makeToken(STRING, string(value))
}

func (l *Lexer) makeToken(typ TokenType, val string) Token {
	return Token{Typ: typ, Value: val, Line: l.line, Col: l.col}
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0, line: 0, col: 0}
}

func (l *Lexer) readNumber() Token {
	start := l.pos

	// Optional leading minus
	if l.peek() == '-' {
		l.advance()
	}

	// Integer part
	if l.peek() == '0' {
		l.advance()
	} else if unicode.IsDigit(rune(l.peek())) {
		for unicode.IsDigit(rune(l.peek())) {
			l.advance()
		}
	}

	// Fractional part
	if l.peek() == '.' {
		l.advance()
		if !unicode.IsDigit(rune(l.peek())) {
			// Invalid: must have digit after '.'
			return l.makeToken(ILLEGAL, l.input[start:l.pos])
		}
		for unicode.IsDigit(rune(l.peek())) {
			l.advance()
		}
	}

	// Exponent part
	if l.peek() == 'e' || l.peek() == 'E' {
		l.advance()
		if l.peek() == '+' || l.peek() == '-' {
			l.advance()
		}
		if !unicode.IsDigit(rune(l.peek())) {
			// Invalid: must have digit after e/E
			return l.makeToken(ILLEGAL, l.input[start:l.pos])
		}
		for unicode.IsDigit(rune(l.peek())) {
			l.advance()
		}
	}

	return l.makeToken(NUMBER, l.input[start:l.pos])
}

func isIdentStart(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		b == '_'
}

func isIdentContinue(b byte) bool {
	return isIdentStart(b) || (b >= '0' && b <= '9')
}

func (l *Lexer) readLiteral() Token {
	if l.pos >= len(l.input) || !isIdentStart(l.peek()) {
		return l.makeToken(ILLEGAL, "unexpected literal start")
	}

	start := l.pos
	// Consume an identifier: [A-Za-z_][A-Za-z0-9_]*
	for l.pos < len(l.input) && isIdentContinue(l.peek()) {
		l.advance()
	}

	lexeme := l.input[start:l.pos]

	switch lexeme {
	case "true":
		return l.makeToken(TRUE, lexeme)
	case "false":
		return l.makeToken(FALSE, lexeme)
	case "null":
		return l.makeToken(NULL, lexeme)
	default:
		// Not a reserved literal — treat as an identifier/keyword candidate.
		return l.makeToken(ILLEGAL, "unknown literal")
	}
}

func (l *Lexer) LexAll() []Token {
	tokens := []Token{}

	for {
		l.skipWhitespace()
		ch := l.peek()
		if ch == 0 {
			tokens = append(tokens, l.makeToken(EOF, ""))
			break
		}

		switch ch {
		case '{':
			l.advance()
			tokens = append(tokens, l.makeToken(LBRACE, "{"))
		case '}':
			l.advance()
			tokens = append(tokens, l.makeToken(RBRACE, "}"))
		case '[':
			l.advance()
			tokens = append(tokens, l.makeToken(LBRACKET, "["))
		case ']':
			l.advance()
			tokens = append(tokens, l.makeToken(RBRACKET, "]"))
		case ':':
			l.advance()
			tokens = append(tokens, l.makeToken(COLON, ":"))
		case ',':
			l.advance()
			tokens = append(tokens, l.makeToken(COMMA, ","))
		case '"':
			tokens = append(tokens, l.readString())
		default:
			// numbers / literals; note: char is a byte—casting to rune is fine for ASCII tests
			if unicode.IsDigit(rune(ch)) || ch == '-' {
				tokens = append(tokens, l.readNumber())
			} else if ch == 't' || ch == 'f' || ch == 'n' {
				tokens = append(tokens, l.readLiteral())
			} else {
				tokens = append(tokens, l.makeToken(ILLEGAL, "unknown type"))
				// Avoid infinite loop on unexpected byte
				l.advance()
			}
		}
	}
	return tokens
}
