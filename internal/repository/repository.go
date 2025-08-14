package repository

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	models "github.com/AleGaliev/kubercontroller/internal/model"
)

var (
	BaseURL = "localhost:8080"
	Shema   = "http"
)

type SendMetrics struct {
	Metrics []models.Metrics
	Client  *http.Client
}

func createURLRequest(name, types, value string) url.URL {
	fullPath := path.Join("update/", types, name, value)
	return url.URL{
		Scheme: Shema,
		Host:   BaseURL,
		Path:   fullPath,
	}
}

func (s SendMetrics) SendMetricsRequest() {
	for _, metric := range s.Metrics {

		fullURL := url.URL{}
		if metric.MType == models.Gauge {
			fullURL = createURLRequest(metric.ID, metric.MType, fmt.Sprint(*metric.Value))
		} else {
			fullURL = createURLRequest(metric.ID, metric.MType, fmt.Sprint(*metric.Delta))
		}

		request, err := http.NewRequest(http.MethodPost, fullURL.String(), nil)
		if err != nil {
			fmt.Println(err)
		}
		response, err := s.Client.Do(request)

		if err != nil {
			fmt.Errorf("FAILED to create request: %w", err)
		}

		defer response.Body.Close()
	}
}
