package main

type NodeType int

const (
	ndError NodeType = iota
)

type Node interface {
	Type() NodeType
	String() string
}
