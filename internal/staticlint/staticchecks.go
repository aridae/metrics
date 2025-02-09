package staticlint

import (
	"github.com/aridae/go-metrics-store/pkg/logger"
	"github.com/aridae/go-metrics-store/pkg/slice"
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

// staticChecks возвращает все анализаторы класса SA и выборочные анализаторы класса simple пакета staticcheck.io
func staticChecks() []*analysis.Analyzer {
	analyzers := make([]*analysis.Analyzer, 0, len(_securityAnalysisStaticChecks)+len(_simpleAnalysisStaticChecks))

	for saName := range _securityAnalysisStaticChecks {
		lintAnalyzer, ok := _allSAChecks[saName]
		if !ok {
			logger.Warnf("static check analyzer for security analysis '%s' does not exist, skipping", saName)
			continue
		}
		analyzers = append(analyzers, lintAnalyzer.Analyzer)
	}

	for sName := range _simpleAnalysisStaticChecks {
		lintAnalyzer, ok := _allSimpleChecks[sName]
		if !ok {
			logger.Warnf("static check gosimple analyzer '%s' does not exist, skipping", sName)
			continue
		}
		analyzers = append(analyzers, lintAnalyzer.Analyzer)
	}

	return analyzers
}

// _allSAChecks все анализаторы класса SA пакета staticcheck.io
var _allSAChecks = slice.KeyBy(staticcheck.Analyzers, func(elem *lint.Analyzer) string { return elem.Analyzer.Name })

// _allSimpleChecks все анализаторы класса gosimple пакета staticcheck.io
var _allSimpleChecks = slice.KeyBy(simple.Analyzers, func(elem *lint.Analyzer) string { return elem.Analyzer.Name })

// _securityAnalysisStaticChecks имена используемых анализаторов класса SA пакета staticcheck.io
var _securityAnalysisStaticChecks = map[string]struct{}{
	"SA1000": {},
	"SA1001": {},
	"SA1002": {},
	"SA1003": {},
	"SA1004": {},
	"SA1005": {},
	"SA1006": {},
	"SA1007": {},
	"SA1008": {},
	"SA1010": {},
	"SA1011": {},
	"SA1012": {},
	"SA1013": {},
	"SA1014": {},
	"SA1015": {},
	"SA1016": {},
	"SA1017": {},
	"SA1018": {},
	"SA1019": {},
	"SA1020": {},
	"SA1021": {},
	"SA1023": {},
	"SA1024": {},
	"SA1025": {},
	"SA1026": {},
	"SA1027": {},
	"SA1028": {},
	"SA1029": {},
	"SA1030": {},
	"SA1031": {},
	"SA1032": {},
	"SA2000": {},
	"SA2001": {},
	"SA2002": {},
	"SA2003": {},
	"SA3000": {},
	"SA3001": {},
	"SA4000": {},
	"SA4001": {},
	"SA4003": {},
	"SA4004": {},
	"SA4005": {},
	"SA4006": {},
	"SA4008": {},
	"SA4009": {},
	"SA4010": {},
	"SA4011": {},
	"SA4012": {},
	"SA4013": {},
	"SA4014": {},
	"SA4015": {},
	"SA4016": {},
	"SA4017": {},
	"SA4018": {},
	"SA4019": {},
	"SA4020": {},
	"SA4021": {},
	"SA4022": {},
	"SA4023": {},
	"SA4024": {},
	"SA4025": {},
	"SA4026": {},
	"SA4027": {},
	"SA4028": {},
	"SA4029": {},
	"SA4030": {},
	"SA4031": {},
	"SA4032": {},
	"SA5000": {},
	"SA5001": {},
	"SA5002": {},
	"SA5003": {},
	"SA5004": {},
	"SA5005": {},
	"SA5007": {},
	"SA5008": {},
	"SA5009": {},
	"SA5010": {},
	"SA5011": {},
	"SA5012": {},
	"SA6000": {},
	"SA6001": {},
	"SA6002": {},
	"SA6003": {},
	"SA6005": {},
	"SA6006": {},
	"SA9001": {},
	"SA9002": {},
	"SA9003": {},
	"SA9004": {},
	"SA9005": {},
	"SA9006": {},
	"SA9007": {},
	"SA9008": {},
	"SA9009": {},
}

// _simpleAnalysisStaticChecks имена выборочных анализаторов класса gosimple пакета staticcheck.io
var _simpleAnalysisStaticChecks = map[string]struct{}{
	"S1000": {},
	"S1001": {},
	"S1003": {},
	"S1010": {},
	"S1012": {},
	"S1018": {},
	"S1019": {},
}
