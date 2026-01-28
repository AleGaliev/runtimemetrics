package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	models "github.com/AleGaliev/runtimemetrics/internal/model"
	"github.com/AleGaliev/runtimemetrics/internal/service/hash"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			var raw json.RawMessage
			if err := json.Unmarshal(bodyBytes, &raw); err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			var metricsList []models.Metrics
			if err := json.Unmarshal(raw, &metricsList); err == nil {
				for _, m := range metricsList {
					if err := MetricValidate(m); err != nil {
						res.WriteHeader(http.StatusBadRequest)
						return
					}
				}
				h.ServeHTTP(res, req)
				return
			}

			var metric models.Metrics
			if err := json.Unmarshal(raw, &metric); err == nil {
				if err := MetricValidate(metric); err != nil {
					res.WriteHeader(http.StatusBadRequest)
					return
				}
				h.ServeHTTP(res, req)
				return
			}

			res.WriteHeader(http.StatusBadRequest)
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
			h.ServeHTTP(res, req)
		}
		return http.HandlerFunc(fn)
	}
}

func IPValidateInterceptor(cidr string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		if cidr != "" {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Errorf(codes.PermissionDenied,
					"metadata not found in context")
			}
			ipMeta := md.Get("x-real-ip")
			if len(ipMeta) == 0 {
				return nil, status.Errorf(codes.PermissionDenied,
					"x-real-ip not found in metadata")
			}
			for _, ip := range ipMeta {
				isInRange, err := IsIPInCIDR(ip, cidr)
				if err != nil {
					return nil, status.Errorf(codes.PermissionDenied, "ip problem format")
				}
				if !isInRange {
					return nil, status.Errorf(codes.PermissionDenied,
						"ip not in cidr")
				}
			}
		}

		resp, err := handler(ctx, req)

		return resp, err
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
