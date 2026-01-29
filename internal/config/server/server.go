package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	defaultAdrHost         string = "localhost:8080"
	defaultAdrHostGrpc     string = ""
	defaultFileStoragePath string = "storage.json"
	defaultDatabaseDSN     string = ""
	defaultHashKey         string = ""
	defaultAuditFile       string = ""
	defaultAuditURL        string = ""
	defaultCryptoKey       string = ""
	defaultStoreInterval   int    = 2
	defaultRestore         bool   = true
	defaultTrustedSubnet   string = ""
)

type ServerConfig struct {
	AdrHost         string `json:"address"`
	AdrHostGrpc     string `json:"address_grpc"`
	FileStoragePath string `json:"store_file"`
	DatabaseDSN     string `json:"database_dsn"`
	HashKey         string `json:"hash_key"`
	AuditFile       string `json:"audit_file"`
	AuditURL        string `json:"audit_url"`
	CryptoKey       string `json:"crypto_key"`
	StoreInterval   int    `json:"store_interval"`
	Restore         bool   `json:"restore"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

func NewServerConfig() (ServerConfig, error) {
	adrHost := flag.String("a", defaultAdrHost, "Endpoint http server")
	adrHostGrpc := flag.String("ag", defaultAdrHostGrpc, "Endpoint grpc server")
	storeInterval := flag.Int("i", defaultStoreInterval, "interval save metrics in storage")
	fileStoragePath := flag.String("f", defaultFileStoragePath, "filepath save metric storage")
	databaseDSN := flag.String("d", defaultDatabaseDSN, "database DSN")
	restore := flag.Bool("r", defaultRestore, "read file storage metrics")
	hashKey := flag.String("k", defaultHashKey, "key server encryp/decrypt")
	auditFile := flag.String("audit-file", defaultAuditFile, "audit file log")
	auditURL := flag.String("audit-url", defaultAuditURL, "audit url")
	cryptoKey := flag.String("crypto-key", defaultCryptoKey, "key agent encryption")
	fileConfig := flag.String("c", "", "config file")
	trustedSubnet := flag.String("t", "", "trusted subnet")
	flag.Parse()

	varCryptoKey, ok := os.LookupEnv("CRYPTO_KEY")
	if ok {
		cryptoKey = &varCryptoKey
	}

	varAdrHost, ok := os.LookupEnv("ADDRESS")
	if ok {
		adrHost = &varAdrHost
	}
	varAdrHostGrpc, ok := os.LookupEnv("ADDRESS_GRPC")
	if ok {
		adrHostGrpc = &varAdrHostGrpc
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
	varHashKey, ok := os.LookupEnv("KEY")
	if ok {
		hashKey = &varHashKey
	}

	varAuditFile, ok := os.LookupEnv("AUDIT_FILE")
	if ok {
		auditFile = &varAuditFile
	}

	varAuditURL, ok := os.LookupEnv("AUDIT_URL")
	if ok {
		auditURL = &varAuditURL
	}
	varTrustedSubnet, ok := os.LookupEnv("TRUSTED_SUBNET")
	if ok {
		trustedSubnet = &varTrustedSubnet
	}

	resultServerConfig := ServerConfig{
		AdrHost:         *adrHost,
		AdrHostGrpc:     *adrHostGrpc,
		StoreInterval:   *storeInterval,
		FileStoragePath: *fileStoragePath,
		Restore:         *restore,
		CryptoKey:       *cryptoKey,
		DatabaseDSN:     *databaseDSN,
		HashKey:         *hashKey,
		AuditFile:       *auditFile,
		AuditURL:        *auditURL,
		TrustedSubnet:   *trustedSubnet,
	}

	if *fileConfig != "" {
		fileServerConfig, err := readConfigFile(*fileConfig)
		if err != nil {
			return ServerConfig{}, err
		}
		resultServerConfig.convertFlagsResult(fileServerConfig)
	}

	return resultServerConfig, nil
}

func readConfigFile(filePath string) (ServerConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("could not open config file: %w", err)
	}

	var flags ServerConfig
	err = json.Unmarshal(data, &flags)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("could not parse config file: %w", err)
	}

	return flags, nil
}

func (flags *ServerConfig) convertFlagsResult(flagsInFile ServerConfig) {
	if flags.AdrHost == defaultAdrHost {
		flags.AdrHost = flagsInFile.AdrHost
	}
	if flags.AdrHostGrpc == defaultAdrHostGrpc {
		flags.AdrHostGrpc = flagsInFile.AdrHostGrpc
	}
	if flags.StoreInterval == defaultStoreInterval {
		flags.StoreInterval = flagsInFile.StoreInterval
	}
	if flags.FileStoragePath == defaultFileStoragePath {
		flags.FileStoragePath = flagsInFile.FileStoragePath
	}
	if flags.AuditFile == defaultAuditFile {
		flags.AuditFile = flagsInFile.AuditFile
	}
	if flags.AuditURL == defaultAuditURL {
		flags.AuditURL = flagsInFile.AuditURL
	}
	if flags.CryptoKey == defaultCryptoKey {
		flags.CryptoKey = flagsInFile.CryptoKey
	}
	if flags.DatabaseDSN == defaultDatabaseDSN {
		flags.DatabaseDSN = flagsInFile.DatabaseDSN
	}
	if flags.HashKey == defaultHashKey {
		flags.HashKey = flagsInFile.HashKey
	}
	if flags.Restore == defaultRestore {
		flags.Restore = flagsInFile.Restore
	}
	if flags.TrustedSubnet == defaultTrustedSubnet {
		flags.TrustedSubnet = flagsInFile.TrustedSubnet
	}
}
