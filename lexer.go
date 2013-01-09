package main

import (
	"bytes"
	"container/list"
	"fmt"
	"strings"
	"text/scanner"
)

const eof int = -1

type Position scanner.Position

type Lexer struct {
	name      string
	scanner   scanner.Scanner
	lookahead *list.List
}

const whitespace = 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' '
const noWhitespace = 0

func NewLexer(name, input string) *Lexer {
	lex := &Lexer{
		name:      name,
		lookahead: list.New(),
	}

	lex.scanner.Init(strings.NewReader(input))
	lex.scanner.Mode = scanner.ScanIdents | scanner.ScanInts

	return lex
}

func (lex *Lexer) PeekToken() Token {
	tok := lex.lexToken()
	lex.lookahead.PushBack(tok)
	return tok
}

func (lex *Lexer) NextToken() Token {
	if lex.lookahead.Front() != nil {
		node := lex.lookahead.Front()
		tok := node.Value.(Token)

		lex.lookahead.Remove(node)

		return tok
	}

	return lex.lexToken()
}

func (lex *Lexer) lexToken() Token {
	tok := Token{
		start: lex.scanner.Pos(),
	}

	scan := lex.scanner.Scan()

	tok.value = lex.scanner.TokenText()

	switch scan {
	case scanner.EOF:
		tok.kind = tkEof

	case scanner.Ident:
		tok.kind = tkIdent

	case scanner.Int:
		tok.kind = tkNumber

	case '{':
		tok.kind = tkOpenBrace

	case '}':
		tok.kind = tkCloseBrace

	case '(':
		tok.kind = tkOpenParen

	case ')':
		tok.kind = tkCloseParen

	case '"':
		tok.kind = tkString
		var eof bool
		tok.value, eof = lex.lexString()

		if eof {
			tok.kind = tkError
			tok.value = "unterminated string constant"
		}

	default:
		tok.kind = tkError
		tok.value = fmt.Sprintf("unexpected character: %c", scan)
	}

	tok.end = lex.scanner.Pos()

	return tok
}

func (lex *Lexer) lexString() (string, bool) {
	// disable skipping whitespace
	lex.scanner.Whitespace = noWhitespace
	defer func() { lex.scanner.Whitespace = whitespace }()

	var buf bytes.Buffer

	for lex.scanner.Peek() != '"' {
		next := lex.scanner.Next()
		if next == scanner.EOF {
			return buf.String(), true
		}

		buf.WriteRune(next)
	}

	// eat last "
	lex.scanner.Next()

	return buf.String(), false
}
