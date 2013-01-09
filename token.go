package main

import (
	"text/scanner"
)

type TokenType int

const (
	tkError TokenType = iota
	tkEof
	tkNumber
	tkIdent
	tkOpenBrace
	tkCloseBrace
	tkOpenParen
	tkCloseParen
	tkString
	tkSemicolon
	tkComma
	tkKeyword
	tkOperator // Composite type of all operators
)

type Token struct {
	kind       TokenType
	value      string
	start, end scanner.Position
}

func (t *Token) Error() Token {
	return Token{
		kind:  tkError,
		value: t.String(),
		start: t.start,
		end:   t.end,
	}
}

func (t TokenType) String() string {
	switch t {
	case tkError:
		return "ERROR"
	case tkEof:
		return "EOF"
	case tkNumber:
		return "Number"
	case tkIdent:
		return "Identifier"
	case tkOpenBrace:
		return "Open Brace"
	case tkCloseBrace:
		return "Close Brace"
	case tkOpenParen:
		return "Open Paren"
	case tkCloseParen:
		return "Close Paren"
	case tkString:
		return "String"
	case tkSemicolon:
		return "Semicolon"
	case tkComma:
		return "Comma"
	case tkKeyword:
		return "Keyword"
	case tkOperator:
		return "Operator"
	}

	return "UnknownType"

}

func (t Token) String() string {
	return t.kind.String() + ": " + t.value
}
