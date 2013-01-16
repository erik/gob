package main

import (
	"fmt"
	"strings"
)

type NodeType int

const (
	ndError NodeType = iota
	ndArrayAccess
	ndBlock
	ndCharacter
	ndExtVarDecl
	ndExtVarInit
	ndFunction
	ndIdent
	ndInteger
	ndString
	ndUnary
	ndVarDecl
)

type Node interface {
	Type() NodeType
	String() string
}

type ArrayAccessNode struct {
	array Node
	index Node
}

func (a ArrayAccessNode) Type() NodeType { return ndArrayAccess }
func (a ArrayAccessNode) String() string {
	return fmt.Sprintf("%s[%s]", a.array, a.index)
}

// '{' node* '}'
type BlockNode struct {
	nodes []Node
}

func (b BlockNode) Type() NodeType { return ndBlock }
func (b BlockNode) String() string {
	str := "{\n"

	for _, node := range b.nodes {
		str += fmt.Sprintf("\t%v\n", node)
	}

	str += "}"
	return str
}

type CharacterNode struct {
	value string
}

func (c CharacterNode) Type() NodeType { return ndCharacter }
func (c CharacterNode) String() string { return fmt.Sprintf("'%s'", c.value) }

type ExternVarDeclNode struct {
	names []string
}

func (e ExternVarDeclNode) Type() NodeType { return ndExtVarDecl }
func (e ExternVarDeclNode) String() string {
	return fmt.Sprintf("extrn %s;", strings.Join(e.names, ", "))
}

// name value ';'
type ExternVarInitNode struct {
	name  string
	value Node
}

func (e ExternVarInitNode) Type() NodeType { return ndExtVarInit }
func (e ExternVarInitNode) String() string {
	return fmt.Sprintf("%s %v;", e.name, e.value)
}

// name '(' (var (',' var)*) ? ')' block
type FunctionNode struct {
	name   string
	params []string
	block  BlockNode
}

func (f FunctionNode) Type() NodeType { return ndFunction }
func (f FunctionNode) String() string {
	return fmt.Sprintf("%s(%s) %s",
		f.name, strings.Join(f.params, ", "), f.block)
}

type IdentNode struct {
	value string
}

func (i IdentNode) Type() NodeType { return ndIdent }
func (i IdentNode) String() string { return i.value }

type IntegerNode struct {
	value string
}

func (i IntegerNode) Type() NodeType { return ndInteger }
func (i IntegerNode) String() string { return i.value }

type StringNode struct {
	value string
}

func (s StringNode) Type() NodeType { return ndString }
func (s StringNode) String() string { return fmt.Sprintf("\"%s\"", s.value) }

type UnaryNode struct {
	oper string
	node Node
}

func (u UnaryNode) Type() NodeType { return ndUnary }
func (u UnaryNode) String() string {
	// TODO: ignores postfix
	return fmt.Sprintf("%s%v", u.oper, u.node)
}

type VarDeclNode struct {
	vars []string
}

func (v VarDeclNode) Type() NodeType { return ndVarDecl }
func (v VarDeclNode) String() string {
	return fmt.Sprintf("auto %s;", strings.Join(v.vars, ", "))
}
