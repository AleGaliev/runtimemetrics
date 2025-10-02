package consumer

import (
	"context"
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
	w := 3
	for i := 0; i < w; i++ {
		go func() {
			for {
				ticker := time.NewTicker(time.Duration(mc.reportInterval) * time.Second)
				defer ticker.Stop()
				select {
				case <-ticker.C:
					mc.Rep.SendMetricsRequest(<-mc.metrics)
				}
			}
		}()
	}
	select {
	case <-ctx.Done():
	}
}
