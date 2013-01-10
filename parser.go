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
}

func NewParser(name string, input io.Reader) *Parser {
	parse := &Parser{
		lex: NewLexer(name, input),
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

func (p *Parser) expect() {
}

func (p *Parser) accept() bool {
	return false
}
