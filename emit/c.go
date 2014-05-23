package emit

import (
	"github.com/erik/gob/parse"
	"bytes"
)


type CEmitter struct {}

func (c CEmitter) Emit(parse.TranslationUnit) string {
	var buf bytes.Buffer

	buf.WriteString("int main() { return -1; }")

	return buf.String()
}
