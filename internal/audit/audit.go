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

func CreateAuditMessage(metrics []models.Metrics, ip string) Audit {
	metricsNames := make([]string, len(metrics))
	for _, metric := range metrics {
		metricsNames = append(metricsNames, metric.ID)
	}
	return Audit{
		Ts:        time.Now().Unix(),
		Metrics:   metricsNames,
		IPAddress: ip,
	}
}
