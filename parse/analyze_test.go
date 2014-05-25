package parse

import (
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
	if err := unit.expectLHS(IdentNode{"foo"}); err != nil {
		t.Errorf("ident node LHS")
	}
	if err := unit.expectLHS(ArrayAccessNode{IdentNode{"abc"}, IntegerNode{2}}); err != nil {
		t.Errorf("array access lhs")
	}
	if err := unit.expectLHS(UnaryNode{"*", IntegerNode{1}, false}); err != nil {
		t.Errorf("unary node lhs")
	}
}

func TestRHS(t *testing.T) {
	// TODO: write me
}

func TestVerifyAssignments(t *testing.T) {
	unit, err := NewParser("", strings.NewReader(`
good() { a = 1; a = 1 + 2; a = (1 + (a = 2)); a = a; a[0] = 1; a[1+2+a] = a;}
bad() { 1 = a; 'this' = 'that';}`)).Parse()

	if err != nil {
		t.Errorf("Parse failed: %v", err)
	} else if err = unit.VerifyAssignments(unit.Funcs[0]); err != nil {
		t.Errorf("verify good assignments failed: %v", err)
	} else if err = unit.VerifyAssignments(unit.Funcs[1]); err == nil {
		t.Errorf("verify bad assignements passed", err)
	}
}
