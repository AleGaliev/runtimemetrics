package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
			request: "/update/gauge/someMetric/527/",
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
	h := MyHandler{}
	handler := http.HandlerFunc(h.ServeHTTP)
	srv := httptest.NewServer(handler)

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
