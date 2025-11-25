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
	contentTypeJSON  = "application/json"
	headerHashSHA256 = "HashSHA256"
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
		r.With(middleware.AuditMiddleware(eventAudit)).
			Post("/", h.ServeHTTPBatchUpdate)
	})

	muxWithMiddlewares.Get("/", h.ListMetrics)
	muxWithMiddlewares.Get("/ping", h.GetPing)

	return muxWithMiddlewares
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
