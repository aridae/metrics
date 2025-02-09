package nomainosexit

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func Test_Analyzer(t *testing.T) {
	// applies an analysis to the packages in ./testdata/ directory
	analysistest.Run(t, analysistest.TestData(), Analyzer,
		"./main-closure-call/...",
		"./main-direct-call/...",
		"./main-indirect-call/...",
		"./non-main-pkg/...",
	)
}
