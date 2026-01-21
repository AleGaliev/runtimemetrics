package middleware

import (
	"net/http"
	"time"
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
