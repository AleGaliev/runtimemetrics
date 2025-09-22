package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AleGaliev/runtimemetrics/internal/agent"
	"github.com/AleGaliev/runtimemetrics/internal/logger"
	"github.com/AleGaliev/runtimemetrics/internal/repository"
	"github.com/AleGaliev/runtimemetrics/internal/service/retry"
)

type flagsAgent struct {
	baseURL        string
	pollInterval   int
	reportInterval int
}

func main() {
	logServer, err := logger.CreateLogger()
	if err != nil {
		panic(errors.Unwrap(err))
	}

	arg, err := initConfig()
	if err != nil {
		panic(errors.Unwrap(err))
	}
	clientCfg := repository.NewClientConfig(logServer, arg.baseURL)

	agentCfg, err := agent.NewAgentConfig(clientCfg, retry.CreateRetry(), arg.pollInterval, arg.reportInterval)

	if err != nil {
		log.Fatalf("error parsing agent config: %v", errors.Unwrap(err))
	}

	for {
		if err := agentCfg.Run(); err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func initConfig() (flagsAgent, error) {
	baseURL := flag.String("a", "localhost:8080", "Endpoint http server")
	varAdrHost, ok := os.LookupEnv("ADDRESS")
	if ok {
		baseURL = &varAdrHost
	}

	pollInterval := flag.Int("p", 2, "Interval poll metrics")
	reportInterval := flag.Int("r", 10, "Interval report metrics")
	varPollInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		StrPollInterval, err := strconv.Atoi(varPollInterval)
		if err != nil {
			return flagsAgent{}, fmt.Errorf("error converting POLL_INTERVAL to int: %v", err)
		}
		pollInterval = &StrPollInterval
	}
	varReportInterval, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		StrReportInterval, err := strconv.Atoi(varReportInterval)
		if err != nil {
			return flagsAgent{}, fmt.Errorf("error converting REPORT_INTERVAL to int: %v", err)
		}
		reportInterval = &StrReportInterval
	}
	flag.Parse()

	return flagsAgent{
		baseURL:        *baseURL,
		pollInterval:   *pollInterval,
		reportInterval: *reportInterval,
	}, nil
}
