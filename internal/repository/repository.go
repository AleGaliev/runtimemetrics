package repository

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	models "github.com/AleGaliev/kubercontroller/internal/model"
)

type HttpSendler struct {
	client  *http.Client
	baseURL *string
	shema   string
}

func NewClientConfig() *HttpSendler {
	baseURL := flag.String("a", "localhost:8080", "Endpoint http server")
	return &HttpSendler{
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
		shema:   "http",
		baseURL: baseURL,
	}
}

func createURLRequest(shema, baseURL, name, types, value string) url.URL {
	fullPath := path.Join("update/", types, name, value)
	return url.URL{
		Scheme: shema,
		Host:   baseURL,
		Path:   fullPath,
	}
}

func (h HttpSendler) SendMetricsRequest(metrics []models.Metrics) error {
	for _, metric := range metrics {

		fullURL := url.URL{}
		if metric.MType == models.Gauge {
			fullURL = createURLRequest(h.shema, *h.baseURL, metric.ID, metric.MType, fmt.Sprint(*metric.Value))
		} else {
			fullURL = createURLRequest(h.shema, *h.baseURL, metric.ID, metric.MType, fmt.Sprint(*metric.Delta))
		}
		request, err := http.NewRequest(http.MethodPost, fullURL.String(), nil)
		if err != nil {
			return fmt.Errorf("Error creating request: %v", err)
		}
		request.Header.Set("Content-Type", "text/plain")
		response, err := h.client.Do(request)

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
