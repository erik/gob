package main

type TokenType int

const (
	tkEof  TokenType = iota
	tkMain           = iota
)

type Token struct {
	tkType TokenType
	string string
}

type Lexer struct {
}
