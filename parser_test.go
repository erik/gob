package main

import (
	"strings"
	"testing"
)

func TestParserAccept(t *testing.T) {
	parser := NewParser("name", strings.NewReader("1 abc"))

	if tok, err := parser.accept(tkNumber, "2"); tok != nil {
		t.Errorf("Accept: value incorrect: %v, %v", tok, err)
	}

	if tok, err := parser.accept(tkNumber, "1"); tok == nil {
		t.Errorf("Accept: correct: %v", err)
	}

	if tok, err := parser.accept(tkNumber, "abc"); tok != nil {
		t.Errorf("Accept: type incorrect: %v, %v", tok, err)
	}

	if tok, err := parser.accept(tkIdent, "abc"); tok == nil {
		t.Errorf("Accept: next correct: %v", err)
	}
}

func TestParserExpect(t *testing.T) {
	parser := NewParser("name", strings.NewReader("1 2 type_incorrect 3"))

	tok, err := parser.expect(tkNumber, "1")
	if tok == nil || err != nil {
		t.Errorf("Expect: %v, %v", tok, err)
	}

	tok, err = parser.expect(tkNumber, "value_incorrect")
	if tok != nil || err == nil {
		t.Errorf("Expect value incorrect: %v", tok)
	}

	tok, err = parser.expect(tkNumber, "type_incorrect")
	if tok != nil || err == nil {
		t.Errorf("Expect type incorrect: %v", tok)
	}

	tok, err = parser.expectType(tkNumber)
	if tok == nil || err != nil {
		t.Errorf("Expect type: %v", err)
	}
}

func TestParserExternalVarDecl(t *testing.T) {
	parser := NewParser("name", strings.NewReader(`
varname 1;
varname 'abcd';
`))

	node, err := parser.parseExternalVariableInit()
	if node == nil || err != nil {
		t.Errorf("Ext var number: %v", err)
	}

	node, err = parser.parseExternalVariableInit()
	if node == nil || err != nil {
		t.Errorf("Ext var character: %v", err)
	}
}

// TODO: flesh out this test
func TestParseFuncDecl(t *testing.T) {
	parser := NewParser("name", strings.NewReader(`main(a,b,c) {}`))

	node, err := parser.parseFuncDeclaration()
	if node == nil || err != nil {
		t.Errorf("Func declaration: %v", err)
	}
}

func TestParseExternDecl(t *testing.T) {
	parser := NewParser("name", strings.NewReader(`extrn a,b,c;`))

	if _, err := parser.parseExternVarDecl(); err != nil {
		t.Errorf("Extern: %v", err)
	}
}
