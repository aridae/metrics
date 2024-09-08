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
	pollInterval    = 2 * time.Second
	reportInterval  = 10 * time.Second
	baseEndpointURL = "http://localhost:8080/update"
)

type gauge float64
type counter int64

func main() {
	pollTick := time.NewTicker(pollInterval)
	reportTick := time.NewTicker(reportInterval)

	httpClient := &http.Client{}

	pollCounter := counter(0)
	gaugeMetrics := make(map[string]gauge)

	for {
		select {
		case <-pollTick.C:
			pollCounter++
			collectMetrics(gaugeMetrics)
		case <-reportTick.C:
			for k, v := range gaugeMetrics {
				metric := fmt.Sprint("/gauge/", k, "/", v)
				reportMetric(httpClient, metric)
			}
			metric := fmt.Sprint("/counter/PollCount/", pollCounter)
			reportMetric(httpClient, metric)
		}
	}
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
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		log.Fatalf("failed to read body: %v", err)
	}
}

func collectMetrics(metrics map[string]gauge) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	metrics["Alloc"] = gauge(rtm.Alloc)
	metrics["TotalAlloc"] = gauge(rtm.TotalAlloc)
	metrics["Sys"] = gauge(rtm.Sys)
	metrics["Lookups"] = gauge(rtm.Lookups)
	metrics["Mallocs"] = gauge(rtm.Mallocs)
	metrics["Frees"] = gauge(rtm.Frees)
	metrics["HeapAlloc"] = gauge(rtm.HeapAlloc)
	metrics["HeapSys"] = gauge(rtm.HeapSys)
	metrics["HeapIdle"] = gauge(rtm.HeapIdle)
	metrics["HeapInuse"] = gauge(rtm.HeapInuse)
	metrics["HeapReleased"] = gauge(rtm.HeapReleased)
	metrics["HeapObjects"] = gauge(rtm.HeapObjects)
	metrics["StackInuse"] = gauge(rtm.StackInuse)
	metrics["StackSys"] = gauge(rtm.StackSys)
	metrics["MSpanInuse"] = gauge(rtm.MSpanInuse)
	metrics["MSpanSys"] = gauge(rtm.MSpanSys)
	metrics["MCacheInuse"] = gauge(rtm.MCacheInuse)
	metrics["MCacheSys"] = gauge(rtm.MCacheSys)
	metrics["BuckHashSys"] = gauge(rtm.BuckHashSys)
	metrics["GCSys"] = gauge(rtm.GCSys)
	metrics["OtherSys"] = gauge(rtm.OtherSys)
	metrics["NextGC"] = gauge(rtm.NextGC)
	metrics["LastGC"] = gauge(rtm.LastGC)
	metrics["PauseTotalNs"] = gauge(rtm.PauseTotalNs)
	metrics["NumGC"] = gauge(rtm.NumGC)
	metrics["NumForcedGC"] = gauge(rtm.NumForcedGC)
	metrics["GCCPUFraction"] = gauge(rtm.GCCPUFraction)

	metrics["RandomValue"] = gauge(rand.Float64())
}
