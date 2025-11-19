package middleware

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/AleGaliev/runtimemetrics/internal/audit"
	models "github.com/AleGaliev/runtimemetrics/internal/model"
	"github.com/AleGaliev/runtimemetrics/internal/observer"
)

func AuditMiddleware(event *observer.Event) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			h.ServeHTTP(res, req)
			data := json.NewDecoder(req.Body)
			ip, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				ip = req.RemoteAddr
			}

			var metricsData []models.Metrics
			if err := data.Decode(&metricsData); err == nil {
				fmt.Println(metricsData)
			}
			audit := audit.CreateAuditMessage(metricsData, ip)
			event.Notify(audit)
		}
		return http.HandlerFunc(fn)
	}
}
