package middleware

import (
	"bytes"
	"io"
	"net/http"
)

type decrypt interface {
	Decrypt(data []byte) ([]byte, error)
}

func MiddlewareDecrypt(decrypt decrypt) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Content-Type") != "application/json" {
				h.ServeHTTP(res, req)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			bodyBytes, err = decrypt.Decrypt(bodyBytes)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			h.ServeHTTP(res, req)
		}
		return http.HandlerFunc(fn)
	}
}
