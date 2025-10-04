package query

type Parser struct {
	tokens []Token
	pos    int
}

func New(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return Token{Typ: EOF}
	}

	return p.tokens[p.pos]
}
