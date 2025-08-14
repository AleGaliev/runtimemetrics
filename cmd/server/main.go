package main

import (
	"flag"
	"net/http"

	"github.com/AleGaliev/kubercontroller/internal/handler"
	"github.com/go-chi/chi/v5"
)

var (
	adrHost = flag.String("a", "localhost:8080", "Endpoint http server")
)

func main() {
	flag.Parse()
	var h handler.MyHandler

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", h.ServeHTTP)
	r.Get("/value/{type}/{name}", h.GetValue)
	r.Get("/", h.ListMetrics)

	err := http.ListenAndServe(*adrHost, r)
	if err != nil {
		panic(err)
	}
}
