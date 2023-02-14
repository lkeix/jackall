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

func run(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		panic("failed to convert inspect Analyzer to Inspector")
	}

	fmt.Println(pass.AllPackageFacts)

	inspect.Nodes(nil, func(n ast.Node, push bool) bool {
		calculateDepenedencyEachPackages(n)
		return true
	})

	return nil, nil
}

func calculateDepenedencyEachPackages(n ast.Node) {
}
