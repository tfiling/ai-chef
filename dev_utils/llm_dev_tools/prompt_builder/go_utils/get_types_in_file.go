package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func parseArgs() string {
	flag.Parse()
	if flag.NArg() != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s <file>\n", os.Args[0])
		os.Exit(1)
	}
	return flag.Arg(0)
}

func main() {
	filePath := parseArgs()

	// Create file set and parse the Go file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}

	// Find all type declarations
	var typeSpecs []*ast.TypeSpec
	ast.Inspect(file, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					typeSpecs = append(typeSpecs, ts)
				}
			}
			return false
		}
		return true
	})

	if len(typeSpecs) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "No types found in %s\n", filePath)
		os.Exit(1)
	}

	// Print each type declaration with its comments
	for i, typeSpec := range typeSpecs {
		if i > 0 {
			fmt.Println() // Add newline between types
		}

		genDecl := &ast.GenDecl{
			Doc:    typeSpec.Comment,
			Tok:    token.TYPE,
			Lparen: token.NoPos,
			Specs:  []ast.Spec{typeSpec},
			Rparen: token.NoPos,
		}

		err = printer.Fprint(os.Stdout, fset, genDecl)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error printing type declaration: %v\n", err)
			os.Exit(1)
		}
	}
}