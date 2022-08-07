package generator

import (
	"fmt"
	"go/ast"
	"x/internal/types"
)

func GenerateExprStatement(body ast.Expr) *ast.ExprStmt {
	return &ast.ExprStmt{
		X: body,
	}
}

func GetArgsString(args types.ArgSpec) string {
	elements := make([]interface{}, 0)
	for k, v := range args {
		if k.IsVarArgs {
			elements = append(elements, v.([]interface{})...)
		} else {
			elements = append(elements, v)
		}
	}
	return getArgString(elements)
}

func getArgString(elements []interface{}) string {
	if len(elements) == 0 {
		return ""
	}
	return processString(elements)
}

func processString(elements []interface{}) string {
	if len(elements) == 1 {
		return fmt.Sprintf("%v", elements[0])
	}
	return fmt.Sprintf("%v, %s", elements[0], processString(elements[1:]))
}
