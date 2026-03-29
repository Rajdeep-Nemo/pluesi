package ast

// Node is the base interface for all nodes in the tree
type Node interface {
	TokenString() string
	String() string
}

// Interface for statements
type Statement interface {
	Node
	statementNode()
}

// Interface for expressions
type Expression interface {
	Node
	expressionNode()
}

// Root node
type Program struct {
	Statements []Statement
}
