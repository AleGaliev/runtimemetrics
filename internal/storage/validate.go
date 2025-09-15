package storage

import (
	"fmt"

	models "github.com/AleGaliev/kubercontroller/internal/model"
)

func MetricValidate(metric models.Metrics) error {

	if metric.ID == "" || (metric.Value == nil && metric.Delta == nil) {
		return fmt.Errorf("metric id or value is required")
	}
	switch metric.MType {

	case models.Gauge:
		if metric.Value == nil {
			return fmt.Errorf("metrics value is nil")
		}

	case models.Counter:
		if metric.Delta == nil {
			return fmt.Errorf("metrics delta is nil")
		}

	default:
		return fmt.Errorf("invalid metric type: %s", metric.MType)
	}

	return nil
}
