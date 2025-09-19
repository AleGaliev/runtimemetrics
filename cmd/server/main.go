package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AleGaliev/runtimemetrics/internal/config/db"
	"github.com/AleGaliev/runtimemetrics/internal/config/server"
	"github.com/AleGaliev/runtimemetrics/internal/filestore"
	"github.com/AleGaliev/runtimemetrics/internal/handler"
	"github.com/AleGaliev/runtimemetrics/internal/logger"
	"github.com/AleGaliev/runtimemetrics/internal/service/retry"
	"github.com/AleGaliev/runtimemetrics/internal/storage"
)

func main() {
	serverConf, err := server.NewServerConfig()
	if err != nil {
		panic(err)
	}

	logServer, err := logger.CreateLogger()
	if err != nil {
		panic(errors.Unwrap(err))
	}

	var r http.Handler

	if serverConf.DatabaseDSN == "" {
		fileStore := filestore.NewFileStore(serverConf.FileStoragePath)

		memStorage, err := storage.CreateStorage(fileStore, serverConf.StoreInterval, serverConf.Restore)
		if err != nil {
			panic(errors.Unwrap(err))
		}
		fmt.Println("mem storage created")
		if serverConf.StoreInterval > 0 && serverConf.DatabaseDSN == "" {
			go func() {
				for {
					time.Sleep(time.Duration(serverConf.StoreInterval) * time.Second)
					if err := memStorage.SaveMetricToFile(); err != nil {
						panic(errors.Unwrap(err))
					}
				}
			}()
		}
		r = handler.CreateMyHandler(memStorage, logServer)
		defer memStorage.SaveMetricToFile()

	} else {
		dbConfig, err := db.NewPostgresDB(serverConf.DatabaseDSN)
		if err != nil {
			panic(errors.Unwrap(err))
		}
		dbMemStorage := storage.NewPostgresDBStorage(dbConfig, retry.CreateRetry())

		if err = dbMemStorage.CreateMigration(); err != nil {
			panic(errors.Unwrap(err))
		}
		r = handler.CreateMyHandler(dbMemStorage, logServer)
		defer dbMemStorage.Close()

	}
	logServer.StartServerLog(serverConf.AdrHost)

	err = http.ListenAndServe(serverConf.AdrHost, r)
	if err != nil {
		panic(errors.Unwrap(err))
	}

}
