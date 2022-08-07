package logger

import (
	"go/ast"
	"x/internal/generator"
	"x/internal/types"
)

func log(spec types.NodeSpec) ([]ast.Node, []string, error) {
	// TODO: consider kotlin dsl style
	return []ast.Node{
		generator.GenerateExprStatement(
			generator.GenerateFuncAST(
				generator.GetFuncSelector("fmt", "Println"),
				generator.GetFuncArgs(spec.Args),
			),
		),
	}, []string{"fmt"}, nil
}
