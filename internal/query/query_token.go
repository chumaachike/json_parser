package query

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	DOT      // .
	IDENT    // e.g. name, friends
	LBRACKET // [
	RBRACKET // ]
	NUMBER   // array index
)

type Token struct {
	Typ   TokenType
	Value string
	Pos   int
}
