package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/AleGaliev/runtimemetrics/internal/agent"
	"github.com/AleGaliev/runtimemetrics/internal/logger"
	"github.com/AleGaliev/runtimemetrics/internal/repository"
	"github.com/AleGaliev/runtimemetrics/internal/service/retry"
)

type flagsAgent struct {
	baseURL        string
	pollInterval   int
	reportInterval int
	hashKey        string
	RateLimit      int
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
	clientCfg := repository.NewClientConfig(logServer, arg.baseURL, arg.hashKey)

	agentCfg, err := agent.NewAgentConfig(clientCfg, retry.CreateRetry(), arg.pollInterval, arg.reportInterval, arg.RateLimit)
	if err != nil {
		log.Fatalf("error parsing agent config: %v", errors.Unwrap(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err = agentCfg.Run(ctx); err != nil {
		panic(err)
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
	rateLimit := flag.Int("l", 10, "Interval poll metrics")
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
	hashKey := flag.String("k", "", "key agent encryption")
	varHashKey, ok := os.LookupEnv("KEY")
	if ok {
		hashKey = &varHashKey
	}
	varRateLimit, ok := os.LookupEnv("RATE_LIMIT")
	if ok {
		StrRateLimit, err := strconv.Atoi(varRateLimit)
		if err != nil {
			return flagsAgent{}, fmt.Errorf("error converting RATE_LIMIT to int: %v", err)
		}
		rateLimit = &StrRateLimit
	}
	flag.Parse()

	return flagsAgent{
		baseURL:        *baseURL,
		pollInterval:   *pollInterval,
		reportInterval: *reportInterval,
		hashKey:        *hashKey,
		RateLimit:      *rateLimit,
	}, nil
}
