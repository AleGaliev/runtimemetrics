package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/AleGaliev/runtimemetrics/internal/agent"
	"github.com/AleGaliev/runtimemetrics/internal/logger"
	"github.com/AleGaliev/runtimemetrics/internal/repository"
	"github.com/AleGaliev/runtimemetrics/internal/service/cripto"
	"github.com/AleGaliev/runtimemetrics/internal/service/retry"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
	serviceName  string = "agent"

	defaultPollInterval   int    = 2
	defaultReportInterval int    = 10
	defaultRateLimit      int    = 10
	defaultBaseURL        string = "localhost:8080"
	defaultHashKey        string = ""
	defaultCryptoKey      string = ""
)

type flagsAgent struct {
	baseURL        string `json:"address"`
	hashKey        string `json:"hash_key"`
	cryptoKey      string `json:"crypto_key"`
	pollInterval   int    `json:"poll_interval"`
	reportInterval int    `json:"report_interval"`
	RateLimit      int    `json:"rate_limit"`
}

func main() {
	logServer, err := logger.CreateLogger()
	if err != nil {
		panic(errors.Unwrap(err))
	}

	logServer.CreateVersionLog(serviceName, buildVersion, buildDate, buildCommit)

	arg, err := initConfig()
	if err != nil {
		panic(errors.Unwrap(err))
	}
	pubKey := &cripto.Cripto{}
	if arg.cryptoKey != "" {
		pubKey, err = cripto.NewCripto("", arg.cryptoKey)
		if err != nil {
			panic(errors.Unwrap(err))
		}
	}

	clientCfg := repository.NewClientConfig(
		repository.WithLogger(logServer),
		repository.WithURL(arg.baseURL),
		repository.WithKeyHash(arg.hashKey),
		repository.WithCripto(pubKey),
	)

	agentCfg, err := agent.NewAgentConfig(clientCfg, retry.CreateRetry(), arg.pollInterval, arg.reportInterval, arg.RateLimit)
	if err != nil {
		log.Fatalf("error parsing agent config: %v", errors.Unwrap(err))
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	defer cancel()
	if err = agentCfg.Run(ctx); err != nil {
		panic(err)
	}
}

func initConfig() (flagsAgent, error) {
	baseURL := flag.String("a", defaultBaseURL, "Endpoint http server")
	pollInterval := flag.Int("p", defaultPollInterval, "Interval poll metrics")
	reportInterval := flag.Int("r", defaultReportInterval, "Interval report metrics")
	rateLimit := flag.Int("l", defaultRateLimit, "Interval poll metrics")
	hashKey := flag.String("k", defaultHashKey, "key agent encryption")
	cryptoKey := flag.String("crypto-key", defaultCryptoKey, "key agent encryption")
	fileConfig := flag.String("c", "", "config file")
	flag.Parse()

	varFileConfig, ok := os.LookupEnv("CONFIG")
	if ok {
		fileConfig = &varFileConfig
	}
	varAdrHost, ok := os.LookupEnv("ADDRESS")
	if ok {
		baseURL = &varAdrHost
	}
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
	varCryptoKey, ok := os.LookupEnv("CRYPTO_KEY")
	if ok {
		cryptoKey = &varCryptoKey
	}

	resultFlags := flagsAgent{
		baseURL:        *baseURL,
		pollInterval:   *pollInterval,
		reportInterval: *reportInterval,
		hashKey:        *hashKey,
		RateLimit:      *rateLimit,
		cryptoKey:      *cryptoKey,
	}

	if *fileConfig != "" {
		confiFlagsInFile, err := readConfigFile(*fileConfig)
		if err != nil {
			return flagsAgent{}, err
		}
		resultFlags.convertFlagsResult(confiFlagsInFile)
	}

	return resultFlags, nil
}

func readConfigFile(filePath string) (flagsAgent, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0o666)
	if err != nil {
		return flagsAgent{}, fmt.Errorf("could not open config file: %w", err)
	}
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return flagsAgent{}, fmt.Errorf("could not read config file: %w", scanner.Err())
	}

	data := scanner.Bytes()
	var flags flagsAgent
	err = json.Unmarshal(data, &flags)
	if err != nil {
		return flagsAgent{}, fmt.Errorf("could not open config file: %w", err)
	}

	return flags, nil
}

func (flags *flagsAgent) convertFlagsResult(flagsInFile flagsAgent) {
	if flags.baseURL == defaultBaseURL {
		flags.baseURL = flagsInFile.baseURL
	}
	if flags.hashKey == defaultHashKey {
		flags.hashKey = flagsInFile.hashKey
	}
	if flags.cryptoKey == defaultCryptoKey {
		flags.cryptoKey = flagsInFile.cryptoKey
	}
	if flags.pollInterval == defaultPollInterval {
		flags.pollInterval = flagsInFile.pollInterval
	}
	if flags.reportInterval == defaultReportInterval {
		flags.reportInterval = flagsInFile.reportInterval
	}
	if flags.RateLimit == defaultRateLimit {
		flags.RateLimit = flagsInFile.RateLimit
	}
}
