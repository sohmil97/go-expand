package parser

import (
	"fmt"
	"go/ast"
	goTypes "go/types"
	"x/internal/types"
)

type FunctionMarker struct {
	Signature types.Signature
	Args      map[types.ParamSpec]interface{}
}

func (p *parser) processFuncMarker(fn *ast.CallExpr, fnObj *goTypes.Func) *FunctionMarker {
	marker := &FunctionMarker{
		Args: map[types.ParamSpec]interface{}{},
	}

	sig := fnObj.Type().(*goTypes.Signature)
	fsig := types.Signature{
		Name:   fnObj.Name(),
		Pkg:    fnObj.Pkg(),
		Params: make([]types.ParamSpec, 0),
	}

	// Extract function args
	args := p.extractFuncArgs(fn.Args, sig.Params(), sig.Variadic())
	for id, value := range args {
		fsig.Params = append(fsig.Params, id)
		marker.Args[id] = value
	}

	// Extract function return values
	fsig.Returns = p.extractFuncReturns(sig.Results())

	marker.Signature = fsig

	return marker
}

func (p *parser) extractFuncArgs(args []ast.Expr, params *goTypes.Tuple, hasVarArgs bool) types.ArgSpec {
	processedArgs := make(types.ArgSpec)
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		isVarArgsParam := hasVarArgs && i == params.Len()-1
		id := types.ParamSpec{
			Identifier: types.Identifier{
				Name: param.Name(),
				Type: param.Type(),
			},
			IsVarArgs: isVarArgsParam,
		}
		// process arg values
		if isVarArgsParam {
			values := make([]interface{}, 0)
			for _, arg := range args[i:] {
				values = append(values, p.processArg(arg))
			}
			processedArgs[id] = values
		} else {
			processedArgs[id] = p.processArg(args[i])
		}
	}
	return processedArgs
}

func (p parser) extractFuncReturns(returns *goTypes.Tuple) []types.ReturnSpec {
	processedReturns := make([]types.ReturnSpec, returns.Len())
	for i := range processedReturns {
		processedReturns[i].Identifier = types.Identifier{
			Name: returns.At(0).Name(),
			Type: returns.At(0).Type(),
		}
		processedReturns[i].IsErr = goTypes.Identical(returns.At(i).Type(), goTypes.Universe.Lookup("error").Type())
		processedReturns[i].IsCleanup = goTypes.Identical(returns.At(i).Type(), goTypes.NewSignatureType(nil, nil, nil, nil, nil, false))
	}
	return processedReturns
}

func (p *parser) processArg(arg ast.Expr) interface{} {
	switch v := arg.(type) {
	case *ast.BasicLit:
		return v.Value
	case *ast.Ident:
		return v.Obj.Data
	default:
		fmt.Println("Unrecognized type")
		return ""
	}
}
