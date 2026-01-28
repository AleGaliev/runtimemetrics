package middleware

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

type Logger interface {
	CreateRequestLog(url, method string, timestamp time.Time)
	CreateErrorLog(service, message string)
}

func MiddlewareHandlerLogger(logger Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			start := time.Now()
			h.ServeHTTP(res, req)
			logger.CreateRequestLog(req.RequestURI, req.Method, start)
		}
		return http.HandlerFunc(fn)
	}
}

func LoggingInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		start := time.Now()

		resp, err := handler(ctx, req)
		logger.CreateRequestLog("grpc", info.FullMethod, start)

		return resp, err
	}
}
