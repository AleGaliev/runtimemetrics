package repository

import (
	agentlogger "github.com/AleGaliev/kubercontroller/internal/logger"
	models "github.com/AleGaliev/kubercontroller/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_createURLRequest(t *testing.T) {
	type args struct {
		shema string
		addr  string
		name  string
		types string
		value string
	}
	tests := []struct {
		name string
		args args
		want url.URL
	}{
		{
			name: "positive",
			args: args{
				shema: "http",
				addr:  "localhost:8080",
				name:  "cpu",
				types: "gauge",
				value: "1234",
			},
			want: url.URL{
				Scheme: "http",
				Host:   "localhost:8080",
				Path:   "update/gauge/cpu/1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, createURLRequest(tt.args.shema, tt.args.addr, tt.args.name, tt.args.types, tt.args.value))
		})
	}
}

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
		client:  server.Client(),
		baseURL: &httpURL.Host,
		shema:   "http",
		logger:  loggerAgent,
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
