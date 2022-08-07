package generator

import "go/ast"

const (
	CONST_IDENTIFIER = iota
	VARIABLE_IDENTIFIER
	SELECTOR_IDENTIFIER
)

func GenerateIdentifierAST(name string, tp int) *ast.Ident {
	switch tp {
	case CONST_IDENTIFIER:
		return createConstIdentifierAST(name)
	case VARIABLE_IDENTIFIER:
		return createVarIdentifierAST(name)
	default:
		return createSelectorIdentifierAST(name)
	}
}

func createConstIdentifierAST(name string) *ast.Ident {
	return &ast.Ident{
		Name: name,
		Obj: &ast.Object{
			Kind: ast.Con,
			Name: name,
		},
	}
}

func createVarIdentifierAST(name string) *ast.Ident {
	return &ast.Ident{
		Name: name,
		Obj: &ast.Object{
			Kind: ast.Var,
			Name: name,
		},
	}
}

func createSelectorIdentifierAST(name string) *ast.Ident {
	return &ast.Ident{
		Name: name,
	}
}
