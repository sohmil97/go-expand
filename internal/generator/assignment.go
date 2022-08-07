package generator

import (
	"go/ast"
	"go/token"
)

func GenerateVarAssignmentAST() {

}

func GenerateSugarAssignmentAST(lhs []ast.Expr, rhs []ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: lhs,
		Tok: token.DEFINE,
		Rhs: rhs,
	}
}
