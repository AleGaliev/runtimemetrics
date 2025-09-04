package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AleGaliev/kubercontroller/internal/config/db"
	"github.com/AleGaliev/kubercontroller/internal/handler"
	"github.com/AleGaliev/kubercontroller/internal/logger"
	"github.com/AleGaliev/kubercontroller/internal/storage"
)

type serverConfig struct {
	adrHost         string
	storeInterval   int
	fileStoragePath string
	restore         bool
}

func createServerConfig() (serverConfig, error) {
	adrHost := flag.String("a", "localhost:8080", "Endpoint http server")
	storeInterval := flag.Int("i", 2, "interval save metrics in storage")
	fileStoragePath := flag.String("f", "storage.json", "filepath save metric storage")
	restore := flag.Bool("r", true, "read file storage metrics")
	//	databaseDSN := flag.String("d", "localhost", "filepath save metric storage")
	flag.Parse()

	varAdrHost, ok := os.LookupEnv("ADDRESS")
	if ok {
		adrHost = &varAdrHost
	}
	varStoreInterval, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		intStoreInterval, err := strconv.Atoi(varStoreInterval)
		if err != nil {
			return serverConfig{}, fmt.Errorf("error converting STORE_INTERVAL to int: %v", err)
		}
		storeInterval = &intStoreInterval
	}
	varFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		fileStoragePath = &varFileStoragePath
	}
	varRestore, ok := os.LookupEnv("RESTORE")
	if ok {
		boolVarRestore, err := strconv.ParseBool(varRestore)
		if err != nil {
			return serverConfig{}, fmt.Errorf("error converting RESTORE to bool: %v", err)
		}
		restore = &boolVarRestore
	}

	return serverConfig{
		adrHost:         *adrHost,
		storeInterval:   *storeInterval,
		fileStoragePath: *fileStoragePath,
		restore:         *restore,
	}, nil
}

func main() {
	serverConf, err := createServerConfig()
	if err != nil {
		panic(err)
	}

	logServer, err := logger.CreateLogger()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	memStotage, err := storage.CreateStorage(serverConf.fileStoragePath, serverConf.storeInterval, serverConf.restore)
	if err != nil {
		panic(err)
	}

	if serverConf.storeInterval > 0 {
		go func() {
			for {
				time.Sleep(time.Duration(serverConf.storeInterval) * time.Second)
				if err := memStotage.SaveMetricToFile(); err != nil {
					panic(err)
				}
			}
		}()
	}
	logServer.StartServerLog(serverConf.adrHost)
	dbConfig := db.NewDatabasePostgresConfig()
	r := handler.CreateMyHandler(memStotage, logServer, dbConfig)

	err = http.ListenAndServe(serverConf.adrHost, r)
	if err != nil {
		panic(err)
	}
}
