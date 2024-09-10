package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

type (
	counter int64
	gauge   float64
)

const (
	baseURLPath         = "/update"
	counterTypeURLParam = "counter"
	gaugeTypeURLParam   = "gauge"
)

var (
	pollInterval   *int64
	reportInterval *int64
	address        *string
)

func init() {
	reportInterval = flag.Int64("r", 10, "частота отправки метрик на сервер (по умолчанию 10 секунд)")
	pollInterval = flag.Int64("p", 2, "частота опроса метрик из пакета runtime (по умолчанию 2 секунды)")
	address = flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080")
}

func main() {
	flag.Parse()

	pollTick := time.NewTicker(time.Second * time.Duration(*pollInterval))
	reportTick := time.NewTicker(time.Second * time.Duration(*reportInterval))

	httpClient := &http.Client{}

	pollCounter := counter(0)
	gaugeMetrics := make(map[string]gauge)

	for {
		select {
		case <-pollTick.C:
			log.Printf("starting metrics collection routine <now:%s>\n", time.Now().UTC())
			pollCounter++
			pollMetrics(gaugeMetrics)
		case <-reportTick.C:
			log.Printf("starting metrics report routine <now:%s>\n", time.Now().UTC())
			reportMetrics(httpClient, gaugeMetrics, pollCounter)
		}
	}
}

func reportMetrics(client *http.Client, gaugeMetrics map[string]gauge, pollCountMetric counter) {
	for metricName, metricVal := range gaugeMetrics {
		metricURLPath := fmt.Sprintf("/%s/%s/%v", gaugeTypeURLParam, metricName, metricVal)
		reportMetric(client, metricURLPath)
	}
	metricURLPath := fmt.Sprintf("/%s/%s/%v", counterTypeURLParam, PollCountMetricName, pollCountMetric)
	reportMetric(client, metricURLPath)
}

func reportMetric(client *http.Client, metricURLPath string) {
	serverURL, _ := url.JoinPath("http://"+*address, baseURLPath, metricURLPath)

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

func pollMetrics(metrics map[string]gauge) {
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
