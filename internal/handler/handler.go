package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/AleGaliev/runtimemetrics/internal/middleware"
	"github.com/AleGaliev/runtimemetrics/internal/service/hash"
	"github.com/go-chi/chi/v5"
)

type metricWriter interface {
	AddMetric(myType, name, value string) error
	UpdateMetrics(r io.Reader) error
	BatchUpdateMetrics(r io.Reader) error
}

type metricReader interface {
	GetMetrics(name string) (string, bool)
	GetAllMetric() (string, error)
	ValueMetrics(r io.Reader) ([]byte, bool, error)
}

type connector interface {
	Connect() error
}

type Storage interface {
	metricWriter
	metricReader
}
type MyHandler struct {
	storage   Storage
	connector connector
	hashKey   string
}

func CreateMyHandler(storage Storage, connector connector, logger middleware.Logger, hashKey string) http.Handler {
	h := &MyHandler{
		storage:   storage,
		connector: connector,
		hashKey:   hashKey,
	}

	mux := chi.NewRouter()

	mux.Route("/update/", func(r chi.Router) {
		r.Post("/", h.ServeHTTPUpdate)
		r.Post("/{type}/{name}/{value}", h.ServeHTTP)
	})

	mux.Route("/value/", func(r chi.Router) {
		r.Post("/", h.ServeHTTPValue)
		r.Get("/{type}/{name}", h.GetValue)
	})

	mux.Route("/updates/", func(r chi.Router) {
		r.Post("/", h.ServeHTTPBatchUpdate)
	})

	mux.Get("/", h.ListMetrics)
	mux.Get("/ping", h.GetPing)

	muxMiddlewareValidate := middleware.MetricValidateMiddleware(mux, hashKey)
	muxGzip := middleware.GzipMiddlewareHandler(muxMiddlewareValidate)
	muxMiddlewareLogger := middleware.MiddlewareHandlerLogger(muxGzip, logger)

	return muxMiddlewareLogger
}

func (h MyHandler) GetPing(res http.ResponseWriter, _ *http.Request) {
	if err := h.connector.Connect(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	successResponse(res, h.hashKey)
}

// ServeHTTPUpdate добавление метрики в формате json
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

// ServeHTTPValue получение метрик в формате json
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

func (h MyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	pathURL := strings.Split(strings.Trim(req.URL.Path, "/"), "/")

	if len(pathURL) < 4 {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, "404 page not found")
		return
	}
	name := chi.URLParam(req, "name")
	myType := chi.URLParam(req, "type")
	value := chi.URLParam(req, "value")

	err := h.storage.AddMetric(myType, name, value)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

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

func (h MyHandler) ListMetrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	body, err := h.storage.GetAllMetric()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprint(res, `
    <!DOCTYPE html>
    <html>
    <body>
        <h1>Metrics List</h1>
		<ul>
    `)
	fmt.Fprintf(res, `%s`, body)

	fmt.Fprint(res, `
	</ul>
    </body>
    </html>
    `)
}

func successResponse(res http.ResponseWriter, hashKey string) {
	response := map[string]interface{}{
		"status":  `success`,
		"message": "Запрос обработан",
	}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
	if hashKey != "" {
		hashSHA256 := hash.CreateHash(hashKey, jsonBytes)
		res.Header().Set("HashSHA256", hashSHA256)
	}
	_, err = res.Write(jsonBytes)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}
