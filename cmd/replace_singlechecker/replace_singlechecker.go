package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"

	"path/filepath"
)

func main() {
	overlay("")
}

func overlay(dir string) (string, error) {
	var buf bytes.Buffer
	_, err := replaceSingleChecker(&buf)

	if err != nil {
		panic(err)
	}
	return "", nil
}

func replaceSingleChecker(w io.Writer) (string, error) {
	pkgDir := filepath.Join(build.Default.GOPATH, "pkg/mod/golang.org/x/tools@v0.6.0/go/analysis/singlechecker")

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	fmt.Println(pkgs)

	if pkgs["singlechecker"] == nil {
		return "", errors.New("not found singlechecker package")
	}

	for _, file := range pkgs["singlechecker"].Files {
		for _, decl := range file.Decls {
			decl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			searchForExit(decl.Body)
			fmt.Println(decl.Name)
		}
	}
	// pkg, err := build.Default.Import("time", srcDir, 0)
	return "", nil
}

func searchForExit(body *ast.BlockStmt) {
	for _, stmt := range body.List {
		expr, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}

		rhs, ok := expr.X.(*ast.CallExpr)
		if !ok {
			continue
		}

		sel, ok := rhs.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			continue
		}

		if ident.Name == "os" && sel.Sel.Name == "Exit" {
			// TODO: exprをコメントアウトする
			ident.Name = "// os"
		}
	}
}
