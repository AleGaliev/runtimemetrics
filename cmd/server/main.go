package main

import (
	"errors"
	"net/http"

	"github.com/AleGaliev/runtimemetrics/internal/config/server"
	"github.com/AleGaliev/runtimemetrics/internal/handler"
	"github.com/AleGaliev/runtimemetrics/internal/logger"
	srv "github.com/AleGaliev/runtimemetrics/internal/server"
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

	memStorage, err := srv.NewServerMemStorage(serverConf)
	if err != nil {
		panic(errors.Unwrap(err))
	}
	r := handler.CreateMyHandler(memStorage.MemStorage, memStorage.DBConfig, logServer)
	defer memStorage.DBConfig.Close()
	logServer.StartServerLog(serverConf.AdrHost)

	err = http.ListenAndServe(serverConf.AdrHost, r)
	if err != nil {
		panic(errors.Unwrap(err))
	}

}
