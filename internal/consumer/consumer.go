package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	models "github.com/AleGaliev/runtimemetrics/internal/model"
)

type Rep interface {
	SendMetricsRequest(metrics []models.Metrics) error
}

type Retry interface {
	RetryConnection(dbFunction func() error) error
}

type MetricsConsumer struct {
	rep            Rep
	Retry          Retry
	metrics        chan []models.Metrics
	workers        int
	reportInterval int
}

func NewMetricsConsumer(workers, reportInterval int, rep Rep, retry Retry, metrics chan []models.Metrics) *MetricsConsumer {
	return &MetricsConsumer{
		metrics:        metrics,
		workers:        workers,
		reportInterval: reportInterval,
		rep:            rep,
		Retry:          retry,
	}
}

func (mc *MetricsConsumer) ConsumerRun(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(mc.reportInterval) * time.Second)
	defer ticker.Stop()
	errCh := make(chan error, 1)
	jobs := make(chan []models.Metrics, 100)
	for w := 1; w <= mc.workers; w++ {
		go func() {
			if err := mc.ConsumerWorker(jobs); err != nil {
				errCh <- fmt.Errorf("consumer worker failed: %w", err)
				log.Printf("Consumer worker failed: %v", err)
			}
		}()
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				metrics := collectInitialMetrics(mc.metrics)
				for _, m := range metrics {
					jobs <- m
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func (mc *MetricsConsumer) ConsumerWorker(metrics <-chan []models.Metrics) error {
	for metric := range metrics {
		if err := mc.Retry.RetryConnection(func() error {
			if err := mc.rep.SendMetricsRequest(metric); err != nil {
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
