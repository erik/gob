package parse

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

func TestParserExternalVarInit(t *testing.T) {
	parser := NewParser("name", strings.NewReader(`
varname 123;
varname 'abcd';
zero ;
`))

	node, err := parser.parseExternalVariableInit()
	if node == nil || err != nil {
		t.Errorf("Ext var number: %v", err)
	}

	node, err = parser.parseExternalVariableInit()
	if node == nil || err != nil {
		t.Errorf("Ext var character: %v", err)
	}

	node, err = parser.parseExternalVariableInit()
	if node == nil || err != nil {
		t.Errorf("Ext var empty: %v", err)
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

// TODO: flesh out this test
func TestParseExternDecl(t *testing.T) {
	parser := NewParser("name", strings.NewReader(`extrn a,b,c;`))

	if _, err := parser.parseExternVarDecl(); err != nil {
		t.Errorf("Extern: %v", err)
	}
}

// TODO: flesh out this test
func TestParseVarDecl(t *testing.T) {
	parser := NewParser("name", strings.NewReader(`auto a,b,c;`))

	if _, err := parser.parseVarDecl(); err != nil {
		t.Errorf("Var: %v", err)
	}
}

func TestParseParen(t *testing.T) {
	parser := NewParser("name", strings.NewReader(`(((('a'))))
((unmatched`))

	node, err := parser.parseParen()
	if err != nil || (*node).String() != "(((('a'))))" {
		t.Errorf("Nested paren: %v", err)
	}

	if node, err := parser.parseParen(); err == nil {
		t.Errorf("Unbalanced paren: %v", *node)
	}
}

func TestParsePrimary(t *testing.T) {
	parser := NewParser("name", strings.NewReader(`
((1)) 123 '123' abc "string"`))

	if _, err := parser.parsePrimary(); err != nil {
		t.Errorf("Paren primary: %v", err)
	}

	if _, err := parser.parsePrimary(); err != nil {
		t.Errorf("Number primary: %v", err)
	}

	if _, err := parser.parsePrimary(); err != nil {
		t.Errorf("Character primary: %v", err)
	}

	if _, err := parser.parsePrimary(); err != nil {
		t.Errorf("Ident primary: %v", err)
	}

	if _, err := parser.parsePrimary(); err != nil {
		t.Errorf("String primary: %v", err)
	}

	parser = NewParser("name", strings.NewReader(`
(func)(1,(ab(c)),3)
((abb(2))[23])[ab(c(d[2]))]
`))
	if _, err := parser.parsePrimary(); err != nil {
		t.Errorf("Complex func call: %v", err)
	}

	if _, err := parser.parsePrimary(); err != nil {
		t.Errorf("Complex array access: %v", err)
	}

}

func TestParseExpression(t *testing.T) {
	parser := NewParser("", strings.NewReader(`
-(!b[2]--)++
a=b+++-(--c)*4
`))

	if _, err := parser.parseExpression(); err != nil {
		t.Errorf("Expression unary: %v", err)
	}

	if _, err := parser.parseExpression(); err != nil {
		t.Errorf("Complex expression: %v", err)
	}
}

func TestParseSubExpression(t *testing.T) {
	parser := NewParser("name", strings.NewReader(`
*1 *abc *(123)
abc[1] abc[(23)]`))

	if _, err := parser.parseSubExpression(); err != nil {
		t.Errorf("Deref num: %v", err)
	}

	if _, err := parser.parseSubExpression(); err != nil {
		t.Errorf("Deref ident: %v", err)
	}

	if _, err := parser.parseSubExpression(); err != nil {
		t.Errorf("Deref primary: %v", err)
	}

	if _, err := parser.parseSubExpression(); err != nil {
		t.Errorf("ArrayAccess num: %v", err)
	}

	if _, err := parser.parseSubExpression(); err != nil {
		t.Errorf("ArrayAccess paren: %v", err)
	}
}

func TestParseBlock(t *testing.T) {
	parser := NewParser("", strings.NewReader(`{a=1;}
{}
{{{1;} 2;} 3;}`))

	if _, err := parser.parseBlock(); err != nil {
		t.Errorf("Simple block: %v", err)
	}

	if _, err := parser.parseBlock(); err != nil {
		t.Errorf("Empty block: %v", err)
	}

	if _, err := parser.parseBlock(); err != nil {
		t.Errorf("Nested block: %v", err)
	}

}

func TestParseStatement(t *testing.T) {
	parser := NewParser("", strings.NewReader(`{{1;}}
a=1+2;
if(a + b == c) statement; else other_statement;
auto a, b, c;
extrn a, b, c;
`))

	if _, err := parser.parseStatement(); err != nil {
		t.Errorf("Block statement: %v", err)
	}

	if _, err := parser.parseStatement(); err != nil {
		t.Errorf("Simple statement: %v", err)
	}

	if _, err := parser.parseStatement(); err != nil {
		t.Errorf("If statement: %v", err)
	}

	if _, err := parser.parseStatement(); err != nil {
		t.Errorf("Auto var decl statement: %v", err)
	}

	if _, err := parser.parseStatement(); err != nil {
		t.Errorf("Extern var decl statement: %v", err)
	}
}

// TODO: I'm only sort of sure about the correctness of these
func TestParseOperatorPrecedence(t *testing.T) {
	parser := NewParser("", strings.NewReader(`
a=b+c---d /* (a = (b + (c-- - d))) */
a+2*--a=b=c /* ((a + (2 * --a)) = (b = c)) */
a=b=c+d=e
`))

	node, err := parser.parseExpression()
	if err != nil {
		t.Errorf("Operator parse: %v", err)
	}

	if str := (*node).(BinaryNode).StringWithPrecedence(); str !=
		"(a = (b + (c-- - d)))" {
		t.Errorf("Bad precedence: %s", str)
	}

	node, err = parser.parseExpression()
	if err != nil {
		t.Errorf("Operator parse: %v", err)
	}

	if str := (*node).(BinaryNode).StringWithPrecedence(); str !=
		"((a + (2 * --a)) = (b = c))" {
		t.Errorf("Bad precedence: %s", str)
	}

	node, err = parser.parseExpression()
	if err != nil {
		t.Errorf("Operator parse: %v", err)
	}

	if str := (*node).(BinaryNode).StringWithPrecedence(); str !=
		"(a = (b = ((c + d) = e)))" {
		t.Errorf("Bad precedence: %s", str)
	}

}

func TestParseIf(t *testing.T) {
	parser := NewParser("", strings.NewReader(`
if (a + b < c) { do_this(); and_this(); }
if (a + b < c) do_that(); else { do_this(); do_that(); }
`))

	if _, err := parser.parseIf(); err != nil {
		t.Errorf("If with no else: %v", err)
	}

	if _, err := parser.parseIf(); err != nil {
		t.Errorf("If with else: %v", err)
	}
}