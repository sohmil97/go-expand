package logger

import (
	"go/ast"
	"go/token"
	"x/internal/types"
)

func log(spec types.NodeSpec) ([]ast.Node, []string, error) {
	const template string = "fmt.print(%s)"
	// fnArgs := ""
	// for k, v := range spec.Args {
	// 	k
	// }

	return []ast.Node{
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "fmt",
					},
					Sel: &ast.Ident{
						Name: "Println",
					},
				},
				// Args: spec.Args,
			},
		},
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "a",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"hello\"",
				},
			},
		},
	}, []string{"fmt"}, nil
}
