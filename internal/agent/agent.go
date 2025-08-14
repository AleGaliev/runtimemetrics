package agent

import (
	"math/rand"
	"net/http"
	"runtime"
	"time"

	models "github.com/AleGaliev/kubercontroller/internal/model"
	"github.com/AleGaliev/kubercontroller/internal/repository"
)

var (
	metRuntime  = runtime.MemStats{}
	sendMetrics = repository.SendMetrics{
		Metrics: []models.Metrics{},
		Client: &http.Client{
			Timeout: time.Duration(2 * time.Second),
		},
	}
	pollCount      int64 = 1
	pollInterval         = 2
	reportInterval       = 10
)

func float64Ptr(f float64) *float64 {
	return &f
}

func Run(counter int) {
	if counter%pollInterval == 0 {
		runtime.ReadMemStats(&metRuntime)

		sendMetrics.Metrics = []models.Metrics{
			{ID: "Alloc", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.Alloc))},
			{ID: "BuckHashSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.BuckHashSys))},
			{ID: "Frees", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.Frees))},
			{ID: "GCCPUFraction", MType: models.Gauge, Value: float64Ptr(metRuntime.GCCPUFraction)},
			{ID: "GCSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.GCSys))},
			{ID: "HeapAlloc", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.HeapAlloc))},
			{ID: "HeapIdle", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.HeapIdle))},
			{ID: "HeapInuse", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.HeapInuse))},
			{ID: "HeapObjects", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.HeapObjects))},
			{ID: "HeapReleased", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.HeapReleased))},
			{ID: "HeapSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.HeapSys))},
			{ID: "LastGC", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.LastGC))},
			{ID: "Lookups", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.Lookups))},
			{ID: "MCacheInuse", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.MCacheInuse))},
			{ID: "MCacheSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.MCacheSys))},
			{ID: "MSpanInuse", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.MSpanInuse))},
			{ID: "MSpanSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.MSpanSys))},
			{ID: "Mallocs", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.Mallocs))},
			{ID: "NextGC", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.NextGC))},
			{ID: "NumForcedGC", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.NumForcedGC))},
			{ID: "NumGC", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.NumGC))},
			{ID: "OtherSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.OtherSys))},
			{ID: "PauseTotalNs", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.PauseTotalNs))},
			{ID: "StackInuse", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.StackInuse))},
			{ID: "StackSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.StackSys))},
			{ID: "Sys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.Sys))},
			{ID: "TotalAlloc", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.TotalAlloc))},
			{ID: "RandomValue", MType: models.Gauge, Value: float64Ptr(rand.Float64())},
			{ID: "PollCount", MType: models.Counter, Delta: &pollCount},
		}
		pollCount++
	}

	if counter%reportInterval == 0 {
		sendMetrics.SendMetricsRequest()
	}

}
