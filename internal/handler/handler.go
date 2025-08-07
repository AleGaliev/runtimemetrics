package handler

import (
	models "github.com/AleGaliev/kubercontroller/internal/model"
	"net/http"
	"strconv"
	"strings"
)

type MyHandler struct{}

func (h MyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pathUrl := strings.Split(strings.Trim(req.URL.Path, "/"), "/")

	if len(pathUrl) < 4 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	name := pathUrl[2]
	MType := pathUrl[1]
	value := pathUrl[3]

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
