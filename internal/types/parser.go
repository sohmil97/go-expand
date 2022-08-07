package types

import (
	"go/ast"
	"go/types"
)

type NodeSpec struct {
	Node ast.Node
	Args ArgSpec
}

type ImportSpec struct {
	Path string
	Name string
}

type Identifier struct {
	Name string
	Type types.Type
}

type ParamSpec struct {
	Identifier
	IsVarArgs bool
}

type ArgSpec map[ParamSpec]interface{}

type ReturnSpec struct {
	Identifier

	IsErr     bool
	IsCleanup bool
}

type Signature struct {
	Name    string
	Pkg     *types.Package
	Params  []ParamSpec
	Returns []ReturnSpec
}
