package repository

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	agentlogger "github.com/AleGaliev/kubercontroller/internal/logger"
	models "github.com/AleGaliev/kubercontroller/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestSendMetrics_SendMetricsRequest(t *testing.T) {
	allocValue := 123.45
	buckHashValue := 67.89
	freesValue := 0.12

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что метод POST
		assert.Equal(t, r.Method, "POST")
	}))
	defer server.Close()
	httpURL, _ := url.Parse(server.URL)
	loggerAgent, err := agentlogger.CreateLogger()
	if err != nil {
		t.Fatal(err)
	}
	clientCfg := HTTPSendler{
		client: server.Client(),
		url:    url.URL{Scheme: httpURL.Scheme, Host: httpURL.Host, Path: httpURL.Path},
		logger: loggerAgent,
	}

	type fields struct {
		Metrics []models.Metrics
		Client  *http.Client
	}
	tests := []struct {
		name   string
		fields fields
	}{{
		name: "positive",
		fields: fields{
			Metrics: []models.Metrics{
				{ID: "Alloc", MType: models.Gauge, Value: &allocValue},
				{ID: "BuckHashSys", MType: models.Gauge, Value: &buckHashValue},
				{ID: "Frees", MType: models.Gauge, Value: &freesValue},
			},
			Client: server.Client(),
		},
	},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientCfg.SendMetricsRequest(tt.fields.Metrics)

		})
	}
}
