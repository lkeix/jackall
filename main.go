package jackall

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
)

type dependenceVec struct {
	from string
	to   string
}

type dependencePair struct {
	in  int
	out int
}

type dependenceVecs []*dependenceVec

func (d dependenceVecs) extractVecEachPackage() map[string]*dependencePair {
	mp := make(map[string]*dependencePair)
	for _, vec := range d {
		p, ok := mp[vec.from]
		if !ok {
			p = &dependencePair{
				in:  0,
				out: 0,
			}
		}

		p.out++
		mp[vec.from] = p

		p, ok = mp[vec.to]
		if !ok {
			p = &dependencePair{
				in:  0,
				out: 0,
			}
		}

		p.in++
		mp[vec.to] = p
	}

	return mp
}
func Run() {
	vec := make(dependenceVecs, 0)

	analyzer := &analysis.Analyzer{
		Name: "Jackall",
		Doc:  "Jackall calculate degree of dependency each packages",
		Run:  wrapRun(&vec),
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}

	singlechecker.Main(analyzer)

	res := vec.extractVecEachPackage()
	for name, r := range res {
		fmt.Printf("degree of dependency in %s package: %.4f\n", name, float64(r.out)/float64(r.in+r.out))
	}

	fmt.Printf("the closer degree of dependency is 1, the more stable package is\n")
	fmt.Printf("the closer degree of dependency is 0, the less stable(unstable) package is\n")
}

// wrapRun bind import dependency for arguments struct
func wrapRun(vec *dependenceVecs) func(pass *analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		// fset := pass.Fset

		for _, f := range pass.Files {
			for _, imprt := range f.Imports {
				name := extractImportPackageName(imprt.Path.Value)
				*vec = append(*vec, &dependenceVec{
					from: f.Name.Name,
					to:   name,
				})
			}
		}

		return nil, nil
	}
}

func extractImportPackageName(path string) string {
	path = strings.ReplaceAll(path, "\"", "")

	reg := regexp.MustCompile(`\/v\d+`)
	if pkgVer := reg.FindString(path); pkgVer != "" {
		path = strings.ReplaceAll(path, pkgVer, "")
	}

	_, pkg := filepath.Split(path)
	return pkg
}
