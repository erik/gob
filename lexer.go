package main

import (
	"container/list"
	"fmt"
	"io"
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

func NewLexer(name string, input io.Reader) *Lexer {
	lex := &Lexer{
		name:      name,
		lookahead: list.New(),
	}

	lex.scanner.Init(input)
	lex.scanner.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanStrings

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
		switch tok.value {
		case "auto", "case", "else", "extrn", "goto", "if",
			"switch", "while":
			tok.kind = tkKeyword

		}

	case scanner.Int:
		tok.kind = tkNumber

	case scanner.String:
		tok.kind = tkString

	case '{':
		tok.kind = tkOpenBrace

	case '}':
		tok.kind = tkCloseBrace

	case '(':
		tok.kind = tkOpenParen

	case ')':
		tok.kind = tkCloseParen

	case ';':
		tok.kind = tkSemicolon

	case ',':
		tok.kind = tkComma

		// XXX: incomplete list
	case '+', '-', '*', '/', '%', '=', '&':
		tok.kind = tkOperator

	default:
		tok.kind = tkError
		tok.value = fmt.Sprintf("unexpected character: %c", scan)
	}

	tok.end = lex.scanner.Pos()

	return tok
}
