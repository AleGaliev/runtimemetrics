package audit

import (
	"time"

	models "github.com/AleGaliev/runtimemetrics/internal/model"
)

type Audit struct {
	Ts        int64    `json:"ts"`
	Metrics   []string `json:"metrics"`
	IPAddress string   `json:"ip_address"`
}

func CreateAuditMessage(ip string, metrics []models.Metrics) Audit {
	metricsNames := make([]string, len(metrics))
	for i, metric := range metrics {
		metricsNames[i] = metric.ID
	}
	return Audit{
		Ts:        time.Now().Unix(),
		Metrics:   metricsNames,
		IPAddress: ip,
	}
}
