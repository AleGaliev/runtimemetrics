package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/AleGaliev/runtimemetrics/internal/middleware"
	"github.com/AleGaliev/runtimemetrics/internal/observer"
	"github.com/AleGaliev/runtimemetrics/internal/service/hash"
	"github.com/go-chi/chi/v5"
)

const (
	contentTypeJSON  = "application/json" // Content type for JSON responses
	headerHashSHA256 = "HashSHA256"       // Header name for SHA256 hash
)

// metricWriter defines the interface for writing metrics
type metricWriter interface {
	AddMetric(myType, name, value string) error
	UpdateMetrics(r io.Reader) error
	BatchUpdateMetrics(r io.Reader) error
}

// metricReader defines the interface for reading metrics
type metricReader interface {
	GetMetrics(name string) (string, bool)
	GetAllMetric() (string, error)
	ValueMetrics(r io.Reader) ([]byte, bool, error)
}

// connector defines the interface for database connection
type connector interface {
	Connect() error
}

// Storage combines both reading and writing capabilities for metrics
type Storage interface {
	metricWriter
	metricReader
}

// MyHandler handles HTTP requests for metrics operations
//
//generate:reset
type MyHandler struct {
	storage   Storage   // Storage for metrics data
	connector connector // Database connector
	hashKey   string    // Key for hash calculation
}

// CreateMyHandler creates and configures a new HTTP handler with routing and middleware
//
// Parameters:
//   - storage: metrics storage implementation
//   - connector: database connection implementation
//   - logger: middleware logger
//   - hashKey: key for hash validation
//   - eventAudit: event audit observer
//
// Returns configured HTTP handler
func CreateMyHandler(storage Storage, connector connector, logger middleware.Logger, hashKey string, eventAudit *observer.Event) http.Handler {
	h := &MyHandler{
		storage:   storage,
		connector: connector,
		hashKey:   hashKey,
	}

	mux := chi.NewRouter()
	muxWithMiddlewares := mux.With(
		middleware.MiddlewareHandlerLogger(logger),
		middleware.GzipMiddlewareHandler(),
		middleware.MetricValidateMiddleware(hashKey),
	)
	muxWithMiddlewares.Route("/", func(r chi.Router) {
		r.Mount("/debug/pprof", pprofRoutes())
	})
	muxWithMiddlewares.Route("/update/", func(r chi.Router) {
		r.Post("/", h.ServeHTTPUpdate)
		r.Post("/{type}/{name}/{value}", h.ServeHTTP)
	})

	muxWithMiddlewares.Route("/value/", func(r chi.Router) {
		r.Post("/", h.ServeHTTPValue)
		r.Get("/{type}/{name}", h.GetValue)
	})

	muxWithMiddlewares.Route("/updates/", func(r chi.Router) {
		r.With(middleware.AuditMiddleware(eventAudit, logger)).
			Post("/", h.ServeHTTPBatchUpdate)
	})

	muxWithMiddlewares.Get("/", h.ListMetrics)
	muxWithMiddlewares.Get("/ping", h.GetPing)

	return muxWithMiddlewares
}

// GetPing handles database connectivity check
// Returns 200 if connection successful, 500 otherwise
func (h MyHandler) GetPing(res http.ResponseWriter, _ *http.Request) {
	if err := h.connector.Connect(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	successResponse(res, h.hashKey)
}

// ServeHTTPUpdate adds a metric in JSON format
// Expects Content-Type: application/json header
// Returns 400 for bad requests, 200 for success
func (h MyHandler) ServeHTTPUpdate(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.storage.UpdateMetrics(req.Body); err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	successResponse(res, h.hashKey)
}

// ServeHTTPBatchUpdate performs batch update of metrics in JSON format
// Expects POST method and Content-Type: application/json header
// Returns 400 for bad requests, 200 for success
func (h MyHandler) ServeHTTPBatchUpdate(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost || req.Header.Get("Content-Type") != "application/json" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.storage.BatchUpdateMetrics(req.Body); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	successResponse(res, h.hashKey)
}

// ServeHTTPValue retrieves metric value in JSON format
// Expects Content-Type: application/json header
// Returns 400 for bad requests, 404 if metric not found, 200 with metric data for success
func (h MyHandler) ServeHTTPValue(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	metrics, ok, err := h.storage.ValueMetrics(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if h.hashKey != "" {
		hashSHA256 := hash.CreateHash(h.hashKey, metrics)
		res.Header().Set("HashSHA256", hashSHA256)
	}

	_, err = res.Write(metrics)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// ServeHTTP adds a metric via URL parameters
// URL format: /update/{type}/{name}/{value}
// Returns 404 for missing parameters, 400 for invalid data
func (h MyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	myType := chi.URLParam(req, "type")
	value := chi.URLParam(req, "value")

	if name == "" || myType == "" || value == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	err := h.storage.AddMetric(myType, name, value)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

// GetValue retrieves a metric value by name
// URL format: /value/{type}/{name}
// Returns 404 if metric not found, 200 with metric value for success
func (h MyHandler) GetValue(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")

	metricName := chi.URLParam(req, "name")

	metric, ok := h.storage.GetMetrics(metricName)
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "%s", metric)
}

// ListMetrics returns HTML page with all metrics
// Returns 500 for internal errors, 200 with HTML content for success
func (h MyHandler) ListMetrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")

	body, err := h.storage.GetAllMetric()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Используем strings.Builder для эффективной конкатенации
	var html strings.Builder
	html.WriteString(`
    <!DOCTYPE html>
    <html>
    <body>
        <h1>Metrics List</h1>
        <ul>
    `)
	html.WriteString(body)
	html.WriteString(`
        </ul>
    </body>
    </html>
    `)

	fmt.Fprint(res, html.String())
}

// successResponse sends a standardized success JSON response
// Includes hash header if hashKey is provided
func successResponse(res http.ResponseWriter, hashKey string) {
	response := map[string]interface{}{
		"status":  "success",
		"message": "Запрос обработан",
	}

	encoder := json.NewEncoder(res)
	if hashKey != "" {
		var buf bytes.Buffer
		if err := encoder.Encode(response); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		hashSHA256 := hash.CreateHash(hashKey, buf.Bytes())
		res.Header().Set(headerHashSHA256, hashSHA256)
		buf.WriteTo(res)
	} else {
		encoder.Encode(response)
	}
}

// pprofRoutes configures pprof endpoints for debugging
// Returns handler with all pprof routes
func pprofRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", pprof.Index)
	r.Get("/cmdline", pprof.Cmdline)
	r.Get("/profile", pprof.Profile)
	r.Get("/symbol", pprof.Symbol)
	r.Get("/trace", pprof.Trace)
	r.Handle("/goroutine", pprof.Handler("goroutine"))
	r.Handle("/heap", pprof.Handler("heap"))
	r.Handle("/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/block", pprof.Handler("block"))
	r.Handle("/mutex", pprof.Handler("mutex"))
	r.Handle("/allocs", pprof.Handler("allocs"))
	return r
}
