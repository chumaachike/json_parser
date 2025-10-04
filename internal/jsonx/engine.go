package jsonx

type JsonX struct {
	lexer  *Lexer
	parser *Parser
}

func New(input string) *JsonX {
	l := NewLexer(input)
	tokens := l.LexAll()
	p := NewParser(tokens)
	return &JsonX{lexer: l, parser: p}
}
