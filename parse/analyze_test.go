package parse

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestAnalyzeDuplicate(t *testing.T) {
	unit, err := NewParser("", strings.NewReader("a; b; c;")).Parse()

	if err != nil {
		t.Errorf("Parse failed: %v", err)
	} else if err = unit.ResolveDuplicates(); err != nil {
		t.Errorf("Resolve duplicates: %v", err)
	}

	unit, err = NewParser("", strings.NewReader("a; b; a;")).Parse()
	if err != nil {
		t.Errorf("Parse failed: %v", err)
	} else if err = unit.ResolveDuplicates(); err == nil {
		t.Errorf("Allowed duplicate variable/variable")
	}

	unit, err = NewParser("", strings.NewReader("a(){} a(){}")).Parse()
	if err != nil {
		t.Errorf("Parse failed: %v", err)
	} else if err = unit.ResolveDuplicates(); err == nil {
		t.Errorf("Allowed duplicate func/func")
	}

	unit, err = NewParser("", strings.NewReader("a; a(){}")).Parse()
	if err != nil {
		t.Errorf("Parse failed: %v", err)
	} else if err = unit.ResolveDuplicates(); err == nil {
		t.Errorf("Allowed duplicate func/variable")
	}
}

func TestLHS(t *testing.T) {
	var unit TranslationUnit

	// Simple lhs cases
	assert.Nil(t, unit.expectLHS(IdentNode{"foo"}), "ident node LHS")
	assert.Nil(t, unit.expectLHS(ArrayAccessNode{IdentNode{"abc"},
		IntegerNode{"2"}}), "array access lhs")
	assert.Nil(t, unit.expectLHS(UnaryNode{"*", IntegerNode{"1"}, false}),
		"unary node lhs")
}

func TestRHS(t *testing.T) {
	// TODO: write me
}

func TestVerifyAssignments(t *testing.T) {
	//var unit TranslationUnit

}
