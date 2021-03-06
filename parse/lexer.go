package parse

import (
	"container/list"
	"fmt"
	"io"
	"strings"
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
	"auto":    true,
	"break":   true,
	"case":    true,
	"default": true,
	"else":    true,
	"extrn":   true,
	"goto":    true,
	"if":      true,
	"return":  true,
	"switch":  true,
	"while":   true,
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

		if _, err := lex.checkEscapes(tok.value); err != nil {
			return tok.Error(), err
		}

	case scanner.Ident:
		// Variable names have one to eight ascii characters,
		// chosen from A-Z, a-z, ., _, 0-9, and start with a
		// non-digit.
		//
		// Will be ignoring "1 to 8" characters limit.
		r := lex.scanner.Peek()
		for strings.ContainsRune("_.", r) || unicode.IsLetter(r) || unicode.IsDigit(r) {

			tok.value += string(lex.scanner.Next())

			r = lex.scanner.Peek()
		}

		if keywords[tok.value] {
			tok.kind = tkKeyword
		} else {
			tok.kind = tkIdent
		}

	case '{':
		tok.kind = tkOpenBrace

	case '}':
		tok.kind = tkCloseBrace

	case '[':
		tok.kind = tkOpenBracket

	case ']':
		tok.kind = tkCloseBracket

	case '(':
		tok.kind = tkOpenParen

	case ')':
		tok.kind = tkCloseParen

	case ';':
		tok.kind = tkSemicolon

	case ':':
		tok.kind = tkColon

	case ',':
		tok.kind = tkComma

	case '?':
		tok.kind = tkTernary

	case '\'':
		tok.kind = tkCharacter
		tok.value = ""

	endstring:
		for {
			switch char := lex.scanner.Next(); char {
			case '\n', scanner.EOF:
				return tok.Error(), NewLexError(lex.scanner.Pos(),
					fmt.Sprintf("unterminated character: %s",
						tok.value))
			case '\'':
				break endstring
			default:
				tok.value += string(char)
			}
		}

		numChars, err := lex.checkEscapes(tok.value)
		if err != nil {
			return tok.Error(), err
		}

		if numChars > 4 {
			return tok.Error(), NewLexError(lex.scanner.Pos(),
				fmt.Sprintf("oversized character literal: %s",
					tok.value))
		}

	case '/':
		if lex.scanner.Peek() == '*' {
			lex.scanner.Next() // eat '*'
		endcomment:
			for {
				switch char := lex.scanner.Next(); char {
				case scanner.EOF:
					return tok.Error(),
						NewLexError(lex.scanner.Pos(),
							"unterminated comment")
				case '*':
					if lex.scanner.Peek() == '/' {
						lex.scanner.Next()
						break endcomment
					}
				}
			}

			return lex.NextToken()
		} else {
			tok.kind = tkOperator
		}

	case '*':
		if lex.scanner.Peek() == '/' {
			lex.scanner.Next() // eat '/'
			return tok.Error(), NewLexError(lex.scanner.Pos(),
				"unexpected end of comment")
		} else {
			tok.kind = tkOperator
		}

		// XXX: some other operators still unhandled
	case '=', '>', '<', '!':
		tok.kind = tkOperator
		if lex.scanner.Peek() == '=' {
			tok.value += string(lex.scanner.Next())

		}

	case '+', '-':
		tok.kind = tkOperator

		// ++ or --
		if tok.value == string(lex.scanner.Peek()) {
			lex.scanner.Next()
			tok.value += tok.value

		}

	case '%', '&', '~':
		tok.kind = tkOperator

	default:
		return tok.Error(), NewLexError(lex.scanner.Pos(),
			fmt.Sprintf("unexpected character: %c", scan))

	}

	tok.end = lex.scanner.Pos()

	return tok, nil
}

// *0	null
// *e	end-of-file
// *(	{
// *)	}
// *t	tab
// **	*
// *'	'
// *"	"
// *n	new line
func (lex *Lexer) checkEscapes(str string) (int, error) {
	escaped := ""

	if str[len(str)-1] == '*' {
		return -1, NewLexError(lex.scanner.Pos(), "invalid escape sequence")
	}

	numChars := 0

	for i := 0; i < len(str); i++ {
		if str[i] == '*' {
			switch str[i+1] {
			case '0', 'e', '(', ')', 't', '*', '\'', '"', 'n':
			default:
				return -1, NewLexError(lex.scanner.Pos(), fmt.Sprintf("invalid escape: %c", str[i+1]))
			}

			i += 1
		} else {
			escaped += string(str[i])
		}

		numChars += 1
	}

	return numChars, nil
}
