package main

import (
	"github.com/MartinKuzma/stutter/pkg/stutter"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(stutter.NewAnalyzer())
}
