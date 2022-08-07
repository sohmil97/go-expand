package generator

import (
	"fmt"
	"go/ast"
	"go/token"
	"x/internal/types"
)

func GenerateFuncAST(selector *ast.SelectorExpr, args []ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:  selector,
		Args: args,
	}
}

func GetFuncSelector(pkg string, name string) *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   GenerateIdentifierAST(pkg, SELECTOR_IDENTIFIER),
		Sel: GenerateIdentifierAST(name, SELECTOR_IDENTIFIER),
	}
}

func GetFuncArgs(args types.ArgSpec) []ast.Expr {
	exprs := make([]ast.Expr, 0)
	for k, v := range args {
		if k.IsVarArgs {
			for _, elm := range v.([]interface{}) {
				value := fmt.Sprint(elm)
				exprs = append(exprs, generateFuncStringArgAST(value))
			}
		} else {
			value := fmt.Sprint(v)
			exprs = append(exprs, generateFuncStringArgAST(value))
		}
	}
	return exprs
}

func generateFuncFloatArgAST(value string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.FLOAT,
		Value: value,
	}
}

func generateFuncIntArgAST(value string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: value,
	}
}

func generateFuncStringArgAST(value string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: value,
	}
}
