package test

import (
	"github.com/boredomist/gob/parse"
	"os"
	"testing"
)

var tests = []string{"convert.b", "copy.b", "lower.b", "snide.b"}

func TestExampleDummy(t *testing.T) {
}

// TODO: parser is not quite ready for this, and I'd rather not have every test fail
func TestExamples(t *testing.T) {
	for _, test := range tests {
		if file, err := os.Open(test); err != nil {
			t.Errorf("failed to open test: %s", err)
		} else {
			p := parse.NewParser(test, file)
			if _, err := p.Parse(); err != nil {
				t.Errorf("%s failed: %v", test, err)
			}
		}
	}
}
