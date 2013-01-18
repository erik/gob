package parse

import (
	"text/scanner"
)

type OperatorBinding int

const (
	opRL OperatorBinding = iota // right to left binding
	opLR                        // left to right binding
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
	tkOpenBracket
	tkCloseBracket
	tkString
	tkSemicolon
	tkComma
	tkCharacter
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
	case tkOpenBracket:
		return "Open bracket"
	case tkCloseBracket:
		return "Close bracket"
	case tkString:
		return "String"
	case tkSemicolon:
		return "Semicolon"
	case tkComma:
		return "Comma"
	case tkCharacter:
		return "Character"
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

func OperatorPrecedence(op string) (prec int, bind OperatorBinding) {
	switch op {
	case "*", "/", "%":
		return 90, opLR
	case "+", "-":
		return 80, opLR
	case ">", "<", "<=", ">=":
		return 70, opLR
	case "==", "!=":
		return 60, opLR
	case "&":
		return 50, opLR
	case "^":
		return 40, opLR
	case "|":
		return 30, opLR
	case "?":
		return 20, opRL
	case "=", "=+", "=-", "=/", "=*": // etc
		return 10, opRL
	}

	return -1, -1
}
