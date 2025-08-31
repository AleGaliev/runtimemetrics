package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	models "github.com/AleGaliev/kubercontroller/internal/model"
)

type Storage struct {
	Metrics map[string]models.Metrics
}

func CreateStorage() *Storage {
	return &Storage{
		Metrics: make(map[string]models.Metrics),
	}
}

func (s *Storage) AddMetric(myType, name, value string) error {
	switch myType {
	case models.Gauge:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		s.Metrics[name] = models.Metrics{
			ID:    name,
			MType: myType,
			Value: &f,
		}
	case models.Counter:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		if metrics, exists := s.Metrics[name]; exists {
			*metrics.Delta += i
		} else {
			s.Metrics[name] = models.Metrics{
				ID:    name,
				MType: myType,
				Delta: &i,
			}
		}
	default:
		return fmt.Errorf("unknown metric type: %s", myType)
	}
	return nil
}

func (s *Storage) GetMetrics(name string) (string, bool) {
	metric, ok := s.Metrics[name]
	if !ok {
		return "", false
	}
	switch metric.MType {
	case models.Gauge:
		return fmt.Sprintf("%g", *metric.Value), true
	case models.Counter:
		return fmt.Sprintf("%d", *metric.Delta), true
	}
	return "", false
}

func (s *Storage) GetAllMetric() string {
	result := ""
	for _, m := range s.Metrics {
		switch m.MType {
		case models.Gauge:
			result += fmt.Sprintf("<li> %s: %g</li>", m.ID, *m.Value)
		case models.Counter:
			result += fmt.Sprintf("<li> %s: %d</li>", m.ID, *m.Delta)

		}
	}
	return result
}

func (s *Storage) UpdateMetrics(r io.Reader) error {
	data := json.NewDecoder(r)
	var metricsData models.Metrics
	if err := data.Decode(&metricsData); err != nil {
		return fmt.Errorf("could not decode metrics: %v", err)
	}

	switch metricsData.MType {

	case models.Gauge:
		if metricsData.Value == nil {
			return fmt.Errorf("metrics value is nil")
		}
		s.Metrics[metricsData.ID] = metricsData

	case models.Counter:

		if metricsData.Delta == nil {
			return fmt.Errorf("metrics delta is nil")
		}

		if metric, exists := s.Metrics[metricsData.ID]; exists {
			*metric.Delta += *metricsData.Delta
		} else {
			s.Metrics[metricsData.ID] = metricsData
		}
	default:
		return fmt.Errorf("unknown metric type: %s", metricsData.MType)
	}

	return nil
}

func (s *Storage) ValueMetrics(r io.Reader) ([]byte, bool, error) {
	data := json.NewDecoder(r)
	var metrics models.Metrics
	if err := data.Decode(&metrics); err != nil {
		return nil, false, fmt.Errorf("could not decode metrics: %v", err)
	}
	if (metrics.MType != models.Counter && metrics.MType != models.Gauge) || metrics.ID == "" {
		return nil, false, fmt.Errorf("invalid metric type: %s", metrics.MType)
	}
	if metrics.Value != nil || metrics.Delta != nil {
		return nil, false, fmt.Errorf("invalid metric type: %s", metrics.MType)
	}

	metric, ok := s.Metrics[metrics.ID]
	if !ok {
		return nil, false, nil
	}
	if metric.MType != metrics.MType {
		return nil, false, fmt.Errorf("invalid metric type: %s", metrics.MType)
	}
	resp, err := json.MarshalIndent(metric, "", "  ")
	if err != nil {
		return nil, false, fmt.Errorf("could not encode metrics: %v", err)
	}
	return resp, true, nil
}
