package emit

import (
	"bufio"
	"fmt"
	"github.com/erik/gob/parse"
	"io"
	"strings"
)

type CEmitter struct {
	writer *bufio.Writer
	indent int
}

func (c CEmitter) Emit(writer io.Writer, unit parse.TranslationUnit) error {
	c.writer = bufio.NewWriter(writer)
	c.indent = 0

	c.EmitHeaders(unit)

	c.EmitLine("\n/* Global variables */")

	for _, v := range unit.Vars {
		c.EmitGlobal(v)
	}

	c.EmitLine("\n/* Function prototypes */")

	for _, f := range unit.Funcs {
		c.EmitFunctionProto(f)
	}

	c.EmitLine("\n/* Function definitions */")

	for _, f := range unit.Funcs {
		c.EmitFunction(f)
	}

	c.writer.Flush()

	return nil
}

func (c *CEmitter) EmitHeaders(unit parse.TranslationUnit) {
	c.EmitLine(fmt.Sprintf(`
/* Generated by gob v%s on TODO:DATE
 *
 * TODO: more interesting info.
 */`, "SOME VERSION IDK"))

	c.EmitLine("")

	c.EmitLine(`#include "bstdlib.h"`)

	c.EmitLine("")
}

func (c *CEmitter) EmitGlobal(v parse.Node) {
	switch v.(type) {
	case parse.ExternVarInitNode:
		var_ := v.(parse.ExternVarInitNode)
		c.EmitLine(fmt.Sprintf("static B_AUTO %v = %v", var_.Name,
			var_.Value))

	case parse.ExternVecInitNode:
		vec := v.(parse.ExternVecInitNode)
		c.EmitPartial(fmt.Sprintf("static B_AUTO %v[][%d] = ",
			vec.Name, vec.Size))

		c.StartBlock()

		for i, val := range vec.Values {
			if i != len(vec.Values)-1 {
				c.EmitLine(fmt.Sprintf("%v,", val))
			} else {
				c.EmitLine(val.String())
			}
		}

		c.EndBlock()

	}
}

func (c *CEmitter) EmitFunctionProto(fn parse.FunctionNode) {
	c.EmitPartial(fmt.Sprintf("static B_AUTO %s(", sanitizeIdentifier(fn.Name)))

	for i, param := range fn.Params {
		c.EmitRaw(fmt.Sprintf("B_AUTO %s", param))

		if i != len(fn.Params)-1 {
			c.EmitRaw(", ")
		}
	}
	c.EmitRaw(");\n")
}

func (c *CEmitter) EmitFunction(fn parse.FunctionNode) {
	c.EmitPartial(fmt.Sprintf("static B_AUTO %s(", sanitizeIdentifier(fn.Name)))

	for i, param := range fn.Params {
		c.EmitRaw(fmt.Sprintf("B_AUTO %s", param))

		if i != len(fn.Params)-1 {
			c.EmitRaw(", ")
		}
	}
	c.EmitRaw(") ")

	c.EmitBlock(fn.Body.(parse.BlockNode))
}

func (c *CEmitter) EmitBlock(block parse.BlockNode) {
	c.StartBlock()

	for _, node := range block.Nodes {
		c.EmitStatement(node)
	}

	c.EndBlock()
}

func (c *CEmitter) EmitStatement(node parse.Node) {
	switch node.(type) {
	case parse.BlockNode:
		c.EmitBlock(node.(parse.BlockNode))
	case parse.BreakNode:
		c.EmitLine("break;")
	case parse.ExternVarDeclNode:
		c.EmitLine(fmt.Sprintf("/* %v */", node))
	case parse.GotoNode:
		c.EmitLine(fmt.Sprintf("goto %s;", node.(parse.GotoNode).Label))
	case parse.IfNode:
		if_ := node.(parse.IfNode)

		c.EmitPartial("if (")
		c.EmitExpression(if_.Cond)

		if _, ok := if_.Body.(parse.BlockNode); ok {
			c.EmitRaw(")\n")
			c.EmitStatement(if_.Body)
		} else {
			c.EmitRaw(")\n")
			c.Indent()
			c.EmitStatement(if_.Body)
			c.Deindent()
		}

		if if_.HasElse {
			c.EmitLine("else")

			if _, ok := if_.ElseBody.(parse.BlockNode); ok {
				c.EmitStatement(if_.ElseBody)
			} else {
				c.Indent()
				c.EmitStatement(if_.ElseBody)
				c.Deindent()
			}
		}
	case parse.LabelNode:
		c.Deindent()
		c.EmitLine(fmt.Sprintf("%s:", node.(parse.LabelNode).Name))
		c.Indent()
	case parse.NullNode:
		c.EmitLine(";")
	case parse.ReturnNode:
		c.EmitPartial("return ")
		c.EmitExpression(node.(parse.ReturnNode).Node)
		c.EmitRaw(";\n")
	case parse.StatementNode:
		c.EmitPartial("")
		c.EmitExpression(node.(parse.StatementNode).Expr)
		c.EmitRaw(";\n")
	case parse.SwitchNode:
		switch_ := node.(parse.SwitchNode)

		c.EmitPartial("switch (")
		c.EmitExpression(switch_.Cond)
		c.EmitRaw(") {\n")

		c.Indent()

		for _, case_ := range switch_.Cases {
			c.Deindent()

			c.EmitLine(fmt.Sprintf("case %v:", case_.Cond))
			c.Indent()
			for _, stmt := range case_.Statements {
				c.EmitStatement(stmt)
			}
		}

		if switch_.DefaultCase != nil {
			for _, stmt := range switch_.DefaultCase {
				c.EmitStatement(stmt)
			}
		}

		c.EndBlock()
	case parse.VarDeclNode:

		c.EmitPartial("B_AUTO ")

		for i, decl := range node.(parse.VarDeclNode).Vars {
			c.EmitRaw(fmt.Sprintf("%s", decl.Name))

			if decl.VecDecl {
				c.EmitRaw(fmt.Sprintf("[%d]", decl.Size))
			}

			if i != len(node.(parse.VarDeclNode).Vars)-1 {
				c.EmitRaw(", ")
			}
		}

		c.EmitRaw(";\n")

	case parse.WhileNode:
		while := node.(parse.WhileNode)

		c.EmitPartial("while (")
		c.EmitExpression(while.Cond)
		c.EmitRaw(")\n")

		if _, ok := while.Body.(parse.BlockNode); ok {
			c.EmitStatement(while.Body)
		} else {
			c.Indent()
			c.EmitStatement(while.Body)
			c.Deindent()
		}
	default:
		fmt.Println(node)
		panic("what are you doing.")
	}
}

