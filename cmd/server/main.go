package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/AleGaliev/runtimemetrics/internal/server"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
	serviceName  string = "server"
)

func main() {
	srv, err := server.New()
	if err != nil {
		log.Fatal(errors.Unwrap(err))
	}
	defer srv.Close()

	srv.LogServer.CreateVersionLog(serviceName, buildVersion, buildDate, buildCommit)
	ctx, cansel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	defer cansel()

	go func() {
		if err := srv.StartHttpServer(); err != nil {
			log.Fatal(errors.Unwrap(err))
		}
	}()

	go func() {
		if err := srv.StartGrpcServer(); err != nil {
			log.Fatal(errors.Unwrap(err))
		}
	}()

	<-ctx.Done()
	if err := srv.Stop(ctx); err != nil {
		log.Fatal(errors.Unwrap(err))
	}
}
