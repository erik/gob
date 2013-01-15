package main

import (
	"fmt"
	"io"
	"strings"
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
	lex    *Lexer
	tokens []Token
	tokIdx int
	nodes  []Node
}

func NewParser(name string, input io.Reader) *Parser {
	parse := &Parser{
		lex:    NewLexer(name, input),
		nodes:  make([]Node, 0, 10),
		tokens: make([]Token, 0, 10),
		tokIdx: -1,
	}

	if _, err := parse.nextToken(); err != nil {
		panic(err)
	}

	return parse
}

func (p *Parser) tokenAt(idx int) Token {
	return p.tokens[idx]
}

func (p *Parser) token() Token {
	return p.tokenAt(p.tokIdx)
}

func (p *Parser) nextToken() (Token, error) {
	p.tokIdx += 1

	if p.tokIdx < len(p.tokens) {
		return p.tokens[p.tokIdx], nil
	}

	tok, err := p.lex.NextToken()
	if err != nil {
		return tok, err
	}

	p.tokens = append(p.tokens, tok)

	return tok, nil
}

func (p *Parser) parseBlock() (*Node, error) {
	if _, err := p.expectType(tkOpenBrace); err != nil {
		return nil, err
	}

	// TODO: rest of block

	if _, err := p.expectType(tkCloseBrace); err != nil {
		return nil, err
	}

	block := BlockNode{}

	var node Node = block

	return &node, nil
}

// zero or more comma separated variables
func (p *Parser) parseVariableList() ([]string, error) {

	var vars []string = nil

	id, err := p.acceptType(tkIdent)
	for id != nil && err == nil {
		vars = append(vars, id.value)

		if tok, err := p.acceptType(tkComma); tok == nil || err != nil {
			break
		}

		if id, err = p.expectType(tkIdent); err != nil {
			return nil, err
		}
	}

	return vars, nil
}

func (p *Parser) parseExternVarDecl() (*Node, error) {
	var err error

	if _, err = p.expect(tkKeyword, "extrn"); err != nil {
		return nil, err
	}

	varNode := ExternVarDeclNode{}

	if varNode.names, err = p.parseVariableList(); err != nil {
		return nil, err
	}

	if _, err = p.expectType(tkSemicolon); err != nil {
		return nil, err
	}

	if len(varNode.names) <= 0 {
		return nil, NewParseError(p.token(),
			"expected at least 1 variable in extrn declaration")
	}

	var node Node = varNode
	return &node, nil
}

func (p *Parser) parseVarDecl() (*Node, error) {
	var err error

	if _, err = p.expect(tkKeyword, "auto"); err != nil {
		return nil, err
	}

	varNode := VarDeclNode{}

	if varNode.vars, err = p.parseVariableList(); err != nil {
		return nil, err
	}

	if _, err = p.expectType(tkSemicolon); err != nil {
		return nil, err
	}

	if len(varNode.vars) <= 0 {
		return nil, NewParseError(p.token(),
			"expected at least 1 variable in auto declaration")
	}

	var node Node = varNode
	return &node, nil
}

func (p *Parser) parseFuncDeclaration() (*Node, error) {
	var err error

	id, err := p.expectType(tkIdent)

	if err != nil {
		return nil, err
	}

	fnNode := FunctionNode{name: id.value}

	if _, err = p.expectType(tkOpenParen); err != nil {
		return nil, err
	}

	if fnNode.params, err = p.parseVariableList(); err != nil {
		return nil, err
	}

	if _, err = p.expectType(tkCloseParen); err != nil {
		return nil, err
	}

	var block *Node

	if block, err = p.parseBlock(); block == nil || err != nil {
		return nil, err
	}

	fnNode.block = (*block).(BlockNode)

	var node Node = fnNode
	return &node, err
}

func (p *Parser) parseExternalVariableInit() (*Node, error) {
	var err error
	var tok *Token

	ident, err := p.expectType(tkIdent)

	if err != nil {
		return nil, err
	}

	retNode := ExternVarInitNode{name: ident.value}

	switch kind, err := p.expectOneOf(tkNumber, tkCharacter); kind {
	case tkNumber:
		tok, err = p.expectType(tkNumber)
		retNode.value = IntegerNode{tok.value}
	case tkCharacter:
		tok, err = p.expectType(tkCharacter)
		retNode.value = CharacterNode{tok.value}
	case tkError:
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	tok, err = p.expectType(tkSemicolon)

	var node Node = retNode
	return &node, err
}

// function declaration or external variable
func (p *Parser) parseTopLevel() (*Node, error) {
	if node, err := p.parseFuncDeclaration(); node != nil {
		return node, err
	} else if node, err := p.parseExternalVariableInit(); node != nil {
		return node, err
	}

	return nil, NewParseError(p.token(), "expected top level decl")
}

func (p *Parser) Parse() error {
	tok, _ := p.lex.NextToken()
	return NewParseError(tok, "Parser not implemented")
}

// doesn't advance tokIdx
func (p *Parser) expectOneOf(t ...TokenType) (TokenType, error) {
	for _, tt := range t {
		if p.token().kind == tt {
			return tt, nil
		}
	}

	types := make([]string, len(t), len(t))

	for i, tt := range t {
		types[i] = fmt.Sprintf("%s", tt)
	}

	return tkError, NewParseError(p.token(),
		fmt.Sprintf("Expected one of: %s", strings.Join(types, ", ")))
}

func (p *Parser) accept(t TokenType, str string) (*Token, error) {
	var tok Token
	var err error = nil

	if p.token().kind == t {
		if str == "" || str == p.token().value {
			tok = p.token()

			// Get next token if we've matched
			_, err := p.nextToken()
			return &tok, err

		}
	}

	return nil, err
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
		if str == "" {
			return nil, NewParseError(p.token(),
				fmt.Sprintf("Expected %v", t))
		} else {
			return nil, NewParseError(p.token(),
				fmt.Sprintf("Expected (%v: %v)", t, str))
		}
	}

	if err != nil {
		return nil, err
	}

	return tok, nil
}
