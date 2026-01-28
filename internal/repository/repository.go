package repository

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	models "github.com/AleGaliev/runtimemetrics/internal/model"
	pb "github.com/AleGaliev/runtimemetrics/internal/proto"
	"github.com/AleGaliev/runtimemetrics/internal/proto/convertor"
	"github.com/AleGaliev/runtimemetrics/internal/service/hash"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type logger interface {
	CreateResponseLog(statusCode int, large int64)
}

type crypto interface {
	Encrypt(data []byte) ([]byte, error)
}

//generate:reset
type HTTPSendler struct {
	client     *http.Client
	grpcClient pb.MetricsClient
	url        *url.URL
	logger     logger
	crypto     crypto
	keyHash    string
	localIP    string
}

type Option func(*HTTPSendler)

func NewClientConfig(opts ...Option) *HTTPSendler {
	ip := getLocalIP()

	clientConfig := &HTTPSendler{
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
		url:     &url.URL{},
		logger:  nil,
		keyHash: "",
		crypto:  nil,
		localIP: ip,
	}

	for _, opt := range opts {
		opt(clientConfig)
	}

	return clientConfig
}

func WithGrpcClient(grpcHost string) Option {
	return func(h *HTTPSendler) {

		if grpcHost == "" {
			return
		}

		conn, err := grpc.NewClient(grpcHost,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any,
				cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

				md := metadata.Pairs("x-real-ip", h.localIP)
				ctx = metadata.NewOutgoingContext(ctx, md)

				return invoker(ctx, method, req, reply, cc, opts...)
			}),
		)

		if err != nil {
			os.Exit(1)
		}

		grpcClient := pb.NewMetricsClient(conn)

		h.grpcClient = grpcClient
	}
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
	request.Header.Set("X-Real-IP", h.localIP)

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

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func (h HTTPSendler) UpdateMetrics(metrics []models.Metrics) error {
	pbMetrics := convertor.ConvertMetricsSliceToProto(metrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := h.grpcClient.UpdateMetrics(ctx, pb.UpdateMetricsRequest_builder{
		Metrics: pbMetrics,
	}.Build())

	if err != nil {
		return fmt.Errorf("failed update metrics: %w", err)
	}
	return nil
}
