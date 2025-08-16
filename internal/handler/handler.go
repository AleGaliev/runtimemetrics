package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Storage interface {
	AddMetric(myType, name, value string) error
	GetMetrics(name string) (string, bool)
	GetAllMetric() string
}
type MyHandler struct {
	Storage Storage
}

func CreateMyHandler(storage Storage) http.Handler {
	h := &MyHandler{
		Storage: storage,
	}

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", h.ServeHTTP)
	r.Get("/value/{type}/{name}", h.GetValue)
	r.Get("/", h.ListMetrics)

	return r
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

func (h MyHandler) ListMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	body := h.Storage.GetAllMetric()

	fmt.Fprint(w, `
    <!DOCTYPE html>
    <html>
    <body>
        <h1>Metrics List</h1>
		<ul>
    `)
	fmt.Fprintf(w, `%s`, body)

	fmt.Fprint(w, `
	</ul>
    </body>
    </html>
    `)
}
