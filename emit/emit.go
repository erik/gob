package emit

import (
	"github.com/erik/gob/parse"
)

type Emitter interface {
	Emit(parse.TranslationUnit) string
}