func (c *CEmitter) EmitExpression(expr parse.Node) {

	// TODO: Put a bit more care into this, there are almost certainly
	//       incompatibilities.
	//
	// TODO: Need to sanitize anything that could touch an IdentNode
	switch expr.(type) {
	case parse.ArrayAccessNode:
		arr := expr.(parse.ArrayAccessNode)
		c.EmitExpression(arr.Array)
		c.EmitRaw("[")
		c.EmitExpression(arr.Index)
		c.EmitRaw("]")

	case parse.BinaryNode:
		bin := expr.(parse.BinaryNode)
		c.EmitExpression(bin.Left)
		c.EmitRaw(" " + bin.Oper + " ")
		c.EmitExpression(bin.Right)

	case parse.IntegerNode:
		c.EmitRaw(expr.String())

	case parse.FunctionCallNode:
		fun := expr.(parse.FunctionCallNode)
		c.EmitExpression(fun.Callable)
		c.EmitRaw("(")
		for i, arg := range fun.Args {
			c.EmitExpression(arg)

			if i != len(fun.Args)-1 {
				c.EmitRaw(", ")
			}
		}

		c.EmitRaw(")")

	case parse.ParenNode:
		c.EmitRaw("(")
		c.EmitExpression(expr.(parse.ParenNode).Node)
		c.EmitRaw(")")

	case parse.TernaryNode:
		ter := expr.(parse.TernaryNode)
		c.EmitExpression(ter.Cond)
		c.EmitRaw(" ? ")
		c.EmitExpression(ter.TrueBody)
		c.EmitRaw(" : ")
		c.EmitExpression(ter.FalseBody)

	case parse.UnaryNode:
		un := expr.(parse.UnaryNode)
		if un.Postfix {
			c.EmitExpression(un.Node)
			c.EmitRaw(un.Oper)
		} else {
			c.EmitRaw(un.Oper)
			c.EmitExpression(un.Node)
		}

	case parse.IdentNode:
		c.EmitRaw(sanitizeIdentifier(expr.String()))

	case parse.CharacterNode, parse.StringNode:
		c.EmitRaw(escapeString(expr.String()))

	default:
		fmt.Println(expr)
		panic("come on now")
	}
}

func (c *CEmitter) EmitRaw(text string) {
	c.writer.WriteString(text)
}

func (c *CEmitter) EmitPartial(line string) {
	c.writer.WriteString(strings.Repeat("\t", c.indent))
	c.writer.WriteString(line)
}

func (c *CEmitter) EmitLine(line string) {
	c.writer.WriteString(strings.Repeat("\t", c.indent))
	c.writer.WriteString(line)
	c.writer.WriteString("\n")
}

func (c *CEmitter) StartBlock() {
	c.EmitLine("{")
	c.Indent()
}

func (c *CEmitter) EndBlock() {
	c.Deindent()
	c.EmitLine("}")
}

func (c *CEmitter) Indent() {
	c.indent += 1
}

func (c *CEmitter) Deindent() {
	c.indent -= 1

	if c.indent < 0 {
		c.indent = 0
	}
}

// Return a C version of the given B identifier
func sanitizeIdentifier(ident string) string {
	return strings.Replace(ident, ".", "_", -1)
}

// *0	null
// *e	end-of-file
// *(	{
// *)	}
// *t	tab
// **	*
// *'	'
// *"	"
// *n	new line
func escapeString(str string) string {
	escaped := ""

	for i := 0; i < len(str); i++ {
		if str[i] == '*' {
			switch str[i+1] {
			case '0':
				escaped += "\\0"
			case 'e':
				// EOT
				escaped += "\\0"
			case '(':
				escaped += "{"
			case ')':
				escaped += "}"
			case 't':
				escaped += "\\t"
			case '*':
				escaped += "*"
			case '\'':
				escaped += "'"
			case '"':
				escaped += "\""
			case 'n':
				escaped += "\\n"
			}

			i += 1
		} else {
			escaped += string(str[i])
		}
	}

	return escaped
}
