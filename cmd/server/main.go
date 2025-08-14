package main

import (
	"net/http"

	"github.com/AleGaliev/kubercontroller/internal/handler"
	"github.com/go-chi/chi/v5"
)

//func main() {
//	var h handler.MyHandler
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/update/", h.ServeHTTP)
//
//	err := http.ListenAndServe(`:8080`, mux)
//	if err != nil {
//		panic(err)
//	}
//}

func main() {
	var h handler.MyHandler

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", h.ServeHTTP)
	r.Get("/value/{type}/{name}", h.GetValue)
	r.Get("/", h.ListMetrics)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
