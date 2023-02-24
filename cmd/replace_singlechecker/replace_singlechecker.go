package main

import (
	"bytes"
	_ "embed"

	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"

	"path/filepath"

	_ "golang.org/x/tools/go/analysis/singlechecker"
)

//go:embed _partials/replace_singlechecker.go
var replace string

func main() {
	dir := filepath.Join(build.Default.GOPATH, "replace_singlechecker")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		panic(err)
	}

	path, err := overlay(dir)
	if err != nil {
		panic(err)
	}

	fmt.Println(path)
}

func overlay(dir string) (string, error) {
	var buf bytes.Buffer
	old, err := replaceSingleChecker(&buf)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(&buf, replace)

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	n := filepath.Join(dir, fmt.Sprintf("replaced_singlechecker.go"))

	if err := os.WriteFile(n, src, 0o600); err != nil {
		return "", err
	}

	v := struct {
		Replace map[string]string
	}{
		Replace: map[string]string{
			old: n,
		},
	}

	overlayPath := filepath.Join(dir, "overlay_singlechecker.json")
	var jsonBytes bytes.Buffer
	if err := json.NewEncoder(&jsonBytes).Encode(v); err != nil {
		return "", err
	}

	if err := os.WriteFile(overlayPath, jsonBytes.Bytes(), 0o600); err != nil {
		return "", err
	}

	return overlayPath, nil
}

func replaceSingleChecker(w io.Writer) (string, error) {
	pkgDir := filepath.Join(build.Default.GOPATH, "pkg/mod/golang.org/x/tools@v0.6.0/go/analysis/singlechecker")

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	if pkgs["singlechecker"] == nil {
		return "", errors.New("not found singlechecker package")
	}

	// refer to https://github.com/tenntenn/testtime/blob/55bcd1f05226591251b90f09b982e7454076e3ab/cmd/testtime/overlay.go
	var (
		path   string
		syntax *ast.File
	)
	for name, file := range pkgs["singlechecker"].Files {
		for _, decl := range file.Decls {
			decl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			if decl.Name.Name == "Main" {
				decl.Name.Name = "_Main"
				path = name
				syntax = file
			}
		}
	}
	if err := format.Node(w, fset, syntax); err != nil {
		return "", err
	}

	return path, nil
}
