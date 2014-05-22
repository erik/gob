package parse

import (
	"fmt"
	"reflect"
)

type SemanticError struct {
	node Node
	msg  string
}

func (s *SemanticError) Error() string {
	return fmt.Sprintf("Semantic error on `%v`: %v", s.node, s.msg)
}

func NewSemanticError(node Node, msg string) error {
	return &SemanticError{node, msg}
}

type TranslationUnit struct {
	file  string
	funcs []FunctionNode
	vars  []Node
}

func (t TranslationUnit) String() string {
	str := fmt.Sprintf("%s:", t.file)

	for _, v := range t.vars {
		str += fmt.Sprintf("%v\n", v)
	}

	str += "\n\n"

	for _, f := range t.funcs {
		str += fmt.Sprintf("%v\n", f)
	}

	return str
}

func (t TranslationUnit) Verify() error {

	if err := t.ResolveDuplicates(); err != nil {
		return err
	}

	for _, fn := range t.funcs {

		if err := t.VerifyFunction(fn); err != nil {
			return err
		}

		if err := t.VerifyAssignments(fn); err != nil {
			return err
		}

		if err := t.ResolveLabels(fn); err != nil {
			return err
		}
	}

	return nil
}

func (t TranslationUnit) expectLHS(node Node) error {
	switch node.(type) {
	case ArrayAccessNode, IdentNode:
		return nil
	case UnaryNode:
		if node.(UnaryNode).oper == "*" {
			return nil
		}
	}

	return NewSemanticError(node, "expected lvalue")
}

func (t TranslationUnit) expectRHS(node Node) error {
	if IsExpr(node) {
		return nil
	}

	return NewSemanticError(node, "expected rvalue")
}

func (t TranslationUnit) expectStatement(node Node) error {
	if IsStatement(node) {
		return nil
	}

	return NewSemanticError(node, "expected statement, got "+reflect.TypeOf(node).Name())
}

func (t TranslationUnit) expectNodeType(node Node, kind reflect.Type) error {
	if reflect.TypeOf(node) != kind {
		return NewSemanticError(node, "expected "+kind.Name())
	}

	return nil
}

func (t TranslationUnit) visitExpressions(node Node, visit func(Node) error) error {
	if IsExpr(node) {
		return visit(node)
	}

	switch node.(type) {
	case BlockNode:
		for _, n := range node.(BlockNode).nodes {
			if err := t.visitExpressions(n, visit); err != nil {
				return err
			}
		}
	case FunctionNode:
		if err := t.visitExpressions(node.(FunctionNode).body, visit); err != nil {
			return err
		}

	case IfNode:
		if err := visit(node.(IfNode).cond); err != nil {
			return err
		}

		if err := t.visitExpressions(node.(IfNode).body, visit); err != nil {
			return err
		}

		if node.(IfNode).hasElse {
			if err := t.visitExpressions(node.(IfNode).elseBody, visit); err != nil {
				return err
			}
		}

	case SwitchNode:
		if err := visit(node.(SwitchNode).cond); err != nil {
			return err
		}

		for _, stmt := range node.(SwitchNode).defaultCase {
			if err := t.visitExpressions(stmt, visit); err != nil {
				return err
			}
		}

		for _, case_ := range node.(SwitchNode).cases {
			if err := visit(case_.cond); err != nil {
				return err
			}

			if err := t.visitExpressions(case_, visit); err != nil {
				return err
			}
		}

	case WhileNode:
		if err := visit(node.(WhileNode).cond); err != nil {
			return err
		}

		if err := t.visitExpressions(node.(WhileNode).body, visit); err != nil {
			return err
		}
	}


	return nil
}

