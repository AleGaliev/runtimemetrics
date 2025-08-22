package agent

import (
	"github.com/AleGaliev/kubercontroller/internal/collector"
	models "github.com/AleGaliev/kubercontroller/internal/model"
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
}

func NewAgentConfig(rep Rep, pollInterval, reportInterval int) (*AgentConfig, error) {
	return &AgentConfig{
		pollCount:      1,
		counter:        1,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		Rep:            rep,
	}, nil
}

func (c *AgentConfig) Run() error {
	metrics := []models.Metrics{}

	if c.counter%c.pollInterval == 0 {
		metrics = collector.PullMetrics(c.pollCount)
		c.pollCount++
	}

	if c.counter%c.reportInterval == 0 {
		if err := c.Rep.SendMetricsRequest(metrics); err != nil {
			return err
		}

	}
	c.counter++
	return nil

}
