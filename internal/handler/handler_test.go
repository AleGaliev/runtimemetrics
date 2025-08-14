package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/AleGaliev/kubercontroller/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestMyHandler_ServeHTTP(t *testing.T) {

	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name    string
		request string
		metod   string
		want    want
	}{
		{
			name:    "positive post request",
			request: "/update/counter/someMetric/527",
			metod:   http.MethodPost,
			want: want{
				statusCode:  200,
				contentType: "text/plain",
			},
		},
		{
			name:    "positive post gauge request",
			request: "/update/gauge/someMetric/527",
			metod:   http.MethodPost,
			want: want{
				statusCode:  200,
				contentType: "text/plain",
			},
		},
		{
			name:    "positive post gauge request",
			request: "/update/gauge/someMetric/527",
			metod:   http.MethodPost,
			want: want{
				statusCode:  200,
				contentType: "text/plain",
			},
		},
		{
			name:    "negativ get request",
			request: "/update/counter/someMetric/527",
			metod:   http.MethodGet,
			want: want{
				statusCode:  405,
				contentType: "text/plain",
			},
		},
		{
			name:    "negativ not name metrics",
			request: "/update/counter/527",
			metod:   http.MethodPost,
			want: want{
				statusCode:  404,
				contentType: "text/plain",
			},
		},
	}
	r := chi.NewRouter()
	h := MyHandler{}
	r.Post("/update/{type}/{name}/{value}", h.ServeHTTP)
	srv := httptest.NewServer(r)

	defer srv.Close()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = test.metod
			req.Header.Set("Content-Type", test.metod)
			req.URL = srv.URL + test.request

			res, err := req.Send()

			assert.Equal(t, test.want.statusCode, res.StatusCode())
			assert.Equal(t, nil, err)
			assert.Equal(t, test.metod, req.Header.Get("Content-Type"))

		})
	}
}

func TestMyHandler_GetValue(t *testing.T) {
	var (
		allocValue          = 123.45
		buckHashValue       = 67.89
		freesValue          = 0.12
		count         int64 = 3
	)

	tests := []struct {
		name    string
		request string
		status  int
		value   string
	}{
		{
			name:    "positive get allocValue",
			request: "/value/gauge/Alloc",
			status:  200,
			value:   fmt.Sprintf("%g", allocValue),
		},
		{
			name:    "positive get buckHashValue",
			request: "/value/gauge/BuckHash",
			status:  200,
			value:   fmt.Sprintf("%g", buckHashValue),
		},
		{
			name:    "positive get freesValue",
			request: "/value/gauge/Frees",
			status:  200,
			value:   fmt.Sprintf("%g", freesValue),
		},
		{
			name:    "positive get counter",
			request: "/value/counter/Count",
			status:  200,
			value:   fmt.Sprintf("%d", count),
		},
		{
			name:    "negativ get allocValue",
			request: "/value/gauge/allocValue",
			status:  404,
			value:   "",
		},
		// TODO: Add test cases.
	}
	models.MemStorage = map[string]models.Metrics{
		"Alloc": {
			ID:    "Alloc",
			MType: models.Gauge,
			Value: &allocValue,
		},
		"BuckHash": {
			ID:    "BuckHash",
			MType: models.Gauge,
			Value: &buckHashValue,
		},
		"Frees": {
			ID:    "Frees",
			MType: models.Gauge,
			Value: &freesValue,
		},
		"Count": {
			ID:    "Count",
			MType: models.Counter,
			Delta: &count,
		},
	}
	r := chi.NewRouter()
	h := MyHandler{}
	r.Get("/value/{type}/{name}", h.GetValue)
	srv := httptest.NewServer(r)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.URL = srv.URL + test.request

			res, err := req.Send()
			assert.Equal(t, test.status, res.StatusCode())
			assert.Equal(t, nil, err)

		})
	}
}

func TestMyHandler_ListMetrics(t *testing.T) {
	var (
		allocValue          = 123.45
		buckHashValue       = 67.89
		freesValue          = 0.12
		count         int64 = 3
	)

	tests := []struct {
		name    string
		request string
		status  int
	}{
		{
			name:    "positive test",
			request: "/",
			status:  200,
		},
	}
	models.MemStorage = map[string]models.Metrics{
		"Alloc": {
			ID:    "Alloc",
			MType: models.Gauge,
			Value: &allocValue,
		},
		"BuckHash": {
			ID:    "BuckHash",
			MType: models.Gauge,
			Value: &buckHashValue,
		},
		"Frees": {
			ID:    "Frees",
			MType: models.Gauge,
			Value: &freesValue,
		},
		"Count": {
			ID:    "Count",
			MType: models.Counter,
			Delta: &count,
		},
	}
	r := chi.NewRouter()
	h := MyHandler{}
	r.Get("/", h.ListMetrics)
	srv := httptest.NewServer(r)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.URL = srv.URL + test.request
			res, err := req.Send()
			defer res.Body()
			assert.Equal(t, test.status, res.StatusCode())
			assert.Equal(t, nil, err)
		})
	}
}
