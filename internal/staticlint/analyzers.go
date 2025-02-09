package staticlint

import (
	"github.com/aridae/go-metrics-store/internal/staticlint/nomainosexit"
	"github.com/fatih/errwrap/errwrap"
	"golang.org/x/tools/go/analysis"
)

func Analyzers() []*analysis.Analyzer {
	staticChecksAnalyzers := staticChecks()
	standardPassesAnalyzers := standardPasses()

	allAnalyzers := make([]*analysis.Analyzer, 0, len(standardPassesAnalyzers)+len(staticChecksAnalyzers)+2)

	// все анализаторы класса SA и выборочные анализаторы класса ST пакета staticcheck.io
	allAnalyzers = append(allAnalyzers, staticChecksAnalyzers...)

	// все стандартные статические анализаторы пакета golang.org/x/tools/go/analysis/passes
	allAnalyzers = append(allAnalyzers, standardPassesAnalyzers...)

	// публичный анализатор "github.com/fatih/errwrap/errwrap"
	allAnalyzers = append(allAnalyzers, errwrap.Analyzer)

	// реализованный с использованием пакета ast анализатор nomainosexit
	allAnalyzers = append(allAnalyzers, nomainosexit.Analyzer)

	return allAnalyzers
}
