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

func (p *Parser) Parse() (*Node, error) {
	tok, _ := p.lex.NextToken()
	return nil, NewParseError(tok, "Parser not implemented")
}

func (p *Parser) accept(t TokenType, str string) (*Token, bool) {
	var tok Token

	if p.token().kind == t {
		if str == "" || str == p.token().value {
			tok = p.token()

			// Get next token if we've matched
			if _, err := p.nextToken(); err != nil {
				// TODO: handle this
				panic(err)
			}

			return &tok, true

		}
	}

	return nil, false
}

func (p *Parser) acceptType(t TokenType) (*Token, bool) {
	return p.accept(t, "")
}

func (p *Parser) expect(t TokenType, str string) (*Token, error) {
	tok, ok := p.accept(t, str)
	if !ok {
		if str == "" {
			return nil, NewParseError(p.token(),
				fmt.Sprintf("Expected %v", t))
		} else {
			return nil, NewParseError(p.token(),
				fmt.Sprintf("Expected (%v: %v)", t, str))
		}
	}

	return tok, nil
}

func (p *Parser) expectOneOf(t ...TokenType) (TokenType, Token, error) {
	tok := p.token()

	for _, tt := range t {
		if p.token().kind == tt {
			p.nextToken()
			return tt, tok, nil
		}
	}

	types := make([]string, len(t), len(t))

	for i, tt := range t {
		types[i] = fmt.Sprintf("%s", tt)
	}

	return tkError, (&tok).Error(), NewParseError(p.token(),
		fmt.Sprintf("Expected one of: %s", strings.Join(types, ", ")))
}

func (p *Parser) expectType(t TokenType) (*Token, error) {
	return p.expect(t, "")
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

	block := BlockNode{}

	for _, ok := p.acceptType(tkCloseBrace); !ok; {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		block.nodes = append(block.nodes, *stmt)
	}

	var node Node = block
	return &node, nil
}

func (p *Parser) parseConstant() (*Node, error) {
	var node Node

	switch kind, tok, err := p.expectOneOf(tkNumber, tkCharacter, tkString); kind {
	case tkNumber:
		node = IntegerNode{tok.value}
		return &node, err
	case tkCharacter:
		node = CharacterNode{tok.value}
		return &node, err
	case tkString:
		node = StringNode{tok.value}
		return &node, err
	default:
		return nil, err
	}

	return nil, NewParseError(p.token(), "The impossible happened")
}

func (p *Parser) parseSubExpression() (*Node, error) {
	unNode := UnaryNode{oper: ""}

	// Unary prefix operator
	if tok, ok := p.acceptType(tkOperator); ok {
		// *, &, -, !, ++, --, and ~.
		switch tok.value {
		case "*", "&", "-", "!", "++", "--", "~":
			unNode = UnaryNode{oper: tok.value, postfix: false}
		default:
			return nil, NewParseError(p.token(), "invalid unary op")
		}
	}

	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	// TODO: this logic is ugly.
	if unNode.oper != "" {
		unNode.node = *expr
		*expr = unNode
	}

	// Unary postfix operator
	if p.token().kind == tkOperator {
		switch p.token().value {
		case "++", "--":
			unNode = UnaryNode{oper: p.token().value,
				node: *expr, postfix: true}
			*expr = unNode

			p.nextToken()
		}
	}

	return expr, nil
}

