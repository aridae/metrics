package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"
)

type (
	counter int64
	gauge   float64
)

const (
	baseURLPath = "/update"
	counterType = "counter"
	gaugeType   = "gauge"
)

var (
	pollInterval   int64
	reportInterval int64
	address        string
	useOldHandler  bool
)

func init() {
	flag.Int64Var(&reportInterval, "r", 10, "частота отправки метрик на сервер (по умолчанию 10 секунд)")
	flag.Int64Var(&pollInterval, "p", 2, "частота опроса метрик из пакета runtime (по умолчанию 2 секунды)")
	flag.StringVar(&address, "a", "localhost:8080", "адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080")
	flag.BoolVar(&useOldHandler, "o", false, "Использовать старый эндпоинт [/update/<type>/<name>/<value>] для сохранения метрики (по умолчанию false)")
}

func main() {
	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		address = envAddress
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		parsedEnv, err := strconv.ParseInt(envReportInterval, 10, 64)
		if err != nil {
			log.Fatalf("invalid REPORT_INTERVAL environment variable, int64 value expected: %v", err)
		}
		reportInterval = parsedEnv
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		parsedEnv, err := strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			log.Fatalf("invalid POLL_INTERVAL environment variable, int64 value expected: %v", err)
		}
		pollInterval = parsedEnv
	}

	pollTick := time.NewTicker(time.Duration(pollInterval) * time.Second)
	reportTick := time.NewTicker(time.Duration(reportInterval) * time.Second)

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
		reportMetric(client, gaugeType, metricName, metricVal)
	}
	reportMetric(client, counterType, PollCountMetricName, pollCountMetric)
}

func reportMetric(client *http.Client, metricType, metricName string, metricVal any) {
	if useOldHandler {
		reportMetricWithURLPath(client, metricType, metricName, metricVal)
		return
	}

	reportMetricWithJSONPayload(client, metricType, metricName, metricVal)
}

func reportMetricWithJSONPayload(client *http.Client, metricType, metricName string, metricVal any) {
	serverURL, _ := url.JoinPath("http://"+address, baseURLPath)

	metric, err := buildMetricJSONPayload(metricType, metricName, metricVal)
	if err != nil {
		log.Fatalf("failed to build metric json-serializable struct: %v", err)
	}

	metricBytes, err := json.Marshal(metric)
	if err != nil {
		log.Fatalf("failed to marshal request body: %v", err)
	}

	mustDoRequest(client, http.MethodPost, serverURL, metricBytes, "application/json")
}

func reportMetricWithURLPath(client *http.Client, metricType, metricName string, metricVal any) {
	metricURLPath := fmt.Sprintf("/%s/%s/%v", metricType, metricName, metricVal)

	serverURL, _ := url.JoinPath("http://"+address, baseURLPath, metricURLPath)

	mustDoRequest(client, http.MethodPost, serverURL, nil, "text/plain")
}

func mustDoRequest(client *http.Client, method string, url string, body []byte, contentType string) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("failed to build http request: %v", err)
	}

	req.Header.Set("Content-Type", contentType)
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

func buildMetricJSONPayload(
	mtype string,
	name string,
	val any,
) (httpmodels.Metric, error) {
	switch mtype {
	case counterType:
		counterVal, ok := val.(counter)
		if !ok {
			return httpmodels.Metric{}, fmt.Errorf("value is not int64")
		}
		int64Val := int64(counterVal)
		return httpmodels.Metric{
			ID:    name,
			MType: mtype,
			Delta: &int64Val,
		}, nil
	case gaugeType:
		gaugeVal, ok := val.(gauge)
		if !ok {
			return httpmodels.Metric{}, fmt.Errorf("value is not float64")
		}
		float64Val := float64(gaugeVal)
		return httpmodels.Metric{
			ID:    name,
			MType: mtype,
			Value: &float64Val,
		}, nil
	default:
		return httpmodels.Metric{}, fmt.Errorf("unsupported metric type: %s", mtype)
	}
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
