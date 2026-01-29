package agent

import (
	"context"

	"github.com/AleGaliev/runtimemetrics/internal/collector"
	"github.com/AleGaliev/runtimemetrics/internal/consumer"
	models "github.com/AleGaliev/runtimemetrics/internal/model"
	"github.com/AleGaliev/runtimemetrics/internal/repository"
	"github.com/AleGaliev/runtimemetrics/internal/service/retry"
)

type Rep interface {
	SendMetricsRequest(metrics []models.Metrics) error
}

type Agent struct {
	Rep            repository.HTTPSender
	BaseURL        *string
	grpcHost       string
	retry          retry.Retry
	pollCount      int
	counter        int
	pollInterval   int
	reportInterval int
	workers        int
}

type Option func(*Agent)

func New(rep repository.HTTPSender, opts ...Option) *Agent {
	agent := &Agent{
		Rep:       rep,
		pollCount: 1,
		counter:   1,
	}

	for _, opt := range opts {
		opt(agent)
	}

	return agent
}

func WithBaseURL(baseURL string) Option {
	return func(agent *Agent) {
		agent.BaseURL = &baseURL
	}
}

func WithGRPCHost(grpcHost string) Option {
	return func(agent *Agent) {
		agent.grpcHost = grpcHost
	}
}

func WithRetry(retry retry.Retry) Option {
	return func(agent *Agent) {
		agent.retry = retry
	}
}

func WithPollCount(pollCount int) Option {
	return func(agent *Agent) {
		agent.pollCount = pollCount
	}
}

func WithCounter(counter int) Option {
	return func(agent *Agent) {
		agent.counter = counter
	}
}

func WithPollInterval(pollInterval int) Option {
	return func(agent *Agent) {
		agent.pollInterval = pollInterval
	}
}

func WithReportInterval(reportInterval int) Option {
	return func(agent *Agent) {
		agent.reportInterval = reportInterval
	}
}

func WithWorkers(workers int) Option {
	return func(agent *Agent) {
		agent.workers = workers
	}
}

//func NewAgentConfig(rep repository.HTTPSendler, retry retry.Retry, pollInterval, reportInterval, workers int) (*AgentConfig, error) {
//	return &AgentConfig{
//		pollCount:      1,
//		counter:        1,
//		pollInterval:   pollInterval,
//		reportInterval: reportInterval,
//		Rep:            rep,
//		retry:          retry,
//		workers:        workers,
//	}, nil
//}

func (c *Agent) Run(ctx context.Context) error {
	metrics := make(chan []models.Metrics, 100)

	pullMetrics := collector.NewMetricsCollector(c.pollInterval, metrics)
	go pullMetrics.CollectMetrics(ctx)

	reportMetrics := consumer.New(metrics,
		consumer.WithGrpc(c.grpcHost),
		consumer.WithRep(c.Rep),
		consumer.WithReportInterval(c.reportInterval),
		consumer.WithRetry(&c.retry),
		consumer.WithWorkers(c.workers),
	)

	if err := reportMetrics.ConsumerRun(ctx); err != nil {
		return err
	}

	return nil
}
