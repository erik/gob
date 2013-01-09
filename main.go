package main

import (
	"fmt"
	"strings"
)

func main() {
	lex := NewLexer("asdf", strings.NewReader("123 \"a string\" a_bc_1 123"))
	for tok, err := lex.NextToken(); tok.kind != tkEof; tok, err = lex.NextToken() {
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(tok)
	}
}
