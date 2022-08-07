package types

import "go/ast"

// Processor changes given ast into another one in place in order to replace the marker usages with actual code.
type Processor func(spec NodeSpec) ([]ast.Node, []string, error)
