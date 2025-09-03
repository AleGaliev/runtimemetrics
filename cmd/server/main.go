package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/AleGaliev/kubercontroller/internal/handler"
	"github.com/AleGaliev/kubercontroller/internal/logger"
	"github.com/AleGaliev/kubercontroller/internal/storage"
)

func main() {
	adrHost := flag.String("a", "localhost:8080", "Endpoint http server")
	storeInterval := flag.Int("i", 300, "interval save metrics in storage")
	fileStoragePath := flag.String("f", "storage.json", "filepath save metric storage")
	restore := flag.Bool("r", false, "read file storage metrics")
	varAdrHost, ok := os.LookupEnv("ADDRESS")
	if ok {
		adrHost = &varAdrHost
	}
	flag.Parse()
	logServer, err := logger.CreateLogger()
	if err != nil {
		panic(err)
	}
	memStotage := storage.CreateStorage(*fileStoragePath, *storeInterval)
	if *restore {
		if err = memStotage.ReadMetricInFile(); err != nil {
			panic(err)
		}
	}
	if *storeInterval > 0 {
		go func() {
			for {
				time.Sleep(time.Duration(*storeInterval) * time.Second)
				if err := memStotage.ReadMetricInFile(); err != nil {
					panic(err)
				}
			}
		}()
	}
	logServer.StartServerLog(*adrHost)
	r := handler.CreateMyHandler(memStotage, logServer)

	err = http.ListenAndServe(*adrHost, r)
	if err != nil {
		panic(err)
	}
}
