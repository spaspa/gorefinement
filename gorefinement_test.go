package gorefinement_test

import (
	"testing"

	"github.com/spaspa/gorefinement"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, gorefinement.Analyzer, "a")
}
