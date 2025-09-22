package agent

import (
	"github.com/AleGaliev/runtimemetrics/internal/collector"
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
}

func NewAgentConfig(rep Rep, retry retry.Retry, pollInterval, reportInterval int) (*AgentConfig, error) {
	return &AgentConfig{
		pollCount:      1,
		counter:        1,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		Rep:            rep,
		retry:          retry,
	}, nil
}

func (c *AgentConfig) Run() error {
	metrics := []models.Metrics{}

	if c.counter%c.pollInterval == 0 {
		metrics = collector.PullMetrics(c.pollCount)
		c.pollCount++
	}

	if c.counter%c.reportInterval == 0 {
		if err := c.retry.RetryConnection(func() error {
			if err := c.Rep.SendMetricsRequest(metrics); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}
	c.counter++
	return nil

}
