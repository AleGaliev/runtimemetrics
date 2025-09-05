package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AleGaliev/kubercontroller/internal/agent"
	"github.com/AleGaliev/kubercontroller/internal/logger"
	"github.com/AleGaliev/kubercontroller/internal/repository"
)

type flagsAgent struct {
	baseURL        string
	pollInterval   int
	reportInterval int
}

func main() {
	logServer, err := logger.CreateLogger()
	if err != nil {
		panic(err)
	}

	arg, err := initСonfig()
	if err != nil {
		panic(err)
	}
	clientCfg := repository.NewClientConfig(logServer, arg.baseURL)
	agentCfg, err := agent.NewAgentConfig(clientCfg, arg.pollInterval, arg.reportInterval)

	if err != nil {
		log.Fatalf("error parsing agent config: %v", err)
	}

	for {
		if err := agentCfg.Run(); err != nil {
			fmt.Println(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func initСonfig() (flagsAgent, error) {
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
