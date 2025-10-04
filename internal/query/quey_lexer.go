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

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0}
}

func (l *Lexer) makeToken(typ TokenType, val string, pos int) Token {
	return Token{Typ: typ, Value: val, Pos: pos}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	for unicode.IsDigit(rune(l.peek())) {
		l.advance()
	}
	return l.makeToken(NUMBER, l.input[start:l.pos], start)
}
func (l *Lexer) readIdent() Token {
	start := l.pos
	for l.pos < len(l.input) {
		ch := l.peek()
		if unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' {
			l.advance()
		} else {
			break
		}
	}
	return l.makeToken(IDENT, l.input[start:l.pos], start)
}

func (l *Lexer) Lex() []Token {
	tokens := []Token{}

	for {
		ch := l.peek()

		if ch == 0 {
			tokens = append(tokens, l.makeToken(EOF, "", l.pos))
			break
		}
		switch ch {
		case '.':
			start := l.pos
			l.advance()
			tokens = append(tokens, l.makeToken(DOT, ".", start))
		case '[':
			start := l.pos
			l.advance()
			tokens = append(tokens, l.makeToken(LBRACKET, "[", start))
		case ']':
			start := l.pos
			l.advance()
			tokens = append(tokens, l.makeToken(RBRACKET, "]", start))
		default:
			if unicode.IsDigit(rune(ch)) {
				tokens = append(tokens, l.readNumber())
			} else if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
				tokens = append(tokens, l.readIdent())
			} else if unicode.IsSpace(rune(ch)) {
				l.advance()
				continue
			} else {
				tokens = append(tokens, l.makeToken(ILLEGAL, string(ch), l.pos))
				// Avoid infinite loop on unexpected byte
				l.advance()
			}
		}

	}
	return tokens
}
