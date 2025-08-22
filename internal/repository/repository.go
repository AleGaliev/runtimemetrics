package repository

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	models "github.com/AleGaliev/kubercontroller/internal/model"
)

type logger interface {
	CreateResponseLog(statusCode int, large int64)
}
type HTTPSendler struct {
	client *http.Client
	//baseURL *string
	//shema   string
	url    *url.URL
	logger logger
}

func NewClientConfig(logger logger) *HTTPSendler {
	baseURL := flag.String("a", "localhost:8080", "Endpoint http server")
	varAdrHost, ok := os.LookupEnv("ADDRESS")
	if ok {
		baseURL = &varAdrHost
	}
	flag.Parse()
	return &HTTPSendler{
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
		url: &url.URL{
			Scheme: "http",
			Host:   *baseURL,
			Path:   "update",
		},
		logger: logger,
	}
}

func (h HTTPSendler) SendMetricsRequest(metrics []models.Metrics) error {
	for _, metric := range metrics {

		jsonMetrics, err := json.Marshal(metric)
		if err != nil {
			return fmt.Errorf("could not marshal metrics: %v", err)
		}

		request, err := http.NewRequest(http.MethodPost, h.url.String(), bytes.NewBuffer(jsonMetrics))
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}
		request.Header.Set("Content-Type", "application/json")

		response, err := h.MiddlewareLoggerDo(request)

		if err != nil {
			return fmt.Errorf("error sending request: %v", err)
		}
		response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("error sending request: %d %s", response.StatusCode, response.Status)
		}
	}

	return nil
}

func (h HTTPSendler) MiddlewareLoggerDo(req *http.Request) (*http.Response, error) {
	response, err := h.client.Do(req)
	h.logger.CreateResponseLog(response.StatusCode, response.ContentLength)
	if err != nil {
		return nil, err
	}
	return response, nil
}
