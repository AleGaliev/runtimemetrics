package storage

import (
	"fmt"
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
