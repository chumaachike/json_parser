package query

import (
	"errors"
	"fmt"
	"strconv"
)

// Parser holds the token stream and current position
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

// Selectors represent query steps
type Selector interface {
	String() string
}

type Field struct {
	Name string
}

func (f Field) String() string { return "." + f.Name }

type Index struct {
	Pos int
}

func (i Index) String() string { return fmt.Sprintf("[%d]", i.Pos) }

// parseValue builds a chain of selectors from tokens
func (p *Parser) parseValue() ([]Selector, error) {
	selectors := []Selector{}

	for {
		tok := p.peek()
		switch tok.Typ {
		case DOT:
			p.advance() // consume "."
			ident, err := p.expect(IDENT)
			if err != nil {
				next := p.peek()
				return nil, fmt.Errorf("expected identifier after '.', found %v at pos %d", next.Typ, next.Pos)
			}
			selectors = append(selectors, Field{Name: ident.Value})

		case IDENT:
			// Support bare identifiers at root (no leading dot)
			ident := p.advance()
			selectors = append(selectors, Field{Name: ident.Value})

		case LBRACKET:
			p.advance() // consume "["
			numTok, err := p.expect(NUMBER)
			if err != nil {
				next := p.peek()
				return nil, fmt.Errorf("expected number after '[', found %v at pos %d", next.Typ, next.Pos)
			}
			if _, err = p.expect(RBRACKET); err != nil {
				next := p.peek()
				return nil, fmt.Errorf("expected closing bracket but found %v at pos %d", next.Typ, next.Pos)
			}
			n, err := strconv.Atoi(numTok.Value)
			if err != nil {
				return nil, fmt.Errorf("cannot parse value as number: %v", err)
			}
			selectors = append(selectors, Index{Pos: n})

		case EOF:
			return selectors, nil

		default:
			return nil, fmt.Errorf("unexpected token: %v at pos %d", tok.Typ, tok.Pos)
		}
	}
}

// Parse evaluates the selectors against a JSON-like structure
func (p *Parser) Parse(json any) (any, error) {
	selectors, err := p.parseValue()
	if err != nil {
		return nil, fmt.Errorf("unable to parse values: %v", err)
	}

	current := json

	for _, selector := range selectors {
		switch v := selector.(type) {
		case Field:
			// Expect a map
			m, ok := current.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("expected map when selecting field %q, got %T", v.Name, current)
			}
			val, exists := m[v.Name]
			if !exists {
				return nil, fmt.Errorf("field %q not found", v.Name)
			}
			current = val

		case Index:
			// Expect a slice/array
			arr, ok := current.([]any)
			if !ok {
				return nil, fmt.Errorf("expected array when selecting index %d, got %T", v.Pos, current)
			}
			if v.Pos < 0 || v.Pos >= len(arr) {
				return nil, fmt.Errorf("index %d out of bounds (len=%d)", v.Pos, len(arr))
			}
			current = arr[v.Pos]

		default:
			return nil, fmt.Errorf("unsupported selector type %T", selector)
		}
	}

	return current, nil
}
