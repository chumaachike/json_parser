package jsonx

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

type Token struct {
	Typ   TokenType
	Value string
	Line  int
	Col   int
}
