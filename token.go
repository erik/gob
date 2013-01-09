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
)

type Token struct {
	kind       TokenType
	value      string
	start, end scanner.Position
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
	}
	return "UnknownType"

}

func (t Token) String() string {
	return t.kind.String() + ": " + t.value
}
