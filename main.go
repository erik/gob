package main

import (
	"fmt"
)

func main() {
	lex := NewLexer("asdf", "123 \"a string\" a_bc_1 123")
	for tok := lex.NextToken(); tok.kind != tkEof; tok = lex.NextToken() {
		fmt.Println(tok)
	}
}
