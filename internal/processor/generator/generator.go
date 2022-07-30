package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"sort"

	"golang.org/x/tools/go/packages"
)

type importSpec struct {
	name    string
	differs bool
}

type generator struct {
	pkg         *packages.Package
	buf         bytes.Buffer
	imports     map[string]importSpec
	anonImports map[string]bool
	values      map[ast.Expr]string
}

func New(pkg *packages.Package) *generator {
	return &generator{
		pkg:         pkg,
		anonImports: make(map[string]bool),
		imports:     make(map[string]importSpec),
		values:      make(map[ast.Expr]string),
	}
}

func (g *generator) GenerateSource(tags string) []byte {
	// if g.buf.Len() == 0 {
	// 	return nil
	// }
	var buf bytes.Buffer
	if len(tags) > 0 {
		tags = fmt.Sprintf(" gen -tags \"%s\"", tags)
	}
	buf.WriteString("// Code generated by Wire. DO NOT EDIT.\n\n")
	buf.WriteString("//go:generate go run -mod=mod github.com/google/wire/cmd/wire" + tags + "\n")
	buf.WriteString("//+build !wireinject\n\n")
	buf.WriteString("package ")
	buf.WriteString(g.pkg.Name)
	buf.WriteString("\n\n")
	if len(g.imports) > 0 {
		buf.WriteString("import (\n")
		imps := make([]string, 0, len(g.imports))
		for path := range g.imports {
			imps = append(imps, path)
		}
		sort.Strings(imps)
		for _, path := range imps {
			// Omit the local package identifier if it matches the package name.
			info := g.imports[path]
			if info.differs {
				fmt.Fprintf(&buf, "\t%s %q\n", info.name, path)
			} else {
				fmt.Fprintf(&buf, "\t%q\n", path)
			}
		}
		buf.WriteString(")\n\n")
	}
	if len(g.anonImports) > 0 {
		buf.WriteString("import (\n")
		anonImps := make([]string, 0, len(g.anonImports))
		for path := range g.anonImports {
			anonImps = append(anonImps, path)
		}
		sort.Strings(anonImps)

		for _, path := range anonImps {
			fmt.Fprintf(&buf, "\t_ %s\n", path)
		}
		buf.WriteString(")\n\n")
	}
	return buf.Bytes()
}
