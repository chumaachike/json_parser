package main

import (
	"fmt"
	"os"

	"flag"
	"io"

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

	validate := flag.Bool("validate", false, "checks if json is valid")
	flag.Parse()
	var input []byte
	var err error
	if flag.NArg() > 0 {
		input, err = os.ReadFile(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	} else {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}

	engine := NewEngine(string(input))

	if *validate {
		_, err := engine.parser.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid json: %v\n", err)
		}
		fmt.Println("Json is valid")
	}

}
