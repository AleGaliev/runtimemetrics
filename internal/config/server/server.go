package server

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type ServerConfig struct {
	AdrHost         string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	DatabaseDSN     string
}

func NewServerConfig() (ServerConfig, error) {
	adrHost := flag.String("a", "localhost:8080", "Endpoint http server")
	storeInterval := flag.Int("i", 2, "interval save metrics in storage")
	fileStoragePath := flag.String("f", "storage.json", "filepath save metric storage")
	databaseDSN := flag.String("d", "", "database DSN")
	restore := flag.Bool("r", true, "read file storage metrics")
	flag.Parse()

	varAdrHost, ok := os.LookupEnv("ADDRESS")
	if ok {
		adrHost = &varAdrHost
	}
	varStoreInterval, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		intStoreInterval, err := strconv.Atoi(varStoreInterval)
		if err != nil {
			return ServerConfig{}, fmt.Errorf("error converting STORE_INTERVAL to int: %v", err)
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
			return ServerConfig{}, fmt.Errorf("error converting RESTORE to bool: %v", err)
		}
		restore = &boolVarRestore
	}
	varDatabaseDSN, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		databaseDSN = &varDatabaseDSN
	}
	return ServerConfig{
		AdrHost:         *adrHost,
		StoreInterval:   *storeInterval,
		FileStoragePath: *fileStoragePath,
		Restore:         *restore,
		DatabaseDSN:     *databaseDSN,
	}, nil
}
