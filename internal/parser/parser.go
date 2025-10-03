package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/chumaachike/json_parser/pkg/token"
)

type Parser struct {
	tokens []token.Token
	pos    int
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) peek() token.Token {
	if p.pos >= len(p.tokens) {
		return token.Token{Typ: token.EOF}
	}

	return p.tokens[p.pos]
}

func (p *Parser) advance() token.Token {
	tok := p.peek()
	p.pos++
	return tok
}

func (p *Parser) expect(tt token.TokenType) (token.Token, error) {
	tok := p.peek()

	if tok.Typ != tt {
		return token.Token{}, errors.New("invalid token type")
	}
	p.pos++
	return tok, nil
}

func (p *Parser) parseValue() (interface{}, error) {
	tok := p.peek()

	switch tok.Typ {
	case token.STRING:
		p.advance()
		return tok.Value, nil

	case token.NUMBER:
		p.advance()
		val, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", tok.Value, err)
		}
		return val, nil

	case token.TRUE:
		p.advance()
		return true, nil

	case token.FALSE:
		p.advance()
		return false, nil

	case token.NULL:
		p.advance()
		return nil, nil

	case token.LBRACE:
		return p.parseObject()

	case token.LBRACKET:
		return p.parseArray()

	default:
		return nil, fmt.Errorf("illegal token type: %v", tok.Typ)
	}
}

func (p *Parser) parseObject() (map[string]any, error) {
	obj := map[string]any{}

	if _, err := p.expect(token.LBRACE); err != nil {
		return nil, err
	}

	if p.peek().Typ == token.RBRACE {
		p.advance()
		return obj, nil
	}

	for {
		keyTok, err := p.expect(token.STRING)
		if err != nil {
			return nil, fmt.Errorf("expected string key: %w", err)

		}

		key := keyTok.Value

		if _, err := p.expect(token.COLON); err != nil {
			return nil, fmt.Errorf("expected ':' after key %q: %w", key, err)
		}

		value, err := p.parseValue()
		if err != nil {
			return nil, fmt.Errorf("error parsing value for key %q: %w", key, err)
		}
		obj[key] = value

		switch p.peek().Typ {
		case token.COMMA:
			p.advance()
			continue
		case token.RBRACE:
			p.advance()
			return obj, nil
		default:
			return nil, fmt.Errorf("unexpected token in object: %v", p.peek())
		}
	}

}

func (p *Parser) parseArray() ([]any, error) {
	arr := []any{}

	if _, err := p.expect(token.LBRACKET); err != nil {
		return nil, err
	}

	if p.peek().Typ == token.RBRACKET {
		p.advance()
		return arr, nil
	}

	for {
		val, err := p.parseValue()
		if err != nil {
			return nil, fmt.Errorf("error parsing element in array: %w", err)
		}
		arr = append(arr, val)

		switch p.peek().Typ {
		case token.COMMA:
			p.advance()
			continue
		case token.RBRACKET:
			p.advance()
			return arr, nil
		default:
			return nil, fmt.Errorf("unexpected token in array: %v", p.peek())
		}

	}

}

func (p *Parser) Parse() (any, error) {
	value, err := p.parseValue()
	if err != nil {
		return nil, fmt.Errorf("unable to parse")
	}
	if p.peek().Typ != token.EOF {
		return nil, fmt.Errorf("unexpected token after JSON value: %v", p.peek())
	}
	return value, nil

}
