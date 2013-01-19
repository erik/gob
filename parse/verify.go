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
	return fmt.Sprintf("Semantic error on %v: %v", s.node, s.msg)
}

func NewSemanticError(node Node, msg string) error {
	return &SemanticError{node, msg}
}

type TranslationUnit struct {
	file  string
	funcs []FunctionNode
	vars  []ExternVarInitNode
}

func (t TranslationUnit) Verify() error {
	return errors.New("Verification not implemented")
}
