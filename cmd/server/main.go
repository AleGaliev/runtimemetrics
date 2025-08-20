package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/AleGaliev/kubercontroller/internal/handler"
	"github.com/AleGaliev/kubercontroller/internal/logger"
	"github.com/AleGaliev/kubercontroller/internal/storage"
)

func main() {
	adrHost := flag.String("a", "localhost:8080", "Endpoint http server")
	varAdrHost, ok := os.LookupEnv("ADDRESS")
	if ok {
		adrHost = &varAdrHost
	}
	flag.Parse()
	logServer, err := logger.CreateLogger()
	if err != nil {
		panic(err)
	}
	memStotage := storage.CreateStorage()
	logServer.StartServerLog(*adrHost)
	r := handler.CreateMyHandler(memStotage, logServer)

	err = http.ListenAndServe(*adrHost, r)
	if err != nil {
		panic(err)
	}
}
