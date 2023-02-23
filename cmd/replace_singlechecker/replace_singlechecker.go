package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"io"

	"path/filepath"
)

func main() {
	var buf bytes.Buffer
	replaceSingleChecker(&buf)
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
	pkgDir := filepath.Join(build.Default.GOPATH, "pkg/mod/golang.org/x/tools@v0.6.0/go/analysis")
	pkg, err := build.Default.Import("singlechecker", pkgDir, 0)
	fmt.Println(pkg)
	fmt.Println(err)
	if err != nil {
		return "", err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkg.Dir, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	if pkgs["singlechecker"] == nil {
		return "", errors.New("not found singlechecker package")
	}

	for _, file := range pkgs["singlechecker"].Files {
		fmt.Println(file)
	}
	// pkg, err := build.Default.Import("time", srcDir, 0)
	return "", nil
}
