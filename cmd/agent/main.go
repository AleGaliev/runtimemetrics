package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/AleGaliev/runtimemetrics/internal/agent"
	"github.com/AleGaliev/runtimemetrics/internal/logger"
	"github.com/AleGaliev/runtimemetrics/internal/repository"
	"github.com/AleGaliev/runtimemetrics/internal/service/crypto"
	"github.com/AleGaliev/runtimemetrics/internal/service/retry"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
	serviceName  string = "agent"

	defaultPollInterval   int    = 10
	defaultReportInterval int    = 2
	defaultRateLimit      int    = 10
	defaultBaseURL        string = "localhost:8080"
	defaultHashKey        string = ""
	defaultCryptoKey      string = ""
	defaultAdrHostGrpc    string = ""
)

type flagsAgent struct {
	BaseURL        string `json:"address"`
	AdrHostGrpc    string `json:"address_grpc"`
	HashKey        string `json:"hash_key"`
	CryptoKey      string `json:"crypto_key"`
	PollInterval   int    `json:"poll_interval"`
	ReportInterval int    `json:"report_interval"`
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
	pubKey := &crypto.Crypto{}
	if arg.CryptoKey != "" {
		pubKey, err = crypto.NewCrypto("", arg.CryptoKey)
		if err != nil {
			panic(errors.Unwrap(err))
		}
	}

	clientCfg := repository.NewClientConfig(
		repository.WithLogger(logServer),
		repository.WithURL(arg.BaseURL),
		repository.WithGrpcClient(arg.AdrHostGrpc),
		repository.WithKeyHash(arg.HashKey),
		repository.WithCrypto(pubKey),
	)

	agent := agent.New(*clientCfg,
		agent.WithBaseURL(arg.BaseURL),
		agent.WithGRPCHost(arg.AdrHostGrpc),
		agent.WithRetry(retry.CreateRetry()),
		agent.WithWorkers(arg.RateLimit),
		agent.WithReportInterval(arg.ReportInterval),
		agent.WithPollInterval(arg.PollInterval),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	defer cancel()
	if err = agent.Run(ctx); err != nil {
		panic(err)
	}
}

func initConfig() (flagsAgent, error) {
	baseURL := flag.String("a", defaultBaseURL, "Endpoint http server")
	adrHostGrpc := flag.String("ag", defaultAdrHostGrpc, "Endpoint grpc server")
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
	varAdrHostGrpc, ok := os.LookupEnv("ADDRESS_GRPC")
	if ok {
		adrHostGrpc = &varAdrHostGrpc
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
		BaseURL:        *baseURL,
		AdrHostGrpc:    *adrHostGrpc,
		PollInterval:   *pollInterval,
		ReportInterval: *reportInterval,
		HashKey:        *hashKey,
		RateLimit:      *rateLimit,
		CryptoKey:      *cryptoKey,
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
	if flags.BaseURL == defaultBaseURL {
		flags.BaseURL = flagsInFile.BaseURL
	}
	if flags.HashKey == defaultHashKey {
		flags.HashKey = flagsInFile.HashKey
	}
	if flags.CryptoKey == defaultCryptoKey {
		flags.CryptoKey = flagsInFile.CryptoKey
	}
	if flags.PollInterval == defaultPollInterval {
		flags.PollInterval = flagsInFile.PollInterval
	}
	if flags.ReportInterval == defaultReportInterval {
		flags.ReportInterval = flagsInFile.ReportInterval
	}
	if flags.RateLimit == defaultRateLimit {
		flags.RateLimit = flagsInFile.RateLimit
	}
	if flags.AdrHostGrpc == defaultAdrHostGrpc {
		flags.AdrHostGrpc = flagsInFile.AdrHostGrpc
	}
}
