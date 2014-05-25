package parse

import (
	"fmt"
	"strings"
)

type Node interface {
	String() string
}

func IsExpr(n Node) bool {
	switch n.(type) {
	case ArrayAccessNode, BinaryNode, IdentNode, IntegerNode, CharacterNode,
		FunctionCallNode, ParenNode, TernaryNode, UnaryNode:
		return true
	}
	return false
}

func IsStatement(n Node) bool {
	if IsExpr(n) {
		return false
	}

	switch n.(type) {
	case BlockNode, BreakNode, CaseNode, ExternVarDeclNode,
		ExternVarInitNode, ExternVecInitNode, FunctionNode, GotoNode,
		IfNode, LabelNode, NullNode, ReturnNode, StatementNode, SwitchNode,
		VarDeclNode, WhileNode:
		return true
	}

	return false
}

type ArrayAccessNode struct {
	array Node
	index Node
}

func (a ArrayAccessNode) String() string {
	return fmt.Sprintf("%s[%s]", a.array, a.index)
}

type BinaryNode struct {
	left  Node
	oper  string
	right Node
}

func (b BinaryNode) String() string {
	return fmt.Sprintf("%v %s %v",
		b.left, b.oper, b.right)
}

// Use parens to make precedence more apparent
func (b BinaryNode) StringWithPrecedence() string {
	var left, right string

	if bin, ok := b.left.(BinaryNode); ok {
		left = bin.StringWithPrecedence()
	} else {
		left = b.left.String()
	}

	if bin, ok := b.right.(BinaryNode); ok {
		right = bin.StringWithPrecedence()
	} else {
		right = b.right.String()
	}

	return fmt.Sprintf("(%v %s %v)",
		left, b.oper, right)
}

// '{' node* '}'
type BlockNode struct {
	Nodes []Node
}

func (b BlockNode) String() string {
	str := "{\n"

	for _, node := range b.Nodes {
		str += fmt.Sprintf("\t%v\n", node)
	}

	str += "}"
	return str
}

type BreakNode struct{}

func (b BreakNode) String() string { return "break;" }

type CharacterNode struct {
	value string
}

func (c CharacterNode) String() string { return fmt.Sprintf("'%s'", c.value) }

type ExternVarDeclNode struct {
	names []string
}

func (e ExternVarDeclNode) String() string {
	return fmt.Sprintf("extrn %s;", strings.Join(e.names, ", "))
}

// name value ';'
type ExternVarInitNode struct {
	Name  string
	Value Node
}

func (e ExternVarInitNode) String() string {
	return fmt.Sprintf("%s %v;", e.Name, e.Value)
}

// name '[' size ']' value+ ';'
type ExternVecInitNode struct {
	Name   string
	Size   int
	Values []Node
}

func (e ExternVecInitNode) String() string {
	vals := make([]string, len(e.Values), len(e.Values))

	for i, val := range e.Values {
		vals[i] = val.String()
	}

	return fmt.Sprintf("%s [%d] %s;", e.Name, e.Size,
		strings.Join(vals, ", "))
}

// name '(' (var (',' var)*) ? ')' block
type FunctionNode struct {
	Name   string
	Params []string
	Body   Node
}

func (f FunctionNode) String() string {
	return fmt.Sprintf("%s(%s) %s",
		f.Name, strings.Join(f.Params, ", "), f.Body)
}

type FunctionCallNode struct {
	callable Node
	args     []Node
}

func (f FunctionCallNode) String() string {
	args := make([]string, len(f.args), len(f.args))
	for i, arg := range f.args {
		args[i] = arg.String()
	}

	return fmt.Sprintf("%s(%s)", f.callable, strings.Join(args, ", "))
}

type GotoNode struct{ Label string }

func (g GotoNode) String() string { return fmt.Sprintf("goto %s;", g.Label) }

type IdentNode struct {
	value string
}

func (i IdentNode) String() string { return i.value }

type IfNode struct {
	Cond     Node
	Body     Node
	HasElse  bool
	ElseBody Node
}

func (i IfNode) String() string {
	var elseStr string = ""

	if i.HasElse {
		elseStr = fmt.Sprintf(" else %v", i.ElseBody)
	}

	return fmt.Sprintf("if(%v) %v%s", i.Cond, i.Body, elseStr)
}

type IntegerNode struct {
	value int
}

func (i IntegerNode) String() string { return fmt.Sprintf("%d", i.value) }

type LabelNode struct{ Name string }

func (l LabelNode) String() string { return fmt.Sprintf("%s:", l.Name) }

type NullNode struct{}

func (n NullNode) String() string { return "" }

type ParenNode struct{ node Node }

func (p ParenNode) String() string { return "(" + p.node.String() + ")" }

type ReturnNode struct{ Node Node }

func (r ReturnNode) String() string { return fmt.Sprintf("return %v;", r.Node) }

type StatementNode struct {
	Expr Node
}

func (s StatementNode) String() string { return fmt.Sprintf("%v;", s.Expr) }

type StringNode struct {
	value string
}

func (s StringNode) String() string { return fmt.Sprintf("\"%s\"", s.value) }

type CaseNode struct {
	Cond       Node
	Statements []Node
}

func (c CaseNode) String() string {
	str := fmt.Sprintf("\tcase %v:", c.Cond)

	for _, stmt := range c.Statements {
		str += fmt.Sprintf("\n\t\t%v", stmt)
	}

	return str
}

type SwitchNode struct {
	Cond        Node
	DefaultCase []Node
	Cases       []CaseNode
}

func (s SwitchNode) String() string {
	str := fmt.Sprintf("switch(%v) {", s.Cond)

	for _, cs := range s.Cases {
		str += "\n" + cs.String()
	}

	if s.DefaultCase != nil {
		str += "\ndefault:"
		for _, stmt := range s.DefaultCase {
			str += fmt.Sprintf("\n\t%v", stmt)
		}
	}

	return str
}

// Yes, I know "ternary" is no more descriptive than binary op,
// but there's only one.
type TernaryNode struct {
	cond      Node
	trueBody  Node
	falseBody Node
}

func (t TernaryNode) String() string {
	return fmt.Sprintf("(%v ? %v : %v)", t.cond, t.trueBody, t.falseBody)
}

type UnaryNode struct {
	oper    string
	node    Node
	postfix bool
}

func (u UnaryNode) String() string {
	if u.postfix {
		return fmt.Sprintf("%v%s", u.node, u.oper)
	}
	return fmt.Sprintf("%s%v", u.oper, u.node)
}

type VarDecl struct {
	Name    string
	VecDecl bool
	Size    int
}

type VarDeclNode struct {
	Vars []VarDecl
}

func (v VarDeclNode) String() string {
	decls := make([]string, 0, len(v.Vars))

	for _, decl := range v.Vars {
		var str string

		if decl.VecDecl {
			str = fmt.Sprintf("%s[%d]", decl.Name, decl.Size)
		} else {
			str = decl.Name
		}

		decls = append(decls, str)
	}

	return fmt.Sprintf("auto %s;", strings.Join(decls, ", "))
}

type WhileNode struct {
	Cond Node
	Body Node
}

func (w WhileNode) String() string {
	return fmt.Sprintf("while(%v) %v", w.Cond, w.Body)
}
