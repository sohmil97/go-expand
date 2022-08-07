package parser

import (
	"context"
	"go/ast"
	goTypes "go/types"
	"strings"
	"x/internal/macro"
	"x/internal/types"

	"golang.org/x/tools/go/packages"
)

const IMPORT_PATH = "x/dsl"

const (
	FUNC_STMT = iota
)

type FileSpec struct {
	Pkg    string
	Syntax *ast.File

	Markers []*MarkerSpec
	Imports []*types.ImportSpec
}

type MarkerSpec struct {
	FunctionMarker *FunctionMarker

	Node ast.Node
	Type int
}

type Parser interface {
	Load(ctx context.Context, directory string) ([]*packages.Package, []error)
	ParseFile(pkg *packages.Package, f *ast.File) (*FileSpec, []error)
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

func (p *parser) ParseFile(pkg *packages.Package, syntax *ast.File) (*FileSpec, []error) {
	//TODO: make parallel
	markers := make([]*MarkerSpec, 0)
	for _, decl := range syntax.Decls {
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

	imports := make([]*types.ImportSpec, 0)
	for _, impt := range syntax.Imports {
		impSpec := &types.ImportSpec{
			Path: impt.Path.Value,
		}
		if impt.Name != nil {
			impSpec.Name = impt.Name.Name
		}
		imports = append(imports, impSpec)
	}

	return &FileSpec{
		Pkg:     pkg.Name,
		Syntax:  syntax,
		Markers: markers,
		Imports: imports,
	}, []error{}
}

func (p *parser) extractMarkers(info *goTypes.Info, fn *ast.FuncDecl) ([]*MarkerSpec, []error) {
	if fn.Body == nil {
		return nil, []error{}
	}
	markers := make([]*MarkerSpec, 0)
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
			markers = append(markers, &MarkerSpec{
				FunctionMarker: p.processFuncMarker(call, obj.(*goTypes.Func)),
				Node:           stmt,
				Type:           FUNC_STMT,
			})
		}
	}

	return markers, []error{}
}

func (p *parser) qualifiedIdentObject(info *goTypes.Info, expr ast.Expr) goTypes.Object {
	switch expr := expr.(type) {
	case *ast.Ident:
		return info.ObjectOf(expr)
	case *ast.SelectorExpr:
		pkgName, ok := expr.X.(*ast.Ident)
		if !ok {
			return nil
		}
		if _, ok := info.ObjectOf(pkgName).(*goTypes.PkgName); !ok {
			return nil
		}
		return info.ObjectOf(expr.Sel)
	default:
		return nil
	}
}

func (p *parser) isDSLKeyword(name string) bool {
	for keyword := range macro.Markers {
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