func (p *Parser) parseExpression() (*Node, error) {
	node, err := p.parseSubExpression()
	if err != nil {
		return nil, err
	}

	if tok, ok := p.acceptType(tkOperator); ok {
		bin := BinaryNode{left: *node, oper: tok.value}
		rhs, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		bin.right = *rhs

		*node = bin
		return node, nil
	}
	return node, nil
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

func (p *Parser) parseExternalVariableInit() (*Node, error) {
	var err error

	ident, err := p.expectType(tkIdent)

	if err != nil {
		return nil, err
	}

	retNode := ExternVarInitNode{name: ident.value}

	constant, err := p.parseConstant()
	if err != nil {
		if _, err = p.expectType(tkSemicolon); err == nil {
			// Empty declarations are zero filled
			retNode.value = IntegerNode{"0"}
			var node Node = retNode
			return &node, err
		}
	} else {
		retNode.value = *constant
	}

	if err != nil {
		return nil, err
	}

	_, err = p.expectType(tkSemicolon)

	var node Node = retNode
	return &node, err
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

	var stmt *Node

	if stmt, err = p.parseStatement(); stmt == nil || err != nil {
		return nil, err
	}

	fnNode.body = *stmt

	var node Node = fnNode
	return &node, err
}

func (p *Parser) parseIdent() (*Node, error) {
	tok, err := p.expectType(tkIdent)

	if err != nil {
		return nil, err
	}

	var node Node = IdentNode{tok.value}
	return &node, nil
}

// TODO: unfinished, untested
func (p *Parser) parseLValue() (*Node, error) {
	if _, ok := p.accept(tkOperator, "*"); ok {
		expr, err := p.parsePrimary()

		if expr == nil {
			return nil, err
		}

		var node Node = UnaryNode{oper: "*", node: *expr}
		return &node, err
	}

	// TODO: this should be more than just idents
	if id, err := p.parseIdent(); err == nil {
		arrayNode := ArrayAccessNode{array: *id}

		if _, err := p.expectType(tkOpenBracket); err != nil {
			return nil, NewParseError(p.token(), "expected lvalue")
		}

		index, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}

		arrayNode.index = *index

		if _, err := p.expectType(tkCloseBracket); err != nil {
			return nil, err
		}

		var node Node = arrayNode
		return &node, nil
	}

	return nil, NewParseError(p.token(), "expected lvalue")
}

func (p *Parser) parseParen() (*Node, error) {
	if _, err := p.expectType(tkOpenParen); err != nil {
		return nil, err
	}

	inner, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.expectType(tkCloseParen); err != nil {
		return nil, err
	}

	var node Node = ParenNode{*inner}
	return &node, nil
}

// TODO: unfinished, untested
func (p *Parser) parsePrimary() (node *Node, err error) {
	if node, err = p.parseParen(); err == nil {
	} else if node, err = p.parseConstant(); err == nil {
	} else if node, err = p.parseIdent(); err == nil {
	} else {
		return nil, NewParseError(p.token(), "expected primary expression")
	}

	// Array access
	if _, ok := p.acceptType(tkOpenBracket); ok {
		array := *node
		index, err := p.parsePrimary()

		if err != nil {
			return nil, err
		}
		if _, err := p.expectType(tkCloseBracket); err != nil {
			return nil, err
		}

		*node = ArrayAccessNode{array: array, index: *index}
		return node, nil
	}

	// Function call
	if _, ok := p.acceptType(tkOpenParen); ok {
		args := make([]Node, 0, 10)

		if p.token().kind != tkCloseParen {
			for {
				arg, err := p.parsePrimary()

				if err != nil {
					return nil, err
				}
				args = append(args, *arg)

				if _, ok := p.acceptType(tkComma); !ok {
					break
				}
			}
		}

		if _, err := p.expectType(tkCloseParen); err != nil {
			return nil, err
		}
		*node = FunctionCallNode{callable: *node, args: args}
		return node, nil
	}

	return node, nil
}

// TODO: unfinished, untested
func (p *Parser) parseRValue() (*Node, error) {
	return nil, nil
}

func (p *Parser) parseStatement() (node *Node, err error) {
	pos := p.tokIdx

	if node, err := p.parseExpression(); err != nil && p.tokIdx != pos {
		return nil, err
	} else if err == nil {
		if _, err := p.expectType(tkSemicolon); err != nil {
			return nil, err
		}
		return node, nil
	}

	if node, err := p.parseBlock(); err != nil && p.tokIdx != pos {
		return nil, err
	} else if err == nil {
		return node, nil
	}

	return nil, NewParseError(p.tokenAt(pos), "expected statement")
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

// zero or more comma separated variables
func (p *Parser) parseVariableList() ([]string, error) {
	var err error
	var vars []string = nil

	id, ok := p.acceptType(tkIdent)
	for id != nil && ok {
		vars = append(vars, id.value)

		if _, ok := p.acceptType(tkComma); !ok {
			break
		}

		if id, err = p.expectType(tkIdent); err != nil {
			return nil, err
		}
	}

	return vars, nil
}

func (p *Parser) tokenAt(idx int) Token {
	return p.tokens[idx]
}

func (p *Parser) token() Token {
	return p.tokenAt(p.tokIdx)
}
