package main

import (
	"github.com/aridae/go-metrics-store/internal/staticlint"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(staticlint.Analyzers()...)
}
