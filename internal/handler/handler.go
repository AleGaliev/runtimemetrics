package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type logger interface {
	CreateRequestLog(url, method string, timestamp time.Time)
}

type Storage interface {
	AddMetric(myType, name, value string) error
	GetMetrics(name string) (string, bool)
	GetAllMetric() string
	UpdateMetrics(r io.Reader) error
	ValueMetrics(r io.Reader) ([]byte, error)
}
type MyHandler struct {
	Storage Storage
	logger  logger
}

func CreateMyHandler(storage Storage, logger logger) http.Handler {
	h := &MyHandler{
		Storage: storage,
		logger:  logger,
	}

	mux := chi.NewRouter()

	mux.Post("/update", h.ServeHTTPUpdate)
	mux.Post("/update/{type}/{name}/{value}", h.ServeHTTP)
	mux.Post("/value", h.ServeHTTPValue)
	mux.Get("/value/{type}/{name}", h.GetValue)
	mux.Get("/", h.ListMetrics)

	muxMiddlewareLogger := h.MiddlewareHandlerLogger(mux)

	return muxMiddlewareLogger
}

func (h MyHandler) ServeHTTPUpdate(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	if err := h.Storage.UpdateMetrics(req.Body); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h MyHandler) ServeHTTPValue(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	metrics, err := h.Storage.ValueMetrics(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(res, "%s", metrics)
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

	err := h.Storage.AddMetric(myType, name, value)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h MyHandler) GetValue(res http.ResponseWriter, req *http.Request) {
	metricName := chi.URLParam(req, "name")

	metric, ok := h.Storage.GetMetrics(metricName)
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "%s", metric)
}

func (h MyHandler) ListMetrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")

	body := h.Storage.GetAllMetric()

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

func (h MyHandler) MiddlewareHandlerLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(res, req)
		h.logger.CreateRequestLog(req.RequestURI, req.Method, start)
	})
}
