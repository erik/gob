package parse

import (
	"fmt"
	"strings"
)

type NodeType int

const (
	ndError NodeType = iota
	ndArrayAccess
	ndBinary
	ndBlock
	ndBreak
	ndCharacter
	ndExtVarDecl
	ndExtVarInit
	ndExtVecInit
	ndFunction
	ndFunctionCall
	ndGoto
	ndIdent
	ndIf
	ndInteger
	ndLabel
	ndNull
	ndParen
	ndReturn
	ndStatement
	ndString
	ndSwitch
	ndTernary
	ndUnary
	ndVarDecl
	ndWhile
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

type BinaryNode struct {
	left  Node
	oper  string
	right Node
}

func (b BinaryNode) Type() NodeType { return ndBinary }
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

type BreakNode struct{}

func (b BreakNode) Type() NodeType { return ndBreak }
func (b BreakNode) String() string { return "break;" }

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

// name '[' size ']' value+ ';'
type ExternVecInitNode struct {
	name   string
	size   string
	values []Node
}

func (e ExternVecInitNode) Type() NodeType { return ndExtVecInit }
func (e ExternVecInitNode) String() string {
	vals := make([]string, 0, len(e.values))

	for i, val := range e.values {
		vals[i] = val.String()
	}

	return fmt.Sprintf("%s [%s] %v;", e.name, e.size,
		strings.Join(vals, ", "))
}

// name '(' (var (',' var)*) ? ')' block
type FunctionNode struct {
	name   string
	params []string
	body   Node
}

func (f FunctionNode) Type() NodeType { return ndFunction }
func (f FunctionNode) String() string {
	return fmt.Sprintf("%s(%s) %s",
		f.name, strings.Join(f.params, ", "), f.body)
}

type FunctionCallNode struct {
	callable Node
	args     []Node
}

func (f FunctionCallNode) Type() NodeType { return ndFunctionCall }
func (f FunctionCallNode) String() string {
	args := make([]string, len(f.args), len(f.args))
	for i, arg := range f.args {
		args[i] = arg.String()
	}

	return fmt.Sprintf("%s(%s)", f.callable, strings.Join(args, ", "))
}

type GotoNode struct{ label Node }

func (g GotoNode) Type() NodeType { return ndGoto }
func (g GotoNode) String() string { return fmt.Sprintf("goto %v;", g.label) }

type IdentNode struct {
	value string
}

func (i IdentNode) Type() NodeType { return ndIdent }
func (i IdentNode) String() string { return i.value }

type IntegerNode struct {
	value string
}

type IfNode struct {
	cond     Node
	body     Node
	hasElse  bool
	elseBody Node
}

func (i IfNode) Type() NodeType { return ndIf }
func (i IfNode) String() string {
	var elseStr string = ""

	if i.hasElse {
		elseStr = fmt.Sprintf(" else %v", i.elseBody)
	}

	return fmt.Sprintf("if(%v) %v%s", i.cond, i.body, elseStr)
}

func (i IntegerNode) Type() NodeType { return ndInteger }
func (i IntegerNode) String() string { return i.value }

type LabelNode struct{ name string }

func (l LabelNode) Type() NodeType { return ndLabel }
func (l LabelNode) String() string { return fmt.Sprintf("%s:", l.name) }

type NullNode struct{}

func (n NullNode) Type() NodeType { return ndNull }
func (n NullNode) String() string { return "" }

type ParenNode struct{ node Node }

func (p ParenNode) Type() NodeType { return ndParen }
func (p ParenNode) String() string { return "(" + p.node.String() + ")" }

type ReturnNode struct{ node Node }

func (r ReturnNode) Type() NodeType { return ndReturn }
func (r ReturnNode) String() string { return fmt.Sprintf("return %v;", r.node) }

type StatementNode struct {
	expr Node
}

func (s StatementNode) Type() NodeType { return ndStatement }
func (s StatementNode) String() string { return fmt.Sprintf("%v;", s.expr) }

type StringNode struct {
	value string
}

func (s StringNode) Type() NodeType { return ndString }
func (s StringNode) String() string { return fmt.Sprintf("\"%s\"", s.value) }

type caseNode struct {
	cond       Node
	statements []Node
}

func (c caseNode) String() string {
	str := fmt.Sprintf("case %v:", c.cond)

	for _, stmt := range c.statements {
		str += fmt.Sprintf("\n\t%v", stmt)
	}

	return str
}

type SwitchNode struct {
	cond        Node
	defaultCase []Node
	cases       []caseNode
}

func (s SwitchNode) Type() NodeType { return ndSwitch }
func (s SwitchNode) String() string {
	str := fmt.Sprintf("switch(%v) {", s.cond)

	for _, cs := range s.cases {
		str += "\n" + cs.String()
	}

	if s.defaultCase != nil {
		str += "\ndefault:"
		for _, stmt := range s.defaultCase {
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

func (t TernaryNode) Type() NodeType { return ndTernary }
func (t TernaryNode) String() string {
	return fmt.Sprintf("(%v ? %v : %v)", t.cond, t.trueBody, t.falseBody)
}

type UnaryNode struct {
	oper    string
	node    Node
	postfix bool
}

func (u UnaryNode) Type() NodeType { return ndUnary }
func (u UnaryNode) String() string {
	if u.postfix {
		return fmt.Sprintf("%v%s", u.node, u.oper)
	}
	return fmt.Sprintf("%s%v", u.oper, u.node)
}

type VarDecl struct {
	name    string
	vecDecl bool
	size    string
}

type VarDeclNode struct {
	vars []VarDecl
}

func (v VarDeclNode) Type() NodeType { return ndVarDecl }
func (v VarDeclNode) String() string {
	decls := make([]string, 0, len(v.vars))

	for _, decl := range v.vars {
		var str string

		if decl.vecDecl {
			str = fmt.Sprintf("%s[%s]", decl.name, decl.size)
		} else {
			str = decl.name
		}

		decls = append(decls, str)
	}

	return fmt.Sprintf("auto %s;", strings.Join(decls, ", "))
}

type WhileNode struct {
	cond Node
	body Node
}

func (w WhileNode) Type() NodeType { return ndWhile }
func (w WhileNode) String() string {
	return fmt.Sprintf("while(%v) %v", w.cond, w.body)
}
