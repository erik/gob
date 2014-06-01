package emit

import (
	"github.com/erik/gob/parse"
	"io"
)

type Emitter interface {
	Emit(io.Writer, parse.TranslationUnit) error
}
