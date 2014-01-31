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
			p := parse.NewParser(test, file)
			if _, err := p.Parse(); err != nil {
				t.Errorf("%s failed: %v", test, err)
			}
		}
	}
}
