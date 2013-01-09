package main

import (
	"fmt"
	"strings"
)

func main() {
	lex := NewLexer("asdf", strings.NewReader("123 \"a string\" a_bc_1 123"))
	for tok := lex.NextToken(); tok.kind != tkEof; tok = lex.NextToken() {
		fmt.Println(tok)
	}
}
