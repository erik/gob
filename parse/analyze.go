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
	File  string
	Funcs []FunctionNode
	Vars  []Node
}

func (t TranslationUnit) String() string {
	str := fmt.Sprintf("%s:", t.File)

	for _, v := range t.Vars {
		str += fmt.Sprintf("%v\n", v)
	}

	str += "\n\n"

	for _, f := range t.Funcs {
		str += fmt.Sprintf("%v\n", f)
	}

	return str
}

func (t TranslationUnit) Verify() error {

	if err := t.ResolveDuplicates(); err != nil {
		return err
	}

	for _, fn := range t.Funcs {

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
		if node.(UnaryNode).Oper == "*" {
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
		for _, n := range node.(BlockNode).Nodes {
			if err := t.visitExpressions(n, visit); err != nil {
				return err
			}
		}
	case FunctionNode:
		if err := t.visitExpressions(node.(FunctionNode).Body, visit); err != nil {
			return err
		}

	case IfNode:
		if err := visit(node.(IfNode).Cond); err != nil {
			return err
		}

		if err := t.visitExpressions(node.(IfNode).Body, visit); err != nil {
			return err
		}

		if node.(IfNode).HasElse {
			if err := t.visitExpressions(node.(IfNode).ElseBody, visit); err != nil {
				return err
			}
		}

	case SwitchNode:
		if err := visit(node.(SwitchNode).Cond); err != nil {
			return err
		}

		for _, stmt := range node.(SwitchNode).DefaultCase {
			if err := t.visitExpressions(stmt, visit); err != nil {
				return err
			}
		}

		for _, case_ := range node.(SwitchNode).Cases {
			if err := visit(case_.Cond); err != nil {
				return err
			}

			if err := t.visitExpressions(case_, visit); err != nil {
				return err
			}
		}

	case WhileNode:
		if err := visit(node.(WhileNode).Cond); err != nil {
			return err
		}

		if err := t.visitExpressions(node.(WhileNode).Body, visit); err != nil {
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
		for _, n := range node.(BlockNode).Nodes {
			if err := t.expectStatement(n); err != nil {
				return err
			}

			if err := t.visitStatements(n, visit); err != nil {
				return err
			}
		}
	case FunctionNode:
		if err := t.expectStatement(node.(FunctionNode).Body); err != nil {
			return err
		}

		if err := t.visitStatements(node.(FunctionNode).Body, visit); err != nil {
			return err
		}

	case GotoNode:
		if err := visit(node); err != nil {
			return err
		}

	case IfNode:
		if err := t.visitStatements(node.(IfNode).Body, visit); err != nil {
			return err
		}

		if node.(IfNode).HasElse {
			if err := t.visitStatements(node.(IfNode).ElseBody, visit); err != nil {
				return err
			}
		}

	case BreakNode, ExternVarDeclNode, ExternVarInitNode,
		ExternVecInitNode, LabelNode, ReturnNode, StatementNode, VarDeclNode:
		if err := visit(node); err != nil {
			return err
		}
	case SwitchNode:

		for _, stmt := range node.(SwitchNode).DefaultCase {
			if err := t.visitStatements(stmt, visit); err != nil {
				return err
			}
		}

		for _, case_ := range node.(SwitchNode).Cases {
			if err := t.visitStatements(case_, visit); err != nil {
				return err
			}
		}

	case WhileNode:
		if err := t.visitStatements(node.(WhileNode).Body, visit); err != nil {
			return err
		}
	}

	return nil
}

func (t TranslationUnit) VerifyFunction(fn FunctionNode) error {

	if err := t.expectNodeType(fn.Body, reflect.TypeOf(BlockNode{})); err != nil {
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

	if err := t.visitStatements(fn.Body, visiter); err != nil {
		return err
	}

	return nil
}

// Verify that all assignments have a proper LHS and RHS
func (t TranslationUnit) VerifyAssignments(fn FunctionNode) error {
	visit := func(node Node) error {
		stmt, ok := node.(StatementNode)

		if !ok {
			return nil
		}
		if bin, ok := stmt.Expr.(BinaryNode); ok {
			if bin.Oper == "=" {
				if err := t.expectLHS(bin.Left); err != nil {
					return err
				}
				if err := t.expectRHS(bin.Right); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return t.visitStatements(fn.Body, visit)
}

// TODO: resolve auto variable declarations within function definitions
func (t TranslationUnit) ResolveDuplicates() error {
	idents := map[string]Node{}

	for _, fn := range t.Funcs {
		if _, ok := idents[fn.Name]; ok {
			return NewSemanticError(fn, "Duplicate function name")
		}

		idents[fn.Name] = fn
	}

	for _, v := range t.Vars {
		var name string

		switch v.(type) {
		case ExternVecInitNode:
			name = v.(ExternVecInitNode).Name
		case ExternVarInitNode:
			name = v.(ExternVarInitNode).Name
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
			if _, ok := labels[node.(LabelNode).Name]; ok {
				return NewSemanticError(node, "duplicate label definition")
			}
			labels[node.(LabelNode).Name] = true
		case GotoNode:
			gotos = append(gotos, node.(GotoNode))
		}
		return nil
	}

	if err := t.visitStatements(fn, visiter); err != nil {
		return err
	}

	for _, node := range gotos {
		if _, ok := labels[node.Label]; !ok {
			return NewSemanticError(node, "unresolved goto")
		}
	}

	return nil
}
