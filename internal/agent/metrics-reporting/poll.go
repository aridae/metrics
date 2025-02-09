package metricsreporting

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"

	gopscpu "github.com/shirou/gopsutil/cpu"
	gopsmem "github.com/shirou/gopsutil/mem"
)

func pollRuntime(_ context.Context) metricsPack {
	pack := metricsPack{
		gauges:   make(map[string]gauge),
		counters: make(map[string]counter),
	}
	pack.counters[PollCountMetricName]++

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	pack.gauges[AllocMetricName] = gauge(rtm.Alloc)
	pack.gauges[TotalAllocMetricName] = gauge(rtm.TotalAlloc)
	pack.gauges[SysMetricName] = gauge(rtm.Sys)
	pack.gauges[LookupsMetricName] = gauge(rtm.Lookups)
	pack.gauges[MallocsMetricName] = gauge(rtm.Mallocs)
	pack.gauges[FreesMetricName] = gauge(rtm.Frees)
	pack.gauges[HeapAllocMetricName] = gauge(rtm.HeapAlloc)
	pack.gauges[HeapSysMetricName] = gauge(rtm.HeapSys)
	pack.gauges[HeapIdleMetricName] = gauge(rtm.HeapIdle)
	pack.gauges[HeapInuseMetricName] = gauge(rtm.HeapInuse)
	pack.gauges[HeapReleasedMetricName] = gauge(rtm.HeapReleased)
	pack.gauges[HeapObjectsMetricName] = gauge(rtm.HeapObjects)
	pack.gauges[StackInuseMetricName] = gauge(rtm.StackInuse)
	pack.gauges[StackSysMetricName] = gauge(rtm.StackSys)
	pack.gauges[MSpanInuseMetricName] = gauge(rtm.MSpanInuse)
	pack.gauges[MSpanSysMetricName] = gauge(rtm.MSpanSys)
	pack.gauges[MCacheInuseMetricName] = gauge(rtm.MCacheInuse)
	pack.gauges[MCacheSysMetricName] = gauge(rtm.MCacheSys)
	pack.gauges[BuckHashSysMetricName] = gauge(rtm.BuckHashSys)
	pack.gauges[GCSysMetricName] = gauge(rtm.GCSys)
	pack.gauges[OtherSysMetricName] = gauge(rtm.OtherSys)
	pack.gauges[NextGCMetricName] = gauge(rtm.NextGC)
	pack.gauges[LastGCMetricName] = gauge(rtm.LastGC)
	pack.gauges[PauseTotalNsMetricName] = gauge(rtm.PauseTotalNs)
	pack.gauges[NumGCMetricName] = gauge(rtm.NumGC)
	pack.gauges[NumForcedGCMetricName] = gauge(rtm.NumForcedGC)
	pack.gauges[GCCPUFractionMetricName] = gauge(rtm.GCCPUFraction)
	pack.gauges[RandomValueMetricName] = gauge(rand.Float64())

	return pack
}

func pollGopsutil(ctx context.Context) (metricsPack, error) {
	pack := metricsPack{
		gauges: make(map[string]gauge),
	}

	memstats, err := gopsmem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return metricsPack{}, fmt.Errorf("failed to get memory stats: %w", err)
	}

	cpustats, err := gopscpu.CountsWithContext(ctx, true)
	if err != nil {
		return metricsPack{}, fmt.Errorf("failed to get cpu stats: %w", err)
	}

	pack.gauges[TotalMemoryMetricName] = gauge(memstats.Total)
	pack.gauges[FreeMemoryMetricName] = gauge(memstats.Free)
	pack.gauges[CPUutilization1MetricName] = gauge(cpustats)

	return pack, nil
}
