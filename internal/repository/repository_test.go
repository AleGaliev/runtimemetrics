package repository

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	models "github.com/AleGaliev/kubercontroller/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_createURLRequest(t *testing.T) {
	type args struct {
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
			assert.Equal(t, tt.want, createURLRequest(tt.args.name, tt.args.types, tt.args.value))
		})
	}
}

func TestSendMetrics_SendMetricsRequest1(t *testing.T) {
	allocValue := 123.45
	buckHashValue := 67.89
	freesValue := 0.12

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что метод POST
		assert.Equal(t, r.Method, "POST")
	}))
	defer server.Close()
	httpURL, _ := url.Parse(server.URL)
	BaseURL = httpURL.Host

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
			s := SendMetrics{
				Metrics: tt.fields.Metrics,
				Client:  tt.fields.Client,
			}
			s.SendMetricsRequest()
		})
	}
}
