package logger

import (
	"go/ast"
	"x/dsl"
)

func log(n ast.Node, spec dsl.NodeSpec) (ast.Node, []string, error) {
	args := spec.Args
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "fmt",
			},
			Sel: &ast.Ident{
				Name: "Println",
			},
		},
		Args: args,
	}, []string{"fmt"}, nil
}
