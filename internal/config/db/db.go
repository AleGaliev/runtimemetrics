package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	connMaxLifetime time.Duration = 30 * time.Minute
	connMaxIdleTime time.Duration = 5 * time.Minute
	maxIdleConns    int           = 5
	maxOpenConns    int           = 25
	connTimeout     time.Duration = 5 * time.Second
	defaultTimeout  time.Duration = 5 * time.Second
	queryMigration                = `
		CREATE TABLE IF NOT EXISTS metrics (
		id VARCHAR(255) NOT NULL,
		mtype VARCHAR(50) NOT NULL,
		delta BIGINT,
		value DOUBLE PRECISION,
		hash VARCHAR(255),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		PRIMARY KEY (id, mtype)
		);
	
		CREATE INDEX IF NOT EXISTS idx_metrics_id ON metrics(id);
		CREATE INDEX IF NOT EXISTS idx_metrics_type ON metrics(mtype);
	`
)

type PostgresDB struct {
	DB             *sql.DB
	DefaultTimeout time.Duration
}

func NewPostgresDB(postgresURL string) (PostgresDB, error) {
	db, err := sql.Open("pgx", postgresURL)
	if err != nil {
		return PostgresDB{}, err
	}
	configureConnectionPool(db)

	return PostgresDB{
		DB:             db,
		DefaultTimeout: defaultTimeout,
	}, nil
}

func configureConnectionPool(db *sql.DB) {
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)
}

func (p PostgresDB) Connect() error {
	if err := p.connectWithTimeout(p.DefaultTimeout); err != nil {
		return fmt.Errorf("connection timeout: %w", err)
	}
	return nil
}

func (p PostgresDB) connectWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := p.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping timeout: %w", err)
	}
	return nil
}

func (p *PostgresDB) CreateMigration() error {
	ctx, cancel := context.WithTimeout(context.Background(), p.DefaultTimeout)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, queryMigration)
	return err
}

func (p *PostgresDB) Close() error {
	return p.DB.Close()
}
