package consumer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	models "github.com/AleGaliev/runtimemetrics/internal/model"
)

type Rep interface {
	SendMetricsRequest(metrics []models.Metrics) error
	UpdateMetrics(metrics []models.Metrics) error
}

type Retry interface {
	RetryConnection(dbFunction func() error) error
}

type Consumer struct {
	rep            Rep
	Retry          Retry
	metrics        chan []models.Metrics
	workers        int
	reportInterval int
	grpc           bool
}

type Option func(*Consumer)

func New(metrics chan []models.Metrics, opts ...Option) *Consumer {
	clientConfig := &Consumer{
		metrics:        metrics,
		workers:        1,
		reportInterval: 10,
		grpc:           false,
	}

	for _, opt := range opts {
		opt(clientConfig)
	}

	return clientConfig
}

func WithRep(rep Rep) Option {
	return func(consumer *Consumer) {
		consumer.rep = rep
	}
}

func WithRetry(retry Retry) Option {
	return func(consumer *Consumer) {
		consumer.Retry = retry
	}
}

func WithWorkers(workers int) Option {
	return func(consumer *Consumer) {
		consumer.workers = workers
	}
}

func WithReportInterval(reportInterval int) Option {
	return func(consumer *Consumer) {
		consumer.reportInterval = reportInterval
	}
}

func WithGrpc(grpcAdr string) Option {
	return func(consumer *Consumer) {
		if grpcAdr != "" {
			consumer.grpc = true
		}
	}
}

func (c *Consumer) ConsumerRun(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(c.reportInterval) * time.Second)
	defer ticker.Stop()
	errCh := make(chan error, 1)
	jobs := make(chan []models.Metrics, 100)
	wg := &sync.WaitGroup{}

	for w := 1; w <= c.workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := c.ConsumerWorker(jobs); err != nil {
				errCh <- fmt.Errorf("consumer worker failed: %w", err)
				log.Printf("Consumer worker failed: %v", err)
			}
		}()
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				for metrics := range c.metrics {
					jobs <- metrics
				}
				close(jobs)
				return
			case <-ticker.C:
				metrics := collectInitialMetrics(c.metrics)
				for _, m := range metrics {
					jobs <- m
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		wg.Wait()
		return nil
	case err := <-errCh:
		return err
	}
}

func (c *Consumer) ConsumerWorker(metrics <-chan []models.Metrics) error {
	for metric := range metrics {
		if c.grpc {
			if err := c.rep.UpdateMetrics(metric); err != nil {
				return err
			}
			continue
		}
		if err := c.Retry.RetryConnection(func() error {
			if err := c.rep.SendMetricsRequest(metric); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func collectInitialMetrics(allMetrics <-chan []models.Metrics) [][]models.Metrics {
	var initialMetrics [][]models.Metrics
	for {
		select {
		case metrics := <-allMetrics:
			initialMetrics = append(initialMetrics, metrics)
		default:
			return initialMetrics
		}
	}
}
