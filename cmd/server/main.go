package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AleGaliev/runtimemetrics/internal/config/server"
	"github.com/AleGaliev/runtimemetrics/internal/handler"
	"github.com/AleGaliev/runtimemetrics/internal/logger"
	"github.com/AleGaliev/runtimemetrics/internal/observer"
	"github.com/AleGaliev/runtimemetrics/internal/repository"
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

	eventAudit := observer.NewEvent()
	auditFile, err := repository.CreateAuditSaveFile(serverConf.AuditFile)
	if err != nil {
		fmt.Println("Error creating audit file", err)
	} else {
		eventAudit.Register(auditFile)
	}
	auditSender, err := repository.NewAuditSender(serverConf.AuditURL)
	if err != nil {
		fmt.Println("Error creating audit sender", err)
	} else {
		eventAudit.Register(auditSender)
	}

	memStorage, err := srv.NewServerMemStorage(serverConf)
	if err != nil {
		panic(errors.Unwrap(err))
	}
	r := handler.CreateMyHandler(memStorage.MemStorage, memStorage.DBConfig, logServer, serverConf.HashKey, eventAudit)
	defer memStorage.DBConfig.Close()
	logServer.StartServerLog(serverConf.AdrHost)

	err = http.ListenAndServe(serverConf.AdrHost, r)
	if err != nil {
		panic(errors.Unwrap(err))
	}

}
