package agent

import (
	"context"

	"github.com/AleGaliev/runtimemetrics/internal/collector"
	"github.com/AleGaliev/runtimemetrics/internal/consumer"
	models "github.com/AleGaliev/runtimemetrics/internal/model"
	"github.com/AleGaliev/runtimemetrics/internal/service/retry"
)

type Rep interface {
	SendMetricsRequest(metrics []models.Metrics) error
}

type AgentConfig struct {
	Rep            Rep
	BaseURL        *string
	pollCount      int64
	counter        int
	pollInterval   int
	reportInterval int
	retry          retry.Retry
	workers        int
}

func NewAgentConfig(rep Rep, retry retry.Retry, pollInterval, reportInterval, workers int) (*AgentConfig, error) {
	return &AgentConfig{
		pollCount:      1,
		counter:        1,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		Rep:            rep,
		retry:          retry,
		workers:        workers,
	}, nil
}

func (c *AgentConfig) Run(ctx context.Context) error {
	metrics := make(chan []models.Metrics, 100)

	pullMetrics := collector.NewMetricsCollector(c.pollInterval, metrics)
	go pullMetrics.CollectMetrics(ctx)

	reportMetrics := consumer.NewMetricsConsumer(c.workers, c.reportInterval, c.Rep, &c.retry, metrics)

	if err := reportMetrics.ConsumerRun(ctx); err != nil {
		return err
	}

	return nil
}
