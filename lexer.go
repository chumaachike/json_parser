package main

import (
	"unicode"
)

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]
	COLON    // :
	COMMA    // ,
	STRING
	NUMBER
	TRUE
	FALSE
	NULL
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
func (l *Lexer) makeToken(typ TokenType, val string) Token {
	return Token{typ: typ, value: val, line: l.line, col: l.col}
}

func (l *Lexer) createTokens() []Token {
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
			// readString() should advance as it consumes characters and return a STRING token
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
