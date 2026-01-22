package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	models "github.com/AleGaliev/runtimemetrics/internal/model"
	"github.com/AleGaliev/runtimemetrics/internal/service/hash"
)

func MetricValidateMiddleware(keyHash string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			if req.Method == http.MethodGet || req.Header.Get("Content-Type") != "application/json" {
				h.ServeHTTP(res, req)
				return
			}
			res.Header().Set("Content-Type", "application/json")

			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			verifiableHash := req.Header.Get("HashSHA256")

			if !hash.CheckHash(keyHash, verifiableHash, bodyBytes) && verifiableHash != "" {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			req.Body.Close()

			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			validationBody := io.NopCloser(bytes.NewBuffer(bodyBytes))

			data := json.NewDecoder(validationBody)
			var metricsData []models.Metrics

			if err := data.Decode(&metricsData); err == nil {
				for _, m := range metricsData {
					if err := MetricValidate(m); err != nil {
						res.WriteHeader(http.StatusBadRequest)
						return
					}
				}
				h.ServeHTTP(res, req)
				return
			}

			validationBody = io.NopCloser(bytes.NewBuffer(bodyBytes))
			data = json.NewDecoder(validationBody)
			var metrics models.Metrics
			if err := data.Decode(&metrics); err == nil {
				if err := MetricValidate(metrics); err != nil {
					res.WriteHeader(http.StatusBadRequest)
					return
				}
			} else {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			h.ServeHTTP(res, req)
		}
		return http.HandlerFunc(fn)
	}
}

func MetricValidate(metric models.Metrics) error {
	if metric.ID == "" {
		return fmt.Errorf("metric id or value is required")
	}

	if metric.Delta == nil && metric.Value == nil {
		return nil
	}
	switch metric.MType {

	case models.Gauge:
		if metric.Value == nil {
			return fmt.Errorf("metrics value is nil")
		}

	case models.Counter:
		if metric.Delta == nil {
			return fmt.Errorf("metrics delta is nil")
		}

	default:
		return fmt.Errorf("invalid metric type: %s", metric.MType)
	}

	return nil
}
func IPValidateMiddleware(cidr string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			if cidr != "" {
				ip := req.Header.Get("X-Real-IP")
				if ip == "" {
					res.WriteHeader(http.StatusForbidden)
				}
				isInRange, err := IsIPInCIDR(ip, cidr)
				if err != nil {
					res.WriteHeader(http.StatusForbidden)
					return
				}
				if !isInRange {
					res.WriteHeader(http.StatusForbidden)
				}
			}
		}
		return http.HandlerFunc(fn)
	}
}

func IsIPInCIDR(ipStr, cidrStr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("ip adress not correct: %s", ipStr)
	}

	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false, fmt.Errorf("CIDR not correct: %s", cidrStr)
	}

	return ipNet.Contains(ip), nil
}
