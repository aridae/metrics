package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

const (
	collectInterval = 2 * time.Second
	reportInterval  = 10 * time.Second

	baseEndpointURL = "http://localhost:8080/update"

	counterTypeURLParam = "counter"
	gaugeTypeURLParam   = "gauge"
)

type gauge float64
type counter int64

func main() {
	collectTick := time.NewTicker(collectInterval)
	reportTick := time.NewTicker(reportInterval)

	httpClient := &http.Client{}

	pollCounter := counter(0)
	gaugeMetrics := make(map[string]gauge)

	for {
		select {
		case <-collectTick.C:
			log.Printf("starting metrics collection routine <now:%s>\n", time.Now().UTC())
			pollCounter++
			collectMetrics(gaugeMetrics)
		case <-reportTick.C:
			log.Printf("starting metrics report routine <now:%s>\n", time.Now().UTC())
			reportMetrics(httpClient, gaugeMetrics, pollCounter)
		}
	}
}

func reportMetrics(client *http.Client, gaugeMetrics map[string]gauge, pollCountMetric counter) {
	for metricName, metricVal := range gaugeMetrics {
		metric := fmt.Sprintf("/%s/%s/%v", gaugeTypeURLParam, metricName, metricVal)
		reportMetric(client, metric)
	}
	metric := fmt.Sprintf("/%s/%s/%v", counterTypeURLParam, PollCountMetricName, pollCountMetric)
	reportMetric(client, metric)
}

func reportMetric(client *http.Client, metric string) {
	serverURL, _ := url.JoinPath(baseEndpointURL, metric)

	data := []byte("")
	req, err := http.NewRequest(http.MethodPost, serverURL, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("failed to build http request: %v", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to do http request: %v", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Fatalf("failed to close resp body: %v", err)
		}
	}()

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		log.Fatalf("failed to read body: %v", err)
	}
}

func collectMetrics(metrics map[string]gauge) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	metrics[AllocMetricName] = gauge(rtm.Alloc)
	metrics[TotalAllocMetricName] = gauge(rtm.TotalAlloc)
	metrics[SysMetricName] = gauge(rtm.Sys)
	metrics[LookupsMetricName] = gauge(rtm.Lookups)
	metrics[MallocsMetricName] = gauge(rtm.Mallocs)
	metrics[FreesMetricName] = gauge(rtm.Frees)
	metrics[HeapAllocMetricName] = gauge(rtm.HeapAlloc)
	metrics[HeapSysMetricName] = gauge(rtm.HeapSys)
	metrics[HeapIdleMetricName] = gauge(rtm.HeapIdle)
	metrics[HeapInuseMetricName] = gauge(rtm.HeapInuse)
	metrics[HeapReleasedMetricName] = gauge(rtm.HeapReleased)
	metrics[HeapObjectsMetricName] = gauge(rtm.HeapObjects)
	metrics[StackInuseMetricName] = gauge(rtm.StackInuse)
	metrics[StackSysMetricName] = gauge(rtm.StackSys)
	metrics[MSpanInuseMetricName] = gauge(rtm.MSpanInuse)
	metrics[MSpanSysMetricName] = gauge(rtm.MSpanSys)
	metrics[MCacheInuseMetricName] = gauge(rtm.MCacheInuse)
	metrics[MCacheSysMetricName] = gauge(rtm.MCacheSys)
	metrics[BuckHashSysMetricName] = gauge(rtm.BuckHashSys)
	metrics[GCSysMetricName] = gauge(rtm.GCSys)
	metrics[OtherSysMetricName] = gauge(rtm.OtherSys)
	metrics[NextGCMetricName] = gauge(rtm.NextGC)
	metrics[LastGCMetricName] = gauge(rtm.LastGC)
	metrics[PauseTotalNsMetricName] = gauge(rtm.PauseTotalNs)
	metrics[NumGCMetricName] = gauge(rtm.NumGC)
	metrics[NumForcedGCMetricName] = gauge(rtm.NumForcedGC)
	metrics[GCCPUFractionMetricName] = gauge(rtm.GCCPUFraction)
	metrics[RandomValueMetricName] = gauge(rand.Float64())
}

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
