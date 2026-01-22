package repository

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	models "github.com/AleGaliev/runtimemetrics/internal/model"
	"github.com/AleGaliev/runtimemetrics/internal/service/hash"
)

type logger interface {
	CreateResponseLog(statusCode int, large int64)
}

type crypto interface {
	Encrypt(data []byte) ([]byte, error)
}

//generate:reset
type HTTPSendler struct {
	client  *http.Client
	url     *url.URL
	logger  logger
	crypto  crypto
	keyHash string
}

type Option func(*HTTPSendler)

func NewClientConfig(opts ...Option) *HTTPSendler {
	clientConfig := &HTTPSendler{
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
		url:     &url.URL{},
		logger:  nil,
		keyHash: "",
		crypto:  nil,
	}

	for _, opt := range opts {
		opt(clientConfig)
	}

	return clientConfig
}

func WithCrypto(crypto crypto) Option {
	return func(h *HTTPSendler) {
		h.crypto = crypto
	}
}

func WithURL(baseURL string) Option {
	return func(h *HTTPSendler) {
		h.url = &url.URL{
			Scheme: "http",
			Host:   baseURL,
			Path:   "updates/",
		}
	}
}

func WithKeyHash(keyHash string) Option {
	return func(h *HTTPSendler) {
		h.keyHash = keyHash
	}
}

func WithLogger(logger logger) Option {
	return func(h *HTTPSendler) {
		h.logger = logger
	}
}

func (h HTTPSendler) SendMetricsRequest(metrics []models.Metrics) error {
	jsonMetrics, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("could not marshal metrics: %v", err)
	}
	jsonMetrics, err = h.crypto.Encrypt(jsonMetrics)
	if err != nil {
		return fmt.Errorf("could not encrypt metrics: %v", err)
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err = gz.Write(jsonMetrics); err != nil {
		return fmt.Errorf("could not gzip metrics: %v", err)
	}
	if err = gz.Close(); err != nil {
		return fmt.Errorf("could not gzip metrics: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, h.url.String(), &buf)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	if h.keyHash != "" {
		hashSendler := hash.CreateHash(h.keyHash, jsonMetrics)
		request.Header.Set(`HashSHA256`, hashSendler)
	}

	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept-Encoding", "gzip")

	response, err := h.MiddlewareLoggerDo(request)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error sending request: %d %s", response.StatusCode, response.Status)
	}

	return nil
}

func (h HTTPSendler) MiddlewareLoggerDo(req *http.Request) (*http.Response, error) {
	response, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	h.logger.CreateResponseLog(response.StatusCode, response.ContentLength)
	return response, nil
}
