package main

import (
	"strings"
	"testing"
)

// Some basic tests to assure that we are working with a somewhat sane lexer
func TestLexSanity(t *testing.T) {
	lex := NewLexer("file", strings.NewReader("a b ¿"))
	if lex.name != "file" {
		t.Errorf("Name incorrect: %v", lex.name)
	}

	tok, err := lex.PeekToken()
	if err != nil || tok.kind != tkIdent || tok.value != "a" {
		t.Errorf("PeekToken: %v", tok)
	}

	tok, err = lex.PeekToken()
	if err != nil || tok.kind != tkIdent || tok.value != "b" {
		t.Errorf("Double PeekToken: %v", tok)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkIdent || tok.value != "a" {
		t.Errorf("NextToken after peek: %v", tok)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkIdent || tok.value != "b" {
		t.Errorf("NextToken after second peek: %v", tok)
	}

	tok, err = lex.NextToken()
	if err == nil || tok.kind != tkError {
		t.Errorf("Bad input: %v", tok)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkEof {
		t.Errorf("Eof: %v", tok)
	}
}

// Test lexing a few basic types
func TestBasicLex(t *testing.T) {
	in := strings.NewReader(`
123 "a string with spaces"
an_identifier_1 auto auto_
'char' 'ch' ''`)

	lex := NewLexer("file", in)

	tok, err := lex.NextToken()
	if err != nil || tok.kind != tkNumber || tok.value != "123" {
		t.Errorf("Number: %v", tok)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkString || tok.value != "a string with spaces" {
		t.Errorf("String: %v", tok)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkIdent || tok.value != "an_identifier_1" {
		t.Errorf("Ident: %v", tok)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkKeyword || tok.value != "auto" {
		t.Errorf("Keyword: %v, %v", tok, err)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkIdent || tok.value != "auto_" {
		t.Errorf("Not keyword: %v, %v", tok, err)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkCharacter || tok.value != "char" {
		t.Errorf("Character: %v, %v", tok, err)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkCharacter || tok.value != "ch" {
		t.Errorf("Short character: %v, %v", tok, err)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkCharacter || tok.value != "" {
		t.Errorf("Empty character: %v, %v", tok, err)
	}
}

// Test operator lexing
func TestLexOp(t *testing.T) {
	lex := NewLexer("", strings.NewReader(`> = >=`))

	tok, err := lex.NextToken()
	if err != nil || tok.kind != tkOperator || tok.value != ">" {
		t.Errorf("GT: %v, %v", tok, err)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkOperator || tok.value != "=" {
		t.Errorf("EQ: %v, %v", tok, err)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkOperator || tok.value != ">=" {
		t.Errorf("GTE: %v, %v", tok, err)
	}
}

func TestComment(t *testing.T) {
	lex := NewLexer("",
		strings.NewReader(`1 /* comment * /* (no nesting) */ 2`))

	if tok, err := lex.NextToken(); err != nil || tok.value != "1" {
		t.Errorf("Comment (pre): %v, %v", tok, err)
	}

	if tok, err := lex.NextToken(); err != nil || tok.value != "2" {
		t.Errorf("Comment (post): %v, %v", tok, err)
	}
}

// Test some exceptional conditions
func TestExceptional(t *testing.T) {
	lex := NewLexer("", strings.NewReader(`"unterminated string`))

	tok, err := lex.NextToken()
	if err == nil || tok.kind != tkError {
		t.Errorf("Unterminated: %v", tok)
	}

	lex = NewLexer("", strings.NewReader(`123abc xyz`))

	tok, err = lex.NextToken()
	if err == nil || tok.kind != tkError {
		t.Errorf("Bad number: %v", tok)
	}

	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkIdent || tok.value != "xyz" {
		t.Errorf("Token after bad number: %v, %v", tok, err)
	}

	lex = NewLexer("", strings.NewReader(`123 abc`))
	tok, err = lex.NextToken()
	if err != nil || tok.kind != tkNumber {
		t.Errorf("Good number: %v, %v", tok, err)
	}

	lex = NewLexer("", strings.NewReader(`'oversizedchar' 'unterminated`))
	tok, err = lex.NextToken()
	if err == nil || tok.kind != tkError {
		t.Errorf("Oversized character: %v", tok)
	}

	tok, err = lex.NextToken()
	if err == nil || tok.kind != tkError {
		t.Errorf("Unterminated character: %v", tok)
	}

	lex = NewLexer("", strings.NewReader(`*/ /* unterminated`))
	// Because Emacs' syntax highlighter is silly: "*/"
	tok, err = lex.NextToken()
	if err == nil || tok.kind != tkError {
		t.Errorf("Unmatched end of comment: %v", tok)
	}

	tok, err = lex.NextToken()
	if err == nil || tok.kind != tkError {
		t.Errorf("Unterminated comment: %v", tok)
	}

}
