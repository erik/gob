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
	return errors.New("Verification not implemented")
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

func (t TranslationUnit) ResolveExternalDeclarations() error {
	// TODO: This should check that all declared extern vars resolve
	return errors.New("Not yet implemented")
}
