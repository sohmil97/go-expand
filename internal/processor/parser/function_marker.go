package parser

import (
	"fmt"
	"go/ast"
	"go/types"
)

type FunctionMarker struct {
	Signature signature
	Args      map[paramSpec]interface{}
}

type identifier struct {
	Name string
	Type types.Type
}

type paramSpec struct {
	identifier
	IsVarArgs bool
}

type argSpec map[paramSpec]interface{}

type returnSpec struct {
	identifier

	IsErr     bool
	IsCleanup bool
}

type signature struct {
	Name    string
	Pkg     *types.Package
	Params  []paramSpec
	Returns []returnSpec
}

func (p *parser) processFuncMarker(fn *ast.CallExpr, fnObj *types.Func) *FunctionMarker {
	marker := &FunctionMarker{
		Args: map[paramSpec]interface{}{},
	}

	sig := fnObj.Type().(*types.Signature)
	fsig := signature{
		Name:   fnObj.Name(),
		Pkg:    fnObj.Pkg(),
		Params: make([]paramSpec, 0),
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

func (p *parser) extractFuncArgs(args []ast.Expr, params *types.Tuple, hasVarArgs bool) argSpec {
	processedArgs := make(argSpec)
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		isVarArgsParam := hasVarArgs && i == params.Len()-1
		id := paramSpec{
			identifier: identifier{
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

func (p parser) extractFuncReturns(returns *types.Tuple) []returnSpec {
	processedReturns := make([]returnSpec, returns.Len())
	for i := range processedReturns {
		processedReturns[i].identifier = identifier{
			Name: returns.At(0).Name(),
			Type: returns.At(0).Type(),
		}
		processedReturns[i].IsErr = types.Identical(returns.At(i).Type(), types.Universe.Lookup("error").Type())
		processedReturns[i].IsCleanup = types.Identical(returns.At(i).Type(), types.NewSignatureType(nil, nil, nil, nil, nil, false))
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
