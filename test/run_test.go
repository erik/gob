package test

import (
	"github.com/erik/gob/parse"
	"os"
	"testing"
)

var tests = []string{"convert.b", "copy.b", "lower.b", "snide.b"}

func TestExamples(t *testing.T) {
	for _, test := range tests {
		if file, err := os.Open(test); err != nil {
			t.Errorf("failed to open test: %s", err)
		} else {
			var unit parse.TranslationUnit
			var err error

			p := parse.NewParser(test, file)
			if unit, err = p.Parse(); err != nil {
				t.Errorf("%s failed to parse: %v", test, err)
			}

			if err = unit.Verify(); err != nil {
				t.Errorf("%s failed to verify: %v\n", test, err)
			}

		}
	}
}
