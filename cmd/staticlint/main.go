package main

import (
	"github.com/aridae/go-metrics-store/internal/staticlint"
	_ "github.com/aridae/go-metrics-store/pkg/logger"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(staticlint.Analyzers()...)
}
