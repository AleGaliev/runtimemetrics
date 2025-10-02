package collector

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	models "github.com/AleGaliev/runtimemetrics/internal/model"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type MetricsCollector struct {
	mu           sync.RWMutex
	metrics      chan []models.Metrics
	pollCount    int64
	pollInterval int
}

func NewMetricsCollector(pollInterval int, metrics chan []models.Metrics) *MetricsCollector {
	return &MetricsCollector{
		metrics:      metrics,
		pollCount:    0,
		pollInterval: pollInterval,
	}
}

func float64Ptr(f float64) *float64 {
	return &f
}

func (mc *MetricsCollector) CollectMetrics(ctx context.Context) {

	tickerSystem := time.NewTicker(time.Duration(mc.pollInterval) * time.Second)
	defer tickerSystem.Stop()
	tickerRuntime := time.NewTicker(time.Duration(mc.pollInterval) * time.Second)
	defer tickerRuntime.Stop()

	go func() {
		var pullCounter int64
		for {
			pullCounter++

			mc.metrics <- collectRuntimeMetrics(pullCounter)

			select {
			case <-ctx.Done():
				return
			case <-tickerRuntime.C:
			}
		}
	}()

	go func() {
		for {
			metrics, _ := collectSystemMetrics()
			mc.metrics <- metrics
			select {
			case <-ctx.Done():
				return
			case <-tickerSystem.C:
			}
		}
	}()

	<-ctx.Done()

}

func collectSystemMetrics() ([]models.Metrics, error) {
	sysMetrics, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("collectSystemMetrics error: %w", err)
	}
	logicalCPUs, err := cpu.Counts(true)
	if err != nil {
		return nil, fmt.Errorf("collectSystemMetrics cpu error: %w", err)
	}
	metrics := []models.Metrics{
		{ID: "TotalMemory", MType: models.Gauge, Value: float64Ptr(float64(sysMetrics.Total))},
		{ID: "FreeMemory", MType: models.Gauge, Value: float64Ptr(float64(sysMetrics.Free))},
		{ID: "CPUutilization1", MType: models.Gauge, Value: float64Ptr(float64(logicalCPUs))},
	}
	return metrics, nil

}

func collectRuntimeMetrics(pollCount int64) []models.Metrics {
	metRuntime := runtime.MemStats{}

	runtime.ReadMemStats(&metRuntime)

	metrics := []models.Metrics{
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
	return metrics
}
