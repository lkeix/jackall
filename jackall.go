package jackall

import (
	"fmt"
	"go/build"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
)

type dependenceVec struct {
	from  string
	to    string
	toRaw string
}

type dependencePair struct {
	in  int
	out int
}

type dependenceVecs []*dependenceVec

func (d dependenceVecs) contain(pkg string) bool {
	for _, v := range d {
		if v.toRaw == pkg {
			return true
		}
	}
	return false
}

func (d dependenceVecs) filter(pkgs []string) dependenceVecs {
	res := make(dependenceVecs, 0)

	for _, vec := range d {
		isThirdParty := false
		for _, pkg := range pkgs {
			if vec.toRaw == pkg {
				isThirdParty = true
			}
		}
		if !isThirdParty {
			res = append(res, vec)
		}
	}

	return res
}

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

	vec, err := removeStdPackages(vec)
	if err != nil {
		panic(err)
	}

	vec, err = removeThirdPartyPackages(vec)
	if err != nil {
		panic(err)
	}

	res := vec.extractVecEachPackage()
	for name, r := range res {
		fmt.Printf("degree of dependency in %s package: %.4f\n", name, float64(r.out)/float64(r.in+r.out))
	}

	fmt.Printf("the closer degree of dependency is 1, the less stable(unstable) package is\n")
	fmt.Printf("the closer degree of dependency is 0, the more stable package is\n")
}

func removeStdPackages(vecs dependenceVecs) (dependenceVecs, error) {
	res := make(dependenceVecs, 0)
	srcDir := filepath.Join(runtime.GOROOT(), "src")

	for _, vec := range vecs {
		if _, err := build.Default.Import(vec.to, srcDir, 0); err != nil {
			res = append(res, vec)
		}
	}
	return res, nil
}

func removeThirdPartyPackages(vecs dependenceVecs) (dependenceVecs, error) {
	cmd := exec.Command("go", "list", "-m", "all")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	rows := strings.Split(string(out), "\n")
	pkgs := make([]string, 0)
	for _, row := range rows {
		col := strings.Split(row, " ")
		if len(col) == 2 {
			pkgs = append(pkgs, col[0])
		}
	}

	return vecs.filter(pkgs), nil
}

// wrapRun bind import dependency for arguments struct
func wrapRun(vec *dependenceVecs) func(pass *analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		// fset := pass.Fset

		for _, f := range pass.Files {
			for _, imprt := range f.Imports {
				name := extractImportPackageName(imprt.Path.Value)
				*vec = append(*vec, &dependenceVec{
					from:  f.Name.Name,
					to:    name,
					toRaw: strings.ReplaceAll(imprt.Path.Value, "\"", ""),
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
