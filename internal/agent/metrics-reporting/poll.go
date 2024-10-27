package metricsreporting

import (
	"context"
	"math/rand"
	"runtime"
)

func (a *Agent) poll(_ context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.counters[PollCountMetricName]++

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	a.gauges[AllocMetricName] = gauge(rtm.Alloc)
	a.gauges[TotalAllocMetricName] = gauge(rtm.TotalAlloc)
	a.gauges[SysMetricName] = gauge(rtm.Sys)
	a.gauges[LookupsMetricName] = gauge(rtm.Lookups)
	a.gauges[MallocsMetricName] = gauge(rtm.Mallocs)
	a.gauges[FreesMetricName] = gauge(rtm.Frees)
	a.gauges[HeapAllocMetricName] = gauge(rtm.HeapAlloc)
	a.gauges[HeapSysMetricName] = gauge(rtm.HeapSys)
	a.gauges[HeapIdleMetricName] = gauge(rtm.HeapIdle)
	a.gauges[HeapInuseMetricName] = gauge(rtm.HeapInuse)
	a.gauges[HeapReleasedMetricName] = gauge(rtm.HeapReleased)
	a.gauges[HeapObjectsMetricName] = gauge(rtm.HeapObjects)
	a.gauges[StackInuseMetricName] = gauge(rtm.StackInuse)
	a.gauges[StackSysMetricName] = gauge(rtm.StackSys)
	a.gauges[MSpanInuseMetricName] = gauge(rtm.MSpanInuse)
	a.gauges[MSpanSysMetricName] = gauge(rtm.MSpanSys)
	a.gauges[MCacheInuseMetricName] = gauge(rtm.MCacheInuse)
	a.gauges[MCacheSysMetricName] = gauge(rtm.MCacheSys)
	a.gauges[BuckHashSysMetricName] = gauge(rtm.BuckHashSys)
	a.gauges[GCSysMetricName] = gauge(rtm.GCSys)
	a.gauges[OtherSysMetricName] = gauge(rtm.OtherSys)
	a.gauges[NextGCMetricName] = gauge(rtm.NextGC)
	a.gauges[LastGCMetricName] = gauge(rtm.LastGC)
	a.gauges[PauseTotalNsMetricName] = gauge(rtm.PauseTotalNs)
	a.gauges[NumGCMetricName] = gauge(rtm.NumGC)
	a.gauges[NumForcedGCMetricName] = gauge(rtm.NumForcedGC)
	a.gauges[GCCPUFractionMetricName] = gauge(rtm.GCCPUFraction)
	a.gauges[RandomValueMetricName] = gauge(rand.Float64())
}
