package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	models "github.com/AleGaliev/kubercontroller/internal/model"
	"github.com/go-chi/chi/v5"
)

type MyHandler struct{}

func (h MyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	pathURL := strings.Split(strings.Trim(req.URL.Path, "/"), "/")

	if len(pathURL) < 4 {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, "404 page not found")
		return
	}
	name := chi.URLParam(req, "name")
	MType := chi.URLParam(req, "type")
	value := chi.URLParam(req, "value")

	switch MType {
	case models.Gauge:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
		}
		models.MemStorage[name] = models.Metrics{
			ID:    name,
			MType: MType,
			Value: &f,
		}

	case models.Counter:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if metrics, exists := models.MemStorage[name]; exists {
			*metrics.Delta += i
		} else {
			models.MemStorage[name] = models.Metrics{
				ID:    name,
				MType: MType,
				Delta: &i,
			}
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
	}
}

func (h MyHandler) GetValue(res http.ResponseWriter, req *http.Request) {
	metricName := chi.URLParam(req, "name")

	metric, ok := models.MemStorage[metricName]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	switch metric.MType {
	case models.Gauge:
		fmt.Fprintf(res, "%g", *metric.Value)
	case models.Counter:
		fmt.Fprintf(res, "%d", *metric.Delta)
	}
	res.WriteHeader(http.StatusOK)
}

func (h MyHandler) ListMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprint(w, `
    <!DOCTYPE html>
    <html>
    <body>
        <h1>Metrics List</h1>
    `)

	for _, m := range models.MemStorage {
		switch m.MType {
		case models.Gauge:
			fmt.Fprintf(w, `<p> %s: %g</p>`, m.ID, *m.Value)
		case models.Counter:
			fmt.Fprintf(w, `<p> %s: %d</p>`, m.ID, *m.Delta)
		}
	}

	fmt.Fprint(w, `
    </body>
    </html>
    `)
}
