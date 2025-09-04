package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabasePostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func (d *DatabasePostgresConfig) ConnectionPostgresString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.Username, d.Password, d.Database, d.SSLMode)
}

func NewDatabasePostgresConfig() *DatabasePostgresConfig {
	return &DatabasePostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "password",
		Database: "postgres",
		SSLMode:  "disable",
	}
}

func (d *DatabasePostgresConfig) PingPostgres() error {
	db, err := sql.Open("postgres", d.ConnectionPostgresString())
	if err != nil {
		fmt.Println("could not connect to postgres:", err)
		return fmt.Errorf("could not connect to postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		fmt.Println("could not ping postgres:", err)
		return fmt.Errorf("could not ping postgres: %w", err)
	}

	return nil
}
