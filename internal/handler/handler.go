package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/AleGaliev/kubercontroller/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type db interface {
	PingPostgres() error
}

type logger interface {
	CreateRequestLog(url, method string, timestamp time.Time)
}

type Storage interface {
	AddMetric(myType, name, value string) error
	GetMetrics(name string) (string, bool)
	GetAllMetric() string
	UpdateMetrics(r io.Reader) error
	ValueMetrics(r io.Reader) ([]byte, bool, error)
}
type MyHandler struct {
	Storage Storage
	logger  logger
	db      db
}

func CreateMyHandler(storage Storage, logger middleware.Logger, db db) http.Handler {
	h := &MyHandler{
		Storage: storage,
		db:      db,
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

	mux.Get("/", h.ListMetrics)
	mux.Get("/ping", h.GetPing)

	muxGzip := middleware.GzipMiddlewareHandler(mux)
	muxMiddlewareLogger := middleware.MiddlewareHandlerLogger(muxGzip, logger)

	return muxMiddlewareLogger
}


func (h MyHandler) GetPing(res http.ResponseWriter, _ *http.Request) {
	if err := h.db.PingPostgres(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":  "success",
		"message": "Postgres ping succeeded",
	}
	json.NewEncoder(res).Encode(response)
}

// ServeHTTPUpdate добавление метрики в формате json
func (h MyHandler) ServeHTTPUpdate(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	if req.Method != http.MethodPost || req.Header.Get("Content-Type") != "application/json" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Storage.UpdateMetrics(req.Body); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":  `success`,
		"message": "Запрос обработан",
	}
	if err := json.NewEncoder(res).Encode(response); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}

}

// ServeHTTPValue получение метрик в формате json
func (h MyHandler) ServeHTTPValue(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	if req.Method != http.MethodPost || req.Header.Get("Content-Type") != "application/json" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	metrics, ok, err := h.Storage.ValueMetrics(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
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

	err := h.Storage.AddMetric(myType, name, value)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h MyHandler) GetValue(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")

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
	res.Header().Set("Content-Type", "text/html")
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
