package processor

import (
	"context"
	"errors"
	"fmt"
	"go/format"
	"path/filepath"
	"strings"

	"x/internal/processor/generator"
	"x/internal/processor/parser"
)

type GenerateResult struct {
	PkgPath    string
	OutputPath string
	Content    []byte
	Errs       []error
}

type Processor interface {
	Generate(ctx context.Context) ([]GenerateResult, []error)
}

type processor struct {
	directory string
	env       []string

	parser parser.Parser
}

type ProcessorOpt func(ps *processor)

func WithDirectory(directory string) ProcessorOpt {
	return func(p *processor) {
		p.directory = directory
	}
}

func WithEnv(env []string) ProcessorOpt {
	return func(p *processor) {
		p.env = env
	}
}

func New(opts ...ProcessorOpt) Processor {
	processor := new(processor)
	for _, opt := range opts {
		opt(processor)
	}

	// config processor parser
	processor.parser = parser.New(
		parser.WithEnv(processor.env),
	)

	return processor
}

// TODO: make parallel
func (p *processor) Generate(ctx context.Context) ([]GenerateResult, []error) {
	pkgs, errs := p.parser.Load(ctx, p.directory)
	if len(errs) > 0 {
		return nil, errs
	}
	results := make([]GenerateResult, 0)
	for _, pkg := range pkgs {
		outDir, err := p.detectOutputDir(pkg.GoFiles)
		if err != nil {
			fmt.Print(err)
		}
		for j, syntax := range pkg.Syntax {
			result := GenerateResult{
				PkgPath: pkg.PkgPath,
			}
			fileNameSegments := strings.Split(strings.Split(pkg.GoFiles[j], ".go")[0], "/")
			fileName := fileNameSegments[len(fileNameSegments)-1]
			result.OutputPath = fmt.Sprintf("%s/%s_gen.go", outDir, fileName)
			fileSpec, errs := p.parser.ParseFile(pkg, syntax)
			if len(errs) > 0 {
				result.Errs = errs
				continue
			}
			p.processFile(fileSpec)
			generator := generator.New(pkg)
			goSrc := generator.GenerateSource("")
			fmtSrc, err := format.Source(goSrc)
			if err != nil {
				// This is likely a bug from a poorly generated source file.
				// Add an error but also the unformatted source.
				result.Errs = append(result.Errs, err)
			} else {
				goSrc = fmtSrc
			}
			result.Content = goSrc

			fmt.Printf("file: %s,\n markers: %#v\n\n", fileName, fileSpec.Markers[0].FunctionMarker)
			results = append(results, result)
		}
	}
	return results, nil
}

func (p *processor) processFile(fileSpec *parser.FileSpec) {
	// for _, m := range fileSpec.Markers {
	// dslProcessor := dsl.Markers[m.Sig.Name]
	// dslProcessor(m.Node, dslImp.NodeSpec{
	// Args: m.Args,
	// })
	// astutil.Apply(m.Node,
	// 	func(c *astutil.Cursor) bool {
	// 		c.Replace(pn)
	// 		return false
	// 	},
	// 	func(c *astutil.Cursor) bool {
	// 		return true
	// 	},
	// )
	// }
}

func (p *processor) detectOutputDir(paths []string) (string, error) {
	if len(paths) == 0 {
		return "", errors.New("no files to derive output directory from")
	}
	dir := filepath.Dir(paths[0])
	for _, p := range paths[1:] {
		if dir2 := filepath.Dir(p); dir2 != dir {
			return "", fmt.Errorf("found conflicting directories %q and %q", dir, dir2)
		}
	}
	return dir, nil
}
