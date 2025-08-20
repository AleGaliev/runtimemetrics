package handler

import (
	"fmt"
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

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", h.ServeHTTP)
	r.Get("/value/{type}/{name}", h.GetValue)
	r.Get("/", h.ListMetrics)

	return r
}

func (h MyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
	pathURL := strings.Split(strings.Trim(req.URL.Path, "/"), "/")

	if len(pathURL) < 4 {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, "404 page not found")
		h.logger.CreateRequestLog(req.RequestURI, req.Method, start)
		return
	}
	name := chi.URLParam(req, "name")
	myType := chi.URLParam(req, "type")
	value := chi.URLParam(req, "value")

	err := h.Storage.AddMetric(myType, name, value)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		h.logger.CreateRequestLog(req.RequestURI, req.Method, start)
		return
	}
	h.logger.CreateRequestLog(req.RequestURI, req.Method, start)
}

func (h MyHandler) GetValue(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
	metricName := chi.URLParam(req, "name")

	metric, ok := h.Storage.GetMetrics(metricName)
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		h.logger.CreateRequestLog(req.RequestURI, req.Method, start)
		return
	}
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "%s", metric)
	h.logger.CreateRequestLog(req.RequestURI, req.Method, start)
}

func (h MyHandler) ListMetrics(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
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
	h.logger.CreateRequestLog(req.RequestURI, req.Method, start)
}
