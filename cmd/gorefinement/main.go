package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"
	"poyo_analyser/gorefinement"
)

func main() { unitchecker.Main(gorefinement.Analyzer) }
