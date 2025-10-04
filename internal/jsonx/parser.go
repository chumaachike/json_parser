package jsonx

import (
	"errors"
	"fmt"
	"strconv"
)

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return Token{Typ: EOF}
	}

	return p.tokens[p.pos]
}

func (p *Parser) advance() Token {
	tok := p.peek()
	p.pos++
	return tok
}

func (p *Parser) expect(tt TokenType) (Token, error) {
	tok := p.peek()

	if tok.Typ != tt {
		return Token{}, errors.New("invalid token type")
	}
	p.pos++
	return tok, nil
}

func (p *Parser) parseValue() (any, error) {
	tok := p.peek()

	switch tok.Typ {
	case STRING:
		p.advance()
		return tok.Value, nil

	case NUMBER:
		p.advance()
		val, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", tok.Value, err)
		}
		return val, nil

	case TRUE:
		p.advance()
		return true, nil

	case FALSE:
		p.advance()
		return false, nil

	case NULL:
		p.advance()
		return nil, nil

	case LBRACE:
		return p.parseObject()

	case LBRACKET:
		return p.parseArray()

	default:
		return nil, fmt.Errorf("illegal token type: %v", tok.Typ)
	}
}

func (p *Parser) parseObject() (map[string]any, error) {
	obj := map[string]any{}

	if _, err := p.expect(LBRACE); err != nil {
		return nil, err
	}

	if p.peek().Typ == RBRACE {
		p.advance()
		return obj, nil
	}

	for {
		keyTok, err := p.expect(STRING)
		if err != nil {
			return nil, fmt.Errorf("expected string key: %w", err)

		}

		key := keyTok.Value

		if _, err := p.expect(COLON); err != nil {
			return nil, fmt.Errorf("expected ':' after key %q: %w", key, err)
		}

		value, err := p.parseValue()
		if err != nil {
			return nil, fmt.Errorf("error parsing value for key %q: %w", key, err)
		}
		obj[key] = value

		switch p.peek().Typ {
		case COMMA:
			p.advance()
			continue
		case RBRACE:
			p.advance()
			return obj, nil
		default:
			return nil, fmt.Errorf("unexpected token in object: %v", p.peek())
		}
	}

}

func (p *Parser) parseArray() ([]any, error) {
	arr := []any{}

	if _, err := p.expect(LBRACKET); err != nil {
		return nil, err
	}

	if p.peek().Typ == RBRACKET {
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
		case COMMA:
			p.advance()
			continue
		case RBRACKET:
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
	if p.peek().Typ != EOF {
		return nil, fmt.Errorf("unexpected token after JSON value: %v", p.peek())
	}
	return value, nil

}
