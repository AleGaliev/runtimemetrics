package main

import (
	"github.com/AleGaliev/kubercontroller/internal/handler"
	"net/http"
)

func main() {
	var h handler.MyHandler

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", h.ServeHTTP)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
