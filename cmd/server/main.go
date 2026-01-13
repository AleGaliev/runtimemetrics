package main

import (
	"errors"
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/AleGaliev/runtimemetrics/internal/config/server"
	"github.com/AleGaliev/runtimemetrics/internal/handler"
	"github.com/AleGaliev/runtimemetrics/internal/logger"
	"github.com/AleGaliev/runtimemetrics/internal/observer"
	"github.com/AleGaliev/runtimemetrics/internal/repository"
	srv "github.com/AleGaliev/runtimemetrics/internal/server"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
	serviceName  string = "server"
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

	logServer.CreateVersionLog(serviceName, buildVersion, buildDate, buildCommit)

	eventAudit := createEventAudit(serverConf.AuditFile, serverConf.AuditURL)

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

func createEventAudit(auditFile, auditURL string) *observer.Event {
	eventAudit := observer.NewEvent()
	eventAuditFile, err := repository.CreateAuditSaveFile(auditFile)
	if err != nil {
		fmt.Println("Error creating audit file", err)
	} else {
		eventAudit.Register(eventAuditFile)
	}

	eventAuditSender, err := repository.NewAuditSender(auditURL)
	if err != nil {
		fmt.Println("Error creating audit sender", err)
	} else {
		eventAudit.Register(eventAuditSender)
	}

	return eventAudit
}
