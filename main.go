package jackall

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "Jackall",
	Doc:  "Jackall calculate degree of dependency each packages",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

type Indicate struct {
	mp map[token.Pos]bool
}

type file struct {
	name          string
	dependencyIn  int64
	dependencyOut int64
}

type files []*file

func (fs files) contains(target string) bool {
	for _, f := range fs {
		if f.name == target {
			return true
		}
	}
	return false
}

type packages map[string]files

var pkgs = make(packages)

func run(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		panic("failed to convert inspect Analyzer to Inspector")
	}

	fset := pass.Fset

	for _, f := range pass.Files {
		fmt.Printf("dependence package of %s\n", fset.File(f.Pos()).Name())

		fout := int64(0)
		for _, imprt := range f.Imports {

			fmt.Printf("\t%s\n", imprt.Path.Value)
			fout++
		}

		if _, ok := pkgs[f.Name.Name]; !ok {
			pkgs[f.Name.Name] = files{
				{
					name:          fset.File(f.Pos()).Name(),
					dependencyIn:  0,
					dependencyOut: fout,
				},
			}
		} else {
			pkg := pkgs[f.Name.Name]
			pkg = append(pkg, &file{
				name:          fset.File(f.Pos()).Name(),
				dependencyIn:  0,
				dependencyOut: fout,
			})
		}
	}

	inspect.Nodes(nil, func(n ast.Node, push bool) bool {
		calculateDepenedencyEachPackages(n)
		return true
	})

	fmt.Println(pkgs)
	return nil, nil
}

func calculatePackageLevelDepenedencyEachPackages() {

}

func calculateDepenedencyEachPackages(n ast.Node) {
}
