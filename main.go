package jackall

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
)

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

func Run() {
	mapper := &dependencyMapper{
		packages: make(map[string]files),
	}

	analyzer := &analysis.Analyzer{
		Name: "Jackall",
		Doc:  "Jackall calculate degree of dependency each packages",
		Run:  wrapRun(mapper),
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}

	singlechecker.Main(analyzer)

}

type dependencyMapper struct {
	packages map[string]files
}

// wrapRun bind import dependency for arguments struct
func wrapRun(deps *dependencyMapper) func(pass *analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		// inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		// if !ok {
		// 	panic("failed to convert inspect Analyzer to Inspector")
		// }

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

		return nil, nil
	}
}
