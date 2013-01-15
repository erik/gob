# gob

Gob is an implementation of the B language, written in Go.

Currently the project is in its infancy and exists only as a
collection of unit tests. Run `go test` to check for sanity.

I aim to get a fully functional B-language compiler out of this
project, with compilation to native code through intermediate C, LLVM
IR, or asm generation, though this is currently undecided. C will
probably come first, just due to it being quite similar to B and
likely the easiest code generator to write.

I chose Go for this project for two reasons:

1. Ken Thompson had a role in the creation of both Go and B.
2. I had just finished marathoning Arrested Development, and found the
name to be too good to pass up.

![](http://i.imgur.com/M7nJp.png)
