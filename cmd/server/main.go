package main

import (
	"net/http"
	"time"

	"github.com/AleGaliev/kubercontroller/internal/config/server"
	"github.com/AleGaliev/kubercontroller/internal/filestore"
	"github.com/AleGaliev/kubercontroller/internal/handler"
	"github.com/AleGaliev/kubercontroller/internal/logger"
	"github.com/AleGaliev/kubercontroller/internal/storage"
)

func main() {
	serverConf, err := server.NewServerConfig()
	if err != nil {
		panic(err)
	}

	logServer, err := logger.CreateLogger()
	if err != nil {
		panic(err)
	}
	fileStore := filestore.NewFileStore(serverConf.FileStoragePath)

	memStorage, err := storage.CreateStorage(fileStore, serverConf.StoreInterval, serverConf.Restore)
	if err != nil {
		panic(err)
	}

	if serverConf.StoreInterval > 0 {
		go func() {
			for {
				time.Sleep(time.Duration(serverConf.StoreInterval) * time.Second)
				if err := memStorage.SaveMetricToFile(); err != nil {
					panic(err)
				}
			}
		}()
	}
	logServer.StartServerLog(serverConf.AdrHost)
	r := handler.CreateMyHandler(memStorage, logServer)

	err = http.ListenAndServe(serverConf.AdrHost, r)
	if err != nil {
		panic(err)
	}
	defer memStorage.SaveMetricToFile()
}
