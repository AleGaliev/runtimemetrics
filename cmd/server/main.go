package main

import (
	"flag"
	"net/http"

	"github.com/AleGaliev/kubercontroller/internal/handler"
	"github.com/AleGaliev/kubercontroller/internal/storage"
)

func main() {
	adrHost := flag.String("a", "localhost:8080", "Endpoint http server")
	flag.Parse()

	memStotage := storage.CreateStorage()
	r := handler.CreateMyHandler(memStotage)

	err := http.ListenAndServe(*adrHost, r)
	if err != nil {
		panic(err)
	}
}
