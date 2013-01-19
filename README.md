# gob
[![Build Status](https://travis-ci.org/boredomist/gob.png?branch=master)](https://travis-ci.org/boredomist/gob)

Gob is an implementation of the B language, written in Go.

Currently the project is in its infancy and exists only as a
collection of unit tests. Run `go test ./...` to check for sanity.

`go build .` (or `go get github.com/boredomist/gob`) will give you an
executable that parses files passed on the command line, and not much
else.

I aim to get a fully functional B-language compiler out of this
project, with compilation to native code through intermediate C, LLVM
IR, or asm generation, though this is currently undecided. C will
probably come first, just due to it being quite similar to B and
likely the easiest code generator to write.

I chose Go for this project for two reasons:

1. Ken Thompson had a role in the creation of both Go and B.
2. I had just finished marathoning Arrested Development, and found the
name to be too good to pass up.

## Status

* Front end (lex, parse, AST generation)
  * Lexer mostly working (needs some simple modification to parse all
    identifiers)
  * Parser mostly complete, a few syntactic structures not yet handled
* Middle end (semantic analysis, optimizations)
  * Not yet implemented
* Back end (code generator)
  * Not yet implemented

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
