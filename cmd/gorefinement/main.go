package main

import (
	"poyo_analyser/gorefinement"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(gorefinement.Analyzer) }

