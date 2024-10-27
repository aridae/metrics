package metricsreporting

type (
	counter int64
	gauge   float64
)

const (
	AllocMetricName         = "Alloc"
	TotalAllocMetricName    = "TotalAlloc"
	SysMetricName           = "Sys"
	LookupsMetricName       = "Lookups"
	MallocsMetricName       = "Mallocs"
	FreesMetricName         = "Frees"
	HeapAllocMetricName     = "HeapAlloc"
	HeapSysMetricName       = "HeapSys"
	HeapIdleMetricName      = "HeapIdle"
	HeapInuseMetricName     = "HeapInuse"
	HeapReleasedMetricName  = "HeapReleased"
	HeapObjectsMetricName   = "HeapObjects"
	StackInuseMetricName    = "StackInuse"
	StackSysMetricName      = "StackSys"
	MSpanInuseMetricName    = "MSpanInuse"
	MSpanSysMetricName      = "MSpanSys"
	MCacheInuseMetricName   = "MCacheInuse"
	MCacheSysMetricName     = "MCacheSys"
	BuckHashSysMetricName   = "BuckHashSys"
	GCSysMetricName         = "GCSys"
	OtherSysMetricName      = "OtherSys"
	NextGCMetricName        = "NextGC"
	LastGCMetricName        = "LastGC"
	PauseTotalNsMetricName  = "PauseTotalNs"
	NumGCMetricName         = "NumGC"
	NumForcedGCMetricName   = "NumForcedGC"
	GCCPUFractionMetricName = "GCCPUFraction"
	RandomValueMetricName   = "RandomValue"
	PollCountMetricName     = "PollCount"
)
