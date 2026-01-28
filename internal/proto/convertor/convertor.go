package convertor

import (
	models "github.com/AleGaliev/runtimemetrics/internal/model"
	pb "github.com/AleGaliev/runtimemetrics/internal/proto"
)

// ConvertProtoSliceToMetricsSimple converts []*pb.Metric to []model.Metrics
func ConvertProtoSliceToMetricsSimple(pbMetrics []*pb.Metric) []models.Metrics {
	result := make([]models.Metrics, 0, len(pbMetrics))

	for _, pbMetric := range pbMetrics {
		if pbMetric == nil {
			continue
		}

		metric := models.Metrics{
			ID: pbMetric.GetId(),
		}

		switch pbMetric.GetType() {
		case pb.Metric_GAUGE:
			metric.MType = models.Gauge
			value := pbMetric.GetValue()
			metric.Value = &value
		case pb.Metric_COUNTER:
			metric.MType = models.Counter
			delta := pbMetric.GetDelta()
			metric.Delta = &delta
		}

		result = append(result, metric)
	}

	return result
}

// convertMetricsToProto converts a single model.Metrics to pb.Metric
func ConvertMetricsToProto(m models.Metrics) *pb.Metric {
	metric := &pb.Metric_builder{
		Id: m.ID,
	}

	switch m.MType {
	case models.Gauge:
		if m.Value != nil {
			metric.Value = *m.Value
			metric.Type = pb.Metric_GAUGE
		}
	case models.Counter:
		if m.Delta != nil {
			metric.Delta = *m.Delta
			metric.Type = pb.Metric_COUNTER
		}
	}
	return metric.Build()
}

// convertMetricsSliceToProto converts []model.Metrics to []*pb.Metric
func ConvertMetricsSliceToProto(metrics []models.Metrics) []*pb.Metric {
	result := make([]*pb.Metric, 0, len(metrics))
	for _, m := range metrics {
		result = append(result, ConvertMetricsToProto(m))
	}
	return result
}
