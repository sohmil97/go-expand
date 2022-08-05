package logger

import (
	"go/ast"
	"go/token"
	"x/dsl"
)

func log(n ast.Node, spec dsl.NodeSpec) ([]ast.Node, []string, error) {
	args := spec.Args
	// return &ast.CallExpr{
	// 	Fun: &ast.SelectorExpr{
	// 		X: &ast.Ident{
	// 			Name: "fmt",
	// 		},
	// 		Sel: &ast.Ident{
	// 			Name: "Println",
	// 		},
	// 	},
	// 	Args: args,
	// }, []string{"fmt"}, nil

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
				Args: args,
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
