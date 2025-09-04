package handler

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

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

func CreateMyHandler(storage Storage, logger logger, db db) http.Handler {
	h := &MyHandler{
		Storage: storage,
		logger:  logger,
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
	muxGzip := h.GzipMiddlewareHandler(mux)
	muxMiddlewareLogger := h.MiddlewareHandlerLogger(muxGzip)

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

// получение метрики в формате json
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
		"status":  "success",
		"message": "Запрос обработан",
	}
	json.NewEncoder(res).Encode(response)
	//res.Write(metrics)
}

// Добавлении метрик в формате json
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
	res.Write(metrics)
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

func (h MyHandler) MiddlewareHandlerLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(res, req)
		h.logger.CreateRequestLog(req.RequestURI, req.Method, start)
	})
}

func (h MyHandler) GzipMiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(res, req)
			return
		}
		if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
			dgz, err := gzip.NewReader(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			req.Body = struct {
				io.Reader
				io.Closer
			}{dgz, req.Body}
			defer dgz.Close()
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(res, gzip.BestSpeed)
		if err != nil {
			io.WriteString(res, err.Error())
			return
		}
		defer gz.Close()

		res.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: res, Writer: gz}, req)
	})
}
