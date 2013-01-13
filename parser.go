package main

import (
	"fmt"
	"io"
)

type ParseError struct {
	tok Token
	msg string
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("Parse error on line %d, at token: %s: %s",
		p.tok.start.Line, p.tok.String(), p.msg)
}

func NewParseError(tok Token, msg string) error {
	return &ParseError{tok, msg}
}

type Parser struct {
	lex   *Lexer
	token Token
	nodes []Node
}

func NewParser(name string, input io.Reader) *Parser {
	parse := &Parser{
		lex:   NewLexer(name, input),
		nodes: make([]Node, 0, 10),
	}

	return parse
}

func (p *Parser) Parse() error {
	tok, _ := p.lex.NextToken()
	return NewParseError(tok, "Parser not implemented")
}

func (p *Parser) nextToken() (Token, error) {
	tok, err := p.lex.NextToken()

	if err != nil {
		return tok, err
	}

	p.token = tok
	return tok, nil
}

func (p *Parser) expectType(t TokenType) {

}

func (p *Parser) expect(t TokenType, str string) {

}

func (p *Parser) accept(t TokenType) (bool, error) {
	return false, nil
}
