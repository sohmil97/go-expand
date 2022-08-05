package generator

import (
	"go/ast"
	"go/token"
)

// TODO create generator utils to provide to dsl constructors
func createAssignment(identifier string, value string, kind token.Token) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{
				Name: identifier,
			},
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.BasicLit{
				Kind:  kind,
				Value: value,
			},
		},
	}
}
