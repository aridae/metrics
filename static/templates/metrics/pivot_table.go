package metrics

import (
	"embed"
)

//go:embed pivot_table.html
var PivotTable embed.FS

const PivotTableHTML = "pivot_table.html"
