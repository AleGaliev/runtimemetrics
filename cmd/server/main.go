package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AleGaliev/kubercontroller/internal/config/server"
	"github.com/AleGaliev/kubercontroller/internal/filestore"
	"github.com/AleGaliev/kubercontroller/internal/handler"
	"github.com/AleGaliev/kubercontroller/internal/logger"
	"github.com/AleGaliev/kubercontroller/internal/postgresdb"
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

	var r http.Handler

	if serverConf.DatabaseDSN == "" {
		fileStore := filestore.NewFileStore(serverConf.FileStoragePath)

		memStorage, err := storage.CreateStorage(fileStore, serverConf.StoreInterval, serverConf.Restore)
		if err != nil {
			panic(err)
		}
		fmt.Println("mem storage created")
		if serverConf.StoreInterval > 0 && serverConf.DatabaseDSN == "" {
			go func() {
				for {
					time.Sleep(time.Duration(serverConf.StoreInterval) * time.Second)
					if err := memStorage.SaveMetricToFile(); err != nil {
						panic(err)
					}
				}
			}()
		}
		r = handler.CreateMyHandler(memStorage, logServer)
		defer memStorage.SaveMetricToFile()

	} else {
		dbConfig, err := postgresdb.NewPostgresDB(serverConf.DatabaseDSN)
		if err != nil {
			panic(err)
		}

		if err = dbConfig.Migrate(); err != nil {
			fmt.Println("failed to migrate database")
			fmt.Println(serverConf.DatabaseDSN)
		}
		r = handler.CreateMyHandler(dbConfig, logServer)
		defer dbConfig.Close()

	}
	logServer.StartServerLog(serverConf.AdrHost)

	err = http.ListenAndServe(serverConf.AdrHost, r)
	if err != nil {
		panic(err)
	}

}
