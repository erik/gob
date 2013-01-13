package main

import (
	"container/list"
	"fmt"
	"io"
	"text/scanner"
	"unicode"
)

const eof int = -1

type Lexer struct {
	name      string
	scanner   scanner.Scanner
	lookahead *list.List
}

var keywords = map[string]bool{
	"auto":   true,
	"break":  true,
	"case":   true,
	"else":   true,
	"extrn":  true,
	"goto":   true,
	"if":     true,
	"switch": true,
	"while":  true,
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

func (lex *Lexer) lexToken() (tok Token, err error) {
	tok = Token{
		start: lex.scanner.Pos(),
	}

	// Remove error handler
	defer func() { lex.scanner.Error = nil }()

	// TODO: this is probably horrible style
	defer func() { recover() }()

	errorHandle := func(s *scanner.Scanner, msg string) {
		tok = tok.Error()
		err = NewLexError(lex.scanner.Pos(), msg)

		// Panic to get ourselves out of the parent func, this is
		// probably terrible form
		panic("ScanErrorHandle")
	}

	lex.scanner.Error = errorHandle

	scan := lex.scanner.Scan()

	tok.value = lex.scanner.TokenText()

	switch scan {
	case scanner.EOF:
		tok.kind = tkEof

	case scanner.Ident:
		if keywords[tok.value] {
			tok.kind = tkKeyword
		} else {
			tok.kind = tkIdent
		}

	case scanner.Int:
		tok.kind = tkNumber
		// TODO: this isn't all inclusive
		if next := lex.scanner.Peek(); unicode.IsLetter(next) {
			lex.scanner.Scan() // run until end of token

			err = NewLexError(
				lex.scanner.Pos(),
				fmt.Sprintf("bad number: %s%s", tok.value,
					lex.scanner.TokenText()))

			return tok.Error(), err
		}

	case scanner.String:
		tok.kind = tkString
		// cut out leading/trailing "
		tok.value = tok.value[1 : len(tok.value)-1]

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

	case '\'':
		tok.kind = tkCharacter
		tok.value = ""

	Loop:
		for {
			switch char := lex.scanner.Next(); char {
			case '\n', scanner.EOF:
				return tok.Error(), NewLexError(lex.scanner.Pos(),
					fmt.Sprintf("unterminated character: %s",
						tok.value))
			case '\'':
				break Loop
			default:
				tok.value = fmt.Sprintf("%s%c", tok.value,
					char)
			}
		}

		if len(tok.value) > 4 {
			return tok.Error(), NewLexError(lex.scanner.Pos(),
				fmt.Sprintf("oversized character literal: %s",
					tok.value))
		}

		// XXX: some other operators still unhandled
	case '=', '>', '<', '!':
		tok.kind = tkOperator
		if lex.scanner.Peek() == '=' {
			tok.value = fmt.Sprintf("%s%c", tok.value,
				lex.scanner.Next())

		}

	case '+', '-', '*', '/', '%', '&':
		tok.kind = tkOperator

	default:
		return tok.Error(), NewLexError(lex.scanner.Pos(),
			fmt.Sprintf("unexpected character: %c", scan))
	}

	tok.end = lex.scanner.Pos()

	return tok, nil
}
