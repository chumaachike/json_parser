package main

import (
	"unicode"
)

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

// / Helpers — ASCII version. If you need Unicode identifiers, see note below.
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
