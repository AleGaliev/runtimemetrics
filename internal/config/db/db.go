package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DatabasePostgresConfig struct {
	PostgresURL string
}

func NewDatabasePostgresConfig(postgresURL string) *DatabasePostgresConfig {
	return &DatabasePostgresConfig{
		PostgresURL: postgresURL,
	}
}

func (d *DatabasePostgresConfig) PingPostgres() error {
	db, err := sql.Open("pgx", d.PostgresURL)
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
