package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/AleGaliev/kubercontroller/internal/handler"
	"github.com/AleGaliev/kubercontroller/internal/storage"
)

func main() {
	adrHost := flag.String("a", "localhost:8080", "Endpoint http server")
	varAdrHost, ok := os.LookupEnv("ADDRESS")
	if ok {
		adrHost = &varAdrHost
	}
	flag.Parse()

	memStotage := storage.CreateStorage()
	r := handler.CreateMyHandler(memStotage)

	err := http.ListenAndServe(*adrHost, r)
	if err != nil {
		panic(err)
	}
}
