package parse

import (
	"fmt"
	"testing"
)

// TODO: this is an incomplete set of tests
var nodefmt = []struct {
	node Node
	str  string
	expr bool
}{
	// ArrayAccessNode
	{ArrayAccessNode{IdentNode{"abc"}, IntegerNode{"2"}}, "abc[2]", true},

	// BinaryNode
	{BinaryNode{IdentNode{"a"}, "==", IdentNode{"b"}}, "a == b", true},

	// IdentNode
	{IdentNode{"abcd"}, "abcd", true},

	// IfNode
	{IfNode{cond: BinaryNode{IdentNode{"a"}, "<", IdentNode{"b"}},
		body: StatementNode{FunctionCallNode{IdentNode{"do_this"},
			[]Node{}}},
		hasElse: false},
		"if(a < b) do_this();",
		false},
	{IfNode{cond: BinaryNode{IdentNode{"a"}, "<", IdentNode{"b"}},
		body: StatementNode{FunctionCallNode{IdentNode{"do_this"},
			[]Node{}}},
		hasElse: true,
		elseBody: StatementNode{FunctionCallNode{IdentNode{"do_that"},
			[]Node{}}}},
		"if(a < b) do_this(); else do_that();",
		false},

	// IntegerNode
	{IntegerNode{"1234567890"}, "1234567890", true},

	// CharacterNode
	{CharacterNode{""}, "''", true},
	{CharacterNode{"1"}, "'1'", true},
	{CharacterNode{"1234"}, "'1234'", true},

	// FunctionNode
	{FunctionNode{"fn", []string{"a", "b", "c"}, BlockNode{}},
		"fn(a, b, c) {\n}", false},
	{FunctionNode{"fn", []string{}, BlockNode{}}, "fn() {\n}", false},

	// FunctionCallNode
	{FunctionCallNode{IdentNode{"fn"}, []Node{IntegerNode{"1"},
		CharacterNode{"123"}}},
		"fn(1, '123')", true},

	// BlockNode
	{BlockNode{[]Node{IntegerNode{"1"}, IntegerNode{"2"},
		IntegerNode{"3"}}},
		"{\n\t1\n\t2\n\t3\n}", false},

	// ExternVarInitNode
	{ExternVarInitNode{"var", IntegerNode{"2"}}, "var 2;", false},

	// ExternVecInitNode
	{ExternVecInitNode{"var", "2", []Node{IntegerNode{"2"}}}, "var [2] 2;", false},
	{ExternVecInitNode{"var", "2", []Node{IntegerNode{"2"}, IntegerNode{"3"}}},
		"var [2] 2, 3;", false},

	// ExternVarDeclNode
	{ExternVarDeclNode{[]string{"a", "b", "c"}}, "extrn a, b, c;", false},

	// StatementNode
	{StatementNode{IntegerNode{"1"}}, "1;", false},

	// UnaryNode
	{UnaryNode{"++", IntegerNode{"1"}, false}, "++1", true},
	{UnaryNode{"++", IntegerNode{"1"}, true}, "1++", true},

	// VarDeclNode
	{VarDeclNode{[]VarDecl{{"a", false, ""},
		{"b", true, "12"},
		{"c", false, ""}}},
		"auto a, b[12], c;", false},

	// WhileNode
	{WhileNode{BinaryNode{IdentNode{"a"}, ">", IdentNode{"b"}},
		StatementNode{BinaryNode{IdentNode{"a"}, "=",
			BinaryNode{IdentNode{"a"}, "-",
				IdentNode{"b"}}}}},
		"while(a > b) a = a - b;", false},
}

func TestNodeString(t *testing.T) {
	for _, test := range nodefmt {
		str := fmt.Sprintf("%v", test.node)
		if str != test.str {
			t.Errorf("expected <%s>, got <%s>", test.str, str)
		}

		if test.expr != IsExpr(test.node) {
			t.Errorf("%v: expected expr = %v", test.node, !test.expr)
		}

	}
}
