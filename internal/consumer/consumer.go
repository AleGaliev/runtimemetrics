package consumer

import (
	"context"
	"fmt"
	"time"

	models "github.com/AleGaliev/runtimemetrics/internal/model"
)

type Rep interface {
	SendMetricsRequest(metrics []models.Metrics) error
}

type MetricsConsumer struct {
	metrics        chan []models.Metrics
	workers        int
	Rep            Rep
	reportInterval int
}

func NewMetricsConsumer(workers, reportInterval int, rep Rep, metrics chan []models.Metrics) *MetricsConsumer {

	return &MetricsConsumer{
		metrics:        metrics,
		workers:        workers,
		reportInterval: reportInterval,
		Rep:            rep,
	}
}

func (mc *MetricsConsumer) ConsumerRun(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(mc.reportInterval) * time.Second)
	defer ticker.Stop()

	jobs := make(chan []models.Metrics, 100)
	for w := 1; w <= mc.workers; w++ {
		go mc.ConsumerWorker(jobs)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for metrics := range mc.metrics {
				jobs <- metrics
			}
		}
	}
}

func (mc *MetricsConsumer) ConsumerWorker(metrics <-chan []models.Metrics) {
	for metric := range metrics {
		err := mc.Rep.SendMetricsRequest(metric)
		if err != nil {
			fmt.Println(err)
		}
	}
}
