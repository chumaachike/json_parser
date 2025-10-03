package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/chumaachike/json_parser/internal/lexer"
	"github.com/chumaachike/json_parser/internal/parser"
)

type Engiine struct {
	lexer  *lexer.Lexer
	parser *parser.Parser
}

func NewEngine(input string) *Engiine {
	l := lexer.New(input)
	tokens := l.LexAll()

	p := parser.New(tokens)
	return &Engiine{lexer: l, parser: p}
}
func main() {
	reader := bufio.NewReader(os.Stdin)
	json_string, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("unable to read err : %v", err)
	}
	engine := NewEngine(string(json_string))
	json, err := engine.parser.Parse()
	if err != nil {
		log.Fatalf("unable to parse token err %v", err)
	}

	fmt.Printf("%T\n", json)

}
