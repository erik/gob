package parse

import (
	"errors"
	"fmt"
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

func (t TranslationUnit) Verify() error {

	for _, fn := range t.funcs {
		if err := t.expectStatement(fn.body); err != nil {
			return err
		}

		// TODO: ...
		if err := t.VerifyAssignments(fn); err != nil {
			return err
		}

		if err := t.ResolveDuplicates(); err != nil {
			return err
		}

		if err := t.ResolveExternalDeclarations(fn); err != nil {
			return err
		}
	}

	return errors.New("Verification not fully implemented")
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

	return NewSemanticError(node, "expected statement")
}

func (t TranslationUnit) VerifyFunction(fn FunctionNode) error {
	return nil
}

// Verify that all assignments have a proper LHS and RHS
func (t TranslationUnit) VerifyAssignments(fn FunctionNode) error {
	if bin, ok := fn.body.(BinaryNode); ok {
		if bin.oper == "=" {
			if err := t.expectLHS(bin.left); err != nil {
				return err
			}
			if err := t.expectRHS(bin.right); err != nil {
				return err
			}
		}

	} else if stmt, ok := fn.body.(StatementNode); ok {
		_ = stmt
	}

	return nil
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

func (t TranslationUnit) ResolveExternalDeclarations(fn FunctionNode) error {
	// TODO: This should check that all declared extern vars resolve
	return errors.New("Not yet implemented")
}
