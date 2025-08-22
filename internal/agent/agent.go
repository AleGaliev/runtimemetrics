package agent

import (
	"flag"
	"fmt"
	"os"
	"strconv"

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
	pollInterval   *int
	reportInterval *int
}

func NewAgentConfig(rep Rep) (*AgentConfig, error) {
	pollInterval := flag.Int("p", 2, "Interval poll metrics")
	reportInterval := flag.Int("r", 10, "Interval report metrics")
	varPollInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		StrPollInterval, err := strconv.Atoi(varPollInterval)
		if err != nil {
			return nil, fmt.Errorf("error converting POLL_INTERVAL to int: %v", err)
		}
		pollInterval = &StrPollInterval
	}
	varReportInterval, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		StrReportInterval, err := strconv.Atoi(varReportInterval)
		if err != nil {
			return nil, fmt.Errorf("error converting REPORT_INTERVAL to int: %v", err)
		}
		reportInterval = &StrReportInterval
	}
	flag.Parse()
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

	if c.counter%*c.pollInterval == 0 {
		metrics = collector.PullMetrics(c.pollCount)
		c.pollCount++
	}

	if c.counter%*c.reportInterval == 0 {
		if err := c.Rep.SendMetricsRequest(metrics); err != nil {
			return err
		}

	}
	c.counter++
	return nil

}
