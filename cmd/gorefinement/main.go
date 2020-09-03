package main

import (
	"github.com/spaspa/gorefinement"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(gorefinement.Analyzer) }