func (t TranslationUnit) visitStatements(node Node, visit func(Node) error) error {

	if err := t.expectStatement(node); err != nil {
		return err
	}

	switch node.(type) {
	case BlockNode:
		for _, n := range node.(BlockNode).nodes {
			if err := t.expectStatement(n); err != nil {
				return err
			}

			if err := t.visitStatements(n, visit); err != nil {
				return err
			}
		}
	case FunctionNode:
		if err := t.expectStatement(node.(FunctionNode).body); err != nil {
			return err
		}

		if err := t.visitStatements(node.(FunctionNode).body, visit); err != nil {
			return err
		}

	case GotoNode:
		if err := visit(node); err != nil {
			return err
		}

	case IfNode:
		if err := t.visitStatements(node.(IfNode).body, visit); err != nil {
			return err
		}

		if node.(IfNode).hasElse {
			if err := t.visitStatements(node.(IfNode).elseBody, visit); err != nil {
				return err
			}
		}

	case BreakNode, ExternVarDeclNode, ExternVarInitNode,
		ExternVecInitNode, LabelNode, ReturnNode, StatementNode, VarDeclNode:
		if err := visit(node); err != nil {
			return err
		}
	case SwitchNode:

		for _, stmt := range node.(SwitchNode).defaultCase {
			if err := t.visitStatements(stmt, visit); err != nil {
				return err
			}
		}

		for _, case_ := range node.(SwitchNode).cases {
			if err := t.visitStatements(case_, visit); err != nil {
				return err
			}
		}

	case WhileNode:
		if err := t.visitStatements(node.(WhileNode).body, visit); err != nil {
			return err
		}
	}

	return nil
}

func (t TranslationUnit) VerifyFunction(fn FunctionNode) error {

	if err := t.expectNodeType(fn.body, reflect.TypeOf(BlockNode{})); err != nil {
		return err
	}

	// Ensure variables are declared at the beginning of functions
	endDecls := false

	visiter := func(stmt Node) error {
		switch stmt.(type) {
		case ExternVarDeclNode, VarDeclNode:
			if endDecls {
				return NewSemanticError(stmt, "var declaration in middle of block")
			}
		default:
			endDecls = true
		}
		return nil
	}

	if err := t.visitStatements(fn.body, visiter); err != nil {
		return err
	}

	return nil
}

// Verify that all assignments have a proper LHS and RHS
func (t TranslationUnit) VerifyAssignments(fn FunctionNode) error {
	visit := func(node Node) error {
		stmt, ok := node.(StatementNode)

		if !ok { return nil }
		if bin, ok := stmt.expr.(BinaryNode); ok {
			if bin.oper == "=" {
				if err := t.expectLHS(bin.left); err != nil {
					return err
				}
				if err := t.expectRHS(bin.right); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return t.visitStatements(fn.body, visit)
}

// TODO: resolve auto variable declarations within function definitions
func (t TranslationUnit) ResolveDuplicates() error {
	idents := map[string]Node{}

	for _, fn := range t.funcs {
		if _, ok := idents[fn.name]; ok {
			return NewSemanticError(fn, "Duplicate function name")
		}

		idents[fn.name] = fn
	}

	for _, v := range t.vars {
		var name string

		switch v.(type) {
		case ExternVecInitNode:
			name = v.(ExternVecInitNode).name
		case ExternVarInitNode:
			name = v.(ExternVarInitNode).name
		default:
			return NewSemanticError(v, "Not variable init")
		}

		if _, ok := idents[name]; ok {
			return NewSemanticError(v, "Duplicate variable name")
		}

		idents[name] = v
	}

	return nil
}

// Make sure all goto jump to valid places
func (t TranslationUnit) ResolveLabels(fn FunctionNode) error {
	labels := map[string]bool{}
	gotos := []GotoNode{}

	visiter := func(node Node) error {
		switch node.(type) {
		case LabelNode:
			if _, ok := labels[node.(LabelNode).name]; ok {
				return NewSemanticError(node, "duplicate label definition")
			}
			labels[node.(LabelNode).name] = true
		case GotoNode:
			gotos = append(gotos, node.(GotoNode))
		}
		return nil
	}

	if err := t.visitStatements(fn, visiter); err != nil {
		return err
	}

	for _, node := range gotos {
		if _, ok := labels[node.label.(IdentNode).value]; !ok {
			return NewSemanticError(node, "unresolved goto")
		}
	}

	return nil
}
