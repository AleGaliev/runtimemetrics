package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/AleGaliev/runtimemetrics/internal/config/db"
	"github.com/AleGaliev/runtimemetrics/internal/config/server"
	"github.com/AleGaliev/runtimemetrics/internal/filestore"
	"github.com/AleGaliev/runtimemetrics/internal/handler"
	"github.com/AleGaliev/runtimemetrics/internal/logger"
	"github.com/AleGaliev/runtimemetrics/internal/middleware"
	"github.com/AleGaliev/runtimemetrics/internal/observer"
	pb "github.com/AleGaliev/runtimemetrics/internal/proto"
	"github.com/AleGaliev/runtimemetrics/internal/repository"
	"github.com/AleGaliev/runtimemetrics/internal/service/crypto"
	"github.com/AleGaliev/runtimemetrics/internal/service/retry"
	"github.com/AleGaliev/runtimemetrics/internal/storage"
	"google.golang.org/grpc"
)

type ServerMemStorage struct {
	MemStorage handler.Storage
	DBConfig   db.PostgresDB
}

func NewServerMemStorage(serverConf server.ServerConfig) (ServerMemStorage, error) {
	if serverConf.DatabaseDSN == "" {
		fileStore := filestore.NewFileStore(serverConf.FileStoragePath)

		memStorage, err := storage.CreateStorage(fileStore, serverConf.StoreInterval, serverConf.Restore)
		if err != nil {
			return ServerMemStorage{}, err
		}

		if serverConf.StoreInterval > 0 && serverConf.DatabaseDSN == "" {
			go func() {
				for {
					time.Sleep(time.Duration(serverConf.StoreInterval) * time.Second)
					if err := memStorage.SaveMetricToFile(); err != nil {
						panic(errors.Unwrap(err))
					}
				}
			}()
		}
		defer memStorage.SaveMetricToFile()
		return ServerMemStorage{
			MemStorage: memStorage,
			DBConfig:   db.PostgresDB{},
		}, nil
	}
	dbConfig, err := db.NewPostgresDB(serverConf.DatabaseDSN)
	if err != nil {
		return ServerMemStorage{}, err
	}

	if err = dbConfig.CreateMigration(); err != nil {
		return ServerMemStorage{}, err
	}

	dbMemStorage := storage.NewPostgresDBStorage(dbConfig, retry.CreateRetry())

	return ServerMemStorage{
		MemStorage: dbMemStorage,
		DBConfig:   dbConfig,
	}, nil
}

type Server struct {
	ServerConf server.ServerConfig
	LogServer  logger.Logger
	EventAudit *observer.Event
	MemStorage ServerMemStorage
	CryptoKey  *crypto.Crypto
	Server     *http.Server
	GrpcServer *grpc.Server
}

func New() (Server, error) {
	serverConf, err := server.NewServerConfig()
	if err != nil {
		return Server{}, err
	}

	logServer, err := logger.CreateLogger()
	if err != nil {
		return Server{}, err
	}

	eventAudit := createEventAudit(serverConf.AuditFile, serverConf.AuditURL)

	memStorage, err := NewServerMemStorage(serverConf)
	if err != nil {
		return Server{}, err
	}

	cryptoKey, err := crypto.NewCrypto(serverConf.CryptoKey, "")
	if err != nil {
		return Server{}, err
	}

	r := handler.CreateMyHandler(memStorage.MemStorage, memStorage.DBConfig, logServer, serverConf.HashKey, eventAudit, cryptoKey, serverConf.TrustedSubnet)

	server := &http.Server{
		Addr:    serverConf.AdrHost,
		Handler: r,
	}

	serverGrpc := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.LoggingInterceptor(logServer),
			middleware.IPValidateInterceptor(serverConf.TrustedSubnet),
		),
	)

	pb.RegisterMetricsServer(serverGrpc, handler.NewUserServer(memStorage.MemStorage))

	return Server{
		ServerConf: serverConf,
		LogServer:  logServer,
		EventAudit: eventAudit,
		MemStorage: memStorage,
		CryptoKey:  cryptoKey,
		Server:     server,
		GrpcServer: serverGrpc,
	}, nil

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

func (s *Server) Close() {
	s.MemStorage.DBConfig.Close()
}

func (s *Server) StartGrpcServer() error {
	s.LogServer.StartServerLog("grpc", s.ServerConf.AdrHostGrpc)
	listen, err := net.Listen("tcp", s.ServerConf.AdrHostGrpc)
	if err != nil {
		return err
	}
	if err = s.GrpcServer.Serve(listen); err != nil {
		return err
	}
	return nil
}

func (s *Server) StartHTTPServer() error {
	s.LogServer.StartServerLog("http", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.Server.Shutdown(ctx)
	s.GrpcServer.GracefulStop()

	return err
}
