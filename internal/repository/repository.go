package repository

import (
	"fmt"
	models "github.com/AleGaliev/kubercontroller/internal/model"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"strconv"
	"strings"
)

const (
	BaseUrl = "localhost:8080"
	Shema   = "http"
)

type SendMetrics struct {
	Url    map[string]url.URL
	Client *http.Client
}

func ConvertMemStatsInNameMetrics(memStats runtime.MemStats) map[string]string {
	return map[string]string{
		"alloc":         strconv.FormatUint(memStats.Alloc, 10),
		"BuckHashSys":   strconv.FormatUint(memStats.BuckHashSys, 10),
		"Frees":         strconv.FormatUint(memStats.Frees, 10),
		"GCCPUFraction": strconv.FormatFloat(memStats.GCCPUFraction, 'f', 10, 64),
		"GCSys":         strconv.FormatUint(memStats.GCSys, 10),
		"HeapAlloc":     strconv.FormatUint(memStats.HeapAlloc, 10),
		"HeapIdle":      strconv.FormatUint(memStats.HeapIdle, 10),
		"HeapInuse":     strconv.FormatUint(memStats.HeapInuse, 10),
		"HeapObjects":   strconv.FormatUint(memStats.HeapObjects, 10),
		"HeapReleased":  strconv.FormatUint(memStats.HeapReleased, 10),
		"HeapSys":       strconv.FormatUint(memStats.HeapSys, 10),
		"LastGC":        strconv.FormatUint(memStats.LastGC, 10),
		"Lookups":       strconv.FormatUint(memStats.Lookups, 10),
		"MCacheInuse":   strconv.FormatUint(memStats.MCacheInuse, 10),
		"MCacheSys":     strconv.FormatUint(memStats.MCacheSys, 10),
		"MSpanInuse":    strconv.FormatUint(memStats.MSpanInuse, 10),
		"MSpanSys":      strconv.FormatUint(memStats.MSpanSys, 10),
		"Mallocs":       strconv.FormatUint(memStats.Mallocs, 10),
		"NextGC":        strconv.FormatUint(memStats.NextGC, 10),
		"NumForcedGC":   strconv.FormatUint(uint64(memStats.NumForcedGC), 10),
		"NumGC":         strconv.FormatUint(uint64(memStats.NumGC), 10),
		"OtherSys":      strconv.FormatUint(memStats.OtherSys, 10),
		"PauseTotalNs":  strconv.FormatUint(memStats.PauseTotalNs, 10),
		"StackInuse":    strconv.FormatUint(memStats.StackInuse, 10),
		"StackSys":      strconv.FormatUint(memStats.StackSys, 10),
		"Sys":           strconv.FormatUint(memStats.Sys, 10),
		"TotalAlloc":    strconv.FormatUint(memStats.TotalAlloc, 10),
	}
}

func (s *SendMetrics) InitMetrics(mericsNameValue map[string]string) {
	for name, value := range mericsNameValue {
		fullUrl := createRequest(name, value)
		s.Url[name] = fullUrl
	}

}

func createRequest(name, metrics string) url.URL {
	fullPath := path.Join("update/", models.Gauge, strings.ToLower(name), metrics)
	return url.URL{
		Scheme: Shema,
		Host:   BaseUrl,
		Path:   fullPath,
	}
}

func (s *SendMetrics) DoRequest() {
	for _, url := range s.Url {
		request, err := http.NewRequest(http.MethodPost, url.String(), nil)
		if err != nil {
			fmt.Println(err)
		}
		response, err := s.Client.Do(request)
		response.Header.Set("Content-Type", "text/plain")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()
	}
}
