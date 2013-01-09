package main

import (
	"container/list"
	"fmt"
	"io"
	"text/scanner"
)

const eof int = -1

type Lexer struct {
	name      string
	scanner   scanner.Scanner
	lookahead *list.List
}

type LexError struct {
	pos scanner.Position
	msg string
}

func (l *LexError) Error() string {
	return fmt.Sprintf("Lex error on line: %d, character: %d: %s",
		l.pos.Line, l.pos.Column, l.msg)

}

func NewLexError(pos scanner.Position, msg string) error {
	return &LexError{pos, msg}
}

func NewLexer(name string, input io.Reader) *Lexer {
	lex := &Lexer{
		name:      name,
		lookahead: list.New(),
	}

	lex.scanner.Init(input)
	lex.scanner.Mode = scanner.ScanIdents | scanner.ScanInts |
		scanner.ScanStrings

	return lex
}

func (lex *Lexer) PeekToken() (Token, error) {
	tok, err := lex.lexToken()

	if err != nil {
		return tok.Error(), err
	}

	lex.lookahead.PushBack(tok)
	return tok, nil
}

func (lex *Lexer) NextToken() (Token, error) {
	if lex.lookahead.Front() != nil {
		node := lex.lookahead.Front()
		tok := node.Value.(Token)

		lex.lookahead.Remove(node)

		return tok, nil
	}

	return lex.lexToken()
}

func (lex *Lexer) lexToken() (Token, error) {
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
		return tok.Error(), NewLexError(lex.scanner.Pos(),
			fmt.Sprintf("unexpected character: %c", scan))
	}

	tok.end = lex.scanner.Pos()

	return tok, nil
}
