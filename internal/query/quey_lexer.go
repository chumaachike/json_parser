package query

import (
	"unicode"
)

type Lexer struct {
	input string
	pos   int
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
	return ch
}

func (l *Lexer) makeToken(typ TokenType, val string) Token {
	return Token{Typ: typ, Value: val}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	for unicode.IsDigit(rune(l.peek())) {
		l.advance()
	}
	return l.makeToken(NUMBER, l.input[start:l.pos])
}
func (l *Lexer) readString() Token {

	ch := l.peek()
	lexeme := string(ch)
	for l.pos < len(l.input) && (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
		lexeme += string(ch)
		l.advance()
	}
	return l.makeToken(IDENT, lexeme)

}

func (l *Lexer) lexall() []Token {
	tokens := []Token{}

	for {
		ch := l.peek()

		if ch == 0 {
			tokens = append(tokens, l.makeToken(EOF, ""))
			break
		}
		switch ch {
		case '.':
			l.advance()
			tokens = append(tokens, l.makeToken(DOT, "."))
		case '[':
			l.advance()
			tokens = append(tokens, l.makeToken(ILLEGAL, "["))
		case ']':
			l.advance()
			tokens = append(tokens, l.makeToken(RBRACKET, "]"))
		default:
			if unicode.IsDigit(rune(ch)) {
				tokens = append(tokens, l.readNumber())
			} else if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
				tokens = append(tokens, l.readString())
			} else {
				tokens = append(tokens, l.makeToken(ILLEGAL, "unknown type"))
				// Avoid infinite loop on unexpected byte
				l.advance()
			}
		}

	}
	return tokens
}
