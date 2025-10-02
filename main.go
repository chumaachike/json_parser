package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	json, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("unable to read err : %v", err)
	}
	lexer := &Lexer{
		input: string(json),
		pos:   0,
		line:  0,
		col:   0,
	}
	tokens := lexer.createTokens()
	fmt.Println(tokens)

}

func (t Token) String() string {
	return fmt.Sprintf(
		"{Type: %v, Lexeme: %q, Line: %d, Col: %d,}",
		t.typ, t.value, t.line, t.col,
	)
}
