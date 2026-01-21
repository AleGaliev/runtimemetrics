package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/AleGaliev/runtimemetrics/internal/audit"
	models "github.com/AleGaliev/runtimemetrics/internal/model"
	"github.com/AleGaliev/runtimemetrics/internal/observer"
)

const (
	serviceNameAuditMiddleware = "AuditMiddleware"
)

func AuditMiddleware(event *observer.Event, logger Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			var bodyBytes []byte
			var metricsData []models.Metrics

			if req.Body != nil {
				bodyBytes, _ = io.ReadAll(req.Body)
				req.Body.Close()

				// Восстанавливаем Body для основного обработчика
				req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Декодируем JSON
				if err := json.Unmarshal(bodyBytes, &metricsData); err != nil {
					// Логируем ошибку, но не прерываем выполнение
					logger.CreateErrorLog(serviceNameAuditMiddleware, fmt.Sprintf("Error unmarshalling request body: %s", err))
				}
			}
			ip, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				ip = req.RemoteAddr
			}

			h.ServeHTTP(res, req)

			audit := audit.CreateAuditMessage(ip, metricsData)
			event.Notify(audit)
		}
		return http.HandlerFunc(fn)
	}
}
