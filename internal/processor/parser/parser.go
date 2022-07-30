package parser

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"strings"
	"x/internal/dsl"

	"golang.org/x/tools/go/packages"
)

const IMPORT_PATH = "x/dsl"

const (
	EXPR_STMT = iota
)

type Marker struct {
	Sig funcSignature

	Args []ast.Expr
	Node ast.Node
	Type int
}

type Parser interface {
	Load(ctx context.Context, directory string) ([]*packages.Package, []error)
	ParseFile(pkg *packages.Package, f *ast.File) ([]*Marker, []error)

	ExtractArgs(n ast.Node) []ast.Expr
}

type parser struct {
	env        []string
	patterns   []string
	buildFlags []string

	mode packages.LoadMode
}

type ParserOpt func(p *parser)

func WithLoadMode(mode packages.LoadMode) ParserOpt {
	return func(p *parser) {
		p.mode = mode
	}
}

func WithPatterns(patterns []string) ParserOpt {
	return func(p *parser) {
		p.patterns = patterns
	}
}

func WithBuildFlags(flags []string) ParserOpt {
	return func(p *parser) {
		p.buildFlags = flags
	}
}

func WithEnv(env []string) ParserOpt {
	return func(p *parser) {
		p.env = env
	}
}

var defaultOpts = []ParserOpt{
	WithLoadMode(packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedDeps |
		packages.NeedImports |
		packages.NeedTypes |
		packages.NeedTypesSizes |
		packages.NeedSyntax |
		packages.NeedTypesInfo, // LoadAllSyntax
	),
	WithPatterns([]string{"."}),
	WithBuildFlags([]string{"-tags=x-preprocessor"}),
}

func New(opts ...ParserOpt) Parser {
	parser := new(parser)
	for _, opt := range defaultOpts {
		opt(parser)
	}
	for _, opt := range opts {
		opt(parser)
	}
	return parser
}

func (p *parser) Load(ctx context.Context, directory string) ([]*packages.Package, []error) {
	cfg := &packages.Config{
		Context:    ctx,
		Mode:       p.mode,
		Dir:        directory,
		Env:        p.env,
		BuildFlags: p.buildFlags,
	}
	pkgs, err := packages.Load(cfg, p.patterns...)
	if err != nil {
		return nil, []error{err}
	}
	var errs []error
	for _, p := range pkgs {
		for _, e := range p.Errors {
			errs = append(errs, e)
		}
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return pkgs, nil
}

func (p *parser) ParseFile(pkg *packages.Package, f *ast.File) ([]*Marker, []error) {
	//TODO: make parallel
	markers := make([]*Marker, 0)
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		em, errs := p.extractMarkers(pkg.TypesInfo, fn)
		if len(errs) > 0 {
			return nil, errs
		}
		markers = append(markers, em...)
	}
	return markers, []error{}
}

func (p *parser) ExtractArgs(n ast.Node) []ast.Expr {
	args := make([]ast.Expr, 0)
	if fnCall, ok := n.(*ast.CallExpr); ok {
		for _, arg := range fnCall.Args {
			switch v := arg.(type) {
			case *ast.BasicLit:
				args = append(args, v)
			case *ast.Ident:
				args = append(args, v)
			default:
				fmt.Println("Unrecognized type")
			}
		}
	}
	return args
}

func (p *parser) extractMarkers(info *types.Info, fn *ast.FuncDecl) ([]*Marker, []error) {
	if fn.Body == nil {
		return nil, []error{}
	}
	markers := make([]*Marker, 0)
	for _, stmt := range fn.Body.List {
		switch stmt := stmt.(type) {
		case *ast.ExprStmt:
			call, ok := stmt.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			obj := p.qualifiedIdentObject(info, call.Fun)
			if obj == nil ||
				obj.Pkg() == nil ||
				!p.isProcessorImport(obj.Pkg().Path()) ||
				!p.isDSLKeyword(obj.Name()) {
				continue
			}
			markers = append(markers, &Marker{
				Sig:  p.processFuncSignature(obj.(*types.Func)),
				Args: p.ExtractArgs(call),
				Node: stmt,
				Type: EXPR_STMT,
			})
		}
	}

	return markers, []error{}
}

func (p *parser) qualifiedIdentObject(info *types.Info, expr ast.Expr) types.Object {
	switch expr := expr.(type) {
	case *ast.Ident:
		return info.ObjectOf(expr)
	case *ast.SelectorExpr:
		pkgName, ok := expr.X.(*ast.Ident)
		if !ok {
			return nil
		}
		if _, ok := info.ObjectOf(pkgName).(*types.PkgName); !ok {
			return nil
		}
		return info.ObjectOf(expr.Sel)
	default:
		return nil
	}
}

type argSignature struct {
	Name string
	Type types.Type
}

type funcSignature struct {
	Name       string
	Pkg        *types.Package
	Params     []argSignature
	Return     argSignature
	HasVarArgs bool
	HasErr     bool
	HasCleanup bool
}

func (p *parser) processFuncSignature(fn *types.Func) funcSignature {
	sig := fn.Type().(*types.Signature)
	fsig := funcSignature{
		Name:       fn.Name(),
		Pkg:        fn.Pkg(),
		Params:     make([]argSignature, 0),
		HasVarArgs: sig.Variadic(),
	}

	// Check function params
	for i := 0; i < sig.Params().Len(); i++ {
		param := sig.Params().At(i)
		fsig.Params = append(fsig.Params, argSignature{
			Name: param.Name(),
			Type: param.Type(),
		})
	}

	// Check function return values
	results := sig.Results()
	if results.Len() > 0 {
		validateType := func(results *types.Tuple, tp types.Type) bool {
			isIdentical := false
			for i := 0; i < results.Len(); i++ {
				isIdentical = types.Identical(results.At(i).Type(), tp)
			}
			return isIdentical
		}

		// TODO: add validations
		out := results.At(0)
		fsig.Return = argSignature{
			Name: out.Name(),
			Type: out.Type(),
		}
		fsig.HasErr = validateType(results, types.Universe.Lookup("error").Type())
		fsig.HasCleanup = validateType(results, types.NewSignatureType(nil, nil, nil, nil, nil, false))
	}

	return fsig
}

func (p *parser) isDSLKeyword(name string) bool {
	for keyword := range dsl.Markers {
		if name == keyword {
			return true
		}
	}
	return false
}

func (p *parser) isProcessorImport(path string) bool {
	const vendorPart = "vendor/"
	if i := strings.LastIndex(path, vendorPart); i != -1 && (i == 0 || path[i-1] == '/') {
		path = path[i+len(vendorPart):]
	}
	return path == IMPORT_PATH
}
