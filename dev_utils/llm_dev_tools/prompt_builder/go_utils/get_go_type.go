package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

func parseArgs() (filePath, typeName string) {
	flag.Parse()
	if flag.NArg() != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s <file>:<type>\n", os.Args[0])
		os.Exit(1)
	}

	arg := flag.Arg(0)
	parts := strings.Split(arg, ":")
	if len(parts) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Invalid input format. Expected <file>:<type>, got %s\n", arg)
		os.Exit(1)
	}
	filePath, typeName = parts[0], parts[1]
	return
}

func main() {
	filePath, typeName := parseArgs()

	// Create file set and parse the Go file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}

	// Find the type declaration
	var typeSpec *ast.TypeSpec
	ast.Inspect(file, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if ts.Name.Name == typeName {
				typeSpec = ts
				return false
			}
		}
		return true
	})

	if typeSpec == nil {
		_, _ = fmt.Fprintf(os.Stderr, "Type %q not found in %s\n", typeName, filePath)
		os.Exit(1)
	}

	// Print the type declaration with its comments
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
