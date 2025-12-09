package audit

import (
	"time"

	models "github.com/AleGaliev/runtimemetrics/internal/model"
)

type Audit struct {
	IPAddress string   `json:"ip_address"`
	Metrics   []string `json:"metrics"`
	TS        int64    `json:"ts"`
}

func CreateAuditMessage(ip string, metrics []models.Metrics) Audit {
	metricsNames := make([]string, len(metrics))
	for i, metric := range metrics {
		metricsNames[i] = metric.ID
	}
	return Audit{
		TS:        time.Now().Unix(),
		Metrics:   metricsNames,
		IPAddress: ip,
	}
}
