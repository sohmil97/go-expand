package dsl

import "go/ast"

type NodeSpec struct {
	Args []ast.Expr
}

// Processor changes given ast into another one in place in order to replace the marker usages with actual code.
type Processor func(call ast.Node, spec NodeSpec) ([]ast.Node, []string, error)
