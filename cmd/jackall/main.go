package main

import (
	"github.com/lkeix/jackall"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(jackall.Analyzer)
}
