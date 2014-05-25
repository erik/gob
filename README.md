# gob
[![Build Status](https://travis-ci.org/erik/gob.png?branch=master)](https://travis-ci.org/erik/gob)

Gob is an implementation of the B language, written in Go.

Currently the project is in its infancy and is probably the most brittle
compiler ever written. Run `go test ./...` to check for sanity.

`go build .` (or `go get github.com/erik/gob`) will give you an executable that
parses B files given to it on the command line and generates C output. You
can't compile this quite yet, because the B standard library wrapper hasn't
been written yet.

`$ gob examples/snide.b`

I aim to get a fully functional B-language compiler out of this
project, with compilation to native code through intermediate C, LLVM
IR, or asm generation, though this is currently undecided. C will
probably come first, just due to it being quite similar to B and
likely the easiest code generator to write.

I chose Go for this project for two reasons:

1. Ken Thompson had a role in the creation of both Go and B.
2. I had just finished marathoning Arrested Development.

## Status

* Front end (lex, parse, AST generation)
  * Lexer complete
  * Parser complete
  * AST definitions complete
  * Unit tests for parser and lexer are lacking, probably accept some incorrect
    syntax constructions.
* Middle end (semantic analysis, optimizations)
  * Semantic analysis is limited, but working.
  * No optimization performed yet.
* Back end (code generator)
  * C code generator is almost functional, needs some supporting library code
    to be entirely working.

## Links

* http://cm.bell-labs.com/cm/cs/who/dmr/btut.html
* http://cm.bell-labs.com/who/dmr/bref.html

## License

Copyright (c) 2013 Erik Price

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
