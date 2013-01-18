package main

import (
	"fmt"
	"testing"
)

var nodefmt = []struct {
	node Node
	str  string
}{
	// ArrayAccessNode
	{ArrayAccessNode{IdentNode{"abc"}, IntegerNode{"2"}}, "abc[2]"},

	// BinaryNode
	{BinaryNode{IdentNode{"a"}, "==", IdentNode{"b"}}, "a == b"},

	// IdentNode
	{IdentNode{"abcd"}, "abcd"},

	// IfNode
	{IfNode{cond: BinaryNode{IdentNode{"a"}, "<", IdentNode{"b"}},
		body: StatementNode{FunctionCallNode{IdentNode{"do_this"},
			[]Node{}}},
		hasElse: false},
		"if(a < b) do_this();"},
	{IfNode{cond: BinaryNode{IdentNode{"a"}, "<", IdentNode{"b"}},
		body: StatementNode{FunctionCallNode{IdentNode{"do_this"},
			[]Node{}}},
		hasElse: true,
		elseBody: StatementNode{FunctionCallNode{IdentNode{"do_that"},
			[]Node{}}}},
		"if(a < b) do_this(); else do_that();"},

	// IntegerNode
	{IntegerNode{"1234567890"}, "1234567890"},

	// CharacterNode
	{CharacterNode{""}, "''"},
	{CharacterNode{"1"}, "'1'"},
	{CharacterNode{"1234"}, "'1234'"},

	// FunctionNode
	{FunctionNode{"fn", []string{"a", "b", "c"}, BlockNode{}},
		"fn(a, b, c) {\n}"},
	{FunctionNode{"fn", []string{}, BlockNode{}}, "fn() {\n}"},

	// FunctionCallNode
	{FunctionCallNode{IdentNode{"fn"}, []Node{IntegerNode{"1"},
		CharacterNode{"123"}}},
		"fn(1, '123')"},

	// BlockNode
	{BlockNode{[]Node{IntegerNode{"1"}, IntegerNode{"2"},
		IntegerNode{"3"}}},
		"{\n\t1\n\t2\n\t3\n}"},

	// ExternVarInitNode
	{ExternVarInitNode{"var", IntegerNode{"2"}}, "var 2;"},

	// ExternVarDeclNode
	{ExternVarDeclNode{[]string{"a", "b", "c"}}, "extrn a, b, c;"},

	// StatementNode
	{StatementNode{IntegerNode{"1"}}, "1;"},

	// UnaryNode
	{UnaryNode{"++", IntegerNode{"1"}, false}, "++1"},
	{UnaryNode{"++", IntegerNode{"1"}, true}, "1++"},

	// VarDeclNode
	{VarDeclNode{[]string{"a", "b", "c"}}, "auto a, b, c;"},
}

func TestNodeString(t *testing.T) {
	for _, test := range nodefmt {
		str := fmt.Sprintf("%v", test.node)
		if str != test.str {
			t.Errorf("expected <%s>, got <%s>", test.str, str)
		}
	}
}
