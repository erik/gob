package main

import (
	"io"
)

type Parser struct {
	lex *Lexer
}

func NewParser(name string, input io.Reader) *Parser {
	parse := &Parser{
		lex: NewLexer(name, input),
	}

	return parse
}

func (p *Parser) Parse() {

}

func (p *Parser) expect() {
}

func (p *Parser) accept() bool {
	return false
}
