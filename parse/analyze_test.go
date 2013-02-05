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
