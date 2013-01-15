package main

import (
	"fmt"
	"testing"
)

var nodefmt = []struct {
	node Node
	str  string
}{
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

	// BlockNode
	{BlockNode{[]Node{IntegerNode{"1"}, IntegerNode{"2"},
		IntegerNode{"3"}}},
		"{\n\t1\n\t2\n\t3\n}"},

	// ExternVarInitNode
	{ExternVarInitNode{"var", IntegerNode{"2"}}, "var 2;"},

	// ExternVarDeclNode
	{ExternVarDeclNode{[]string{"a", "b", "c"}}, "extrn a, b, c;"},

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
