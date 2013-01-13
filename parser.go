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

	tok, err := parse.lex.NextToken()

	if err != nil {
		panic(err)
	}

	parse.token = tok

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

func (p *Parser) accept(t TokenType, str string) (*Token, error) {
	var tok *Token = nil
	var err error = nil

	if p.token.kind == t {
		if str == "" || str == tok.value {
			tok = &p.token
		}
	}

	// Get next token if we've matched
	if tok != nil {
		next, e := p.lex.NextToken()
		err = e
		p.token = next
	}

	return tok, err
}

func (p *Parser) acceptType(t TokenType) (*Token, error) {
	return p.accept(t, "")
}

func (p *Parser) expectType(t TokenType) (*Token, error) {
	return p.expect(t, "")
}

func (p *Parser) expect(t TokenType, str string) (*Token, error) {
	tok, err := p.accept(t, str)

	if tok == nil {
		return nil, NewParseError(p.token,
			fmt.Sprintf("Expected %v (%v)", t, str))
	}

	if err != nil {
		return nil, err
	}

	return tok, nil
}
