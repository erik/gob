package main

import (
	"fmt"
	"strings"
)

type NodeType int

const (
	ndError NodeType = iota
	ndBlock
	ndCharacter
	ndExtVarDecl
	ndInteger
	ndFunction
)

type Node interface {
	Type() NodeType
	String() string
}

// name '(' (var (',' var)*) ? ')' block
type FunctionNode struct {
	name   string
	params []string
	block  BlockNode
}

// '{' node* '}'
type BlockNode struct {
	nodes []Node
}

// name value ';'
type ExternVarDeclNode struct {
	name  string
	value Node
}

type IntegerNode struct {
	value string
}

type CharacterNode struct {
	value string
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

func (f FunctionNode) Type() NodeType { return ndFunction }
func (f FunctionNode) String() string {
	return fmt.Sprintf("%s(%s) %s",
		f.name, strings.Join(f.params, ", "), f.block)
}

func (e ExternVarDeclNode) Type() NodeType { return ndExtVarDecl }
func (e ExternVarDeclNode) String() string {
	return fmt.Sprintf("%s %v;", e.name, e.value)
}

func (i IntegerNode) Type() NodeType { return ndInteger }
func (i IntegerNode) String() string { return i.value }

func (c CharacterNode) Type() NodeType { return ndCharacter }
func (c CharacterNode) String() string { return fmt.Sprintf("'%s'", c.value) }
