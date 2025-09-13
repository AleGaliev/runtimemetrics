package postgresdb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	models "github.com/AleGaliev/kubercontroller/internal/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	queryMigration = `
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
	queryUpgrad = `
		INSERT INTO metrics (id, mtype, delta, value, hash, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW())
        ON CONFLICT (id, mtype) 
        DO UPDATE SET 
            delta = EXCLUDED.delta,
            value = EXCLUDED.value,
            hash = EXCLUDED.hash,
            updated_at = NOW()
        WHERE metrics.id = EXCLUDED.id AND metrics.mtype = EXCLUDED.mtype
		`
	queryGet = `
		SELECT id, mtype, delta, value, hash
        FROM metrics 
        WHERE id = $1 AND mtype = $2`
	queryGetСontent = `
		SELECT id
        FROM metrics 
        WHERE id = $1`
	queryGetAll = `SELECT id, mtype, delta, value, hash FROM metrics`
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(PostgresURL string) (*PostgresDB, error) {
	db, err := sql.Open("pgx", PostgresURL)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) Connect() error {
	if err := p.db.Ping(); err != nil {
		return fmt.Errorf("could not ping postgres: %w", err)
	}
	return nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

func (p *PostgresDB) Migrate() error {
	driver, err := postgres.WithInstance(p.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}

func (p *PostgresDB) AddMetric(myType, name, value string) error {
	metrics := models.Metrics{
		ID:    name,
		MType: myType,
	}
	switch myType {
	case models.Gauge:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		metrics.Value = &f
	case models.Counter:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		var deltaOld int64
		err = p.db.QueryRow(queryGet, name, metrics.MType).Scan(&deltaOld)

		if errors.Is(err, sql.ErrNoRows) {
			metrics.Delta = &i
		} else if err != nil {
			return err
		} else {
			i += deltaOld
			metrics.Delta = &i
		}
	default:
		return fmt.Errorf("unknown metric type: %s", myType)
	}

	_, err := p.db.Exec(queryUpgrad, metrics.ID, metrics.MType, metrics.Delta, metrics.Value, metrics.Hash)
	if err != nil {
		return fmt.Errorf("failed to add metric: %w", err)
	}
	return nil
}

func (p *PostgresDB) GetMetrics(name string) (string, bool) {
	check := false
	var content string
	err := p.db.QueryRow(queryGetСontent, name).Scan(&content)
	if errors.Is(err, sql.ErrNoRows) {
		check = true
	}
	return content, check
}

func (p *PostgresDB) GetAllMetric() (string, error) {
	result := ""
	rows, err := p.db.Query(queryGetAll)
	if err != nil {
		return "", fmt.Errorf("failed to query all metrics: %w", err)
	}

	for rows.Next() {
		metric := models.Metrics{}
		err := rows.Scan(
			&metric.ID,
			&metric.MType,
			&metric.Delta,
			&metric.Value,
			&metric.Hash,
		)
		if err != nil {
			return "", fmt.Errorf("failed to scan metrics: %w", err)
		}
		switch metric.MType {
		case models.Gauge:
			result += fmt.Sprintf("<li> %s: %g</li>", metric.ID, *metric.Value)
		case models.Counter:
			result += fmt.Sprintf("<li> %s: %d</li>", metric.ID, *metric.Delta)
		}
	}

	if err = rows.Err(); err != nil {
		return "", fmt.Errorf("error during rows iteration: %w", err)
	}

	return result, nil
}
func (p *PostgresDB) UpdateMetrics(r io.Reader) error {
	data := json.NewDecoder(r)
	var metricsData models.Metrics
	if err := data.Decode(&metricsData); err != nil {
		return fmt.Errorf("could not decode metrics: %v", err)
	}

	switch metricsData.MType {

	case models.Gauge:
		if metricsData.Value == nil {
			return fmt.Errorf("metrics value is nil")
		}

	case models.Counter:

		if metricsData.Delta == nil {
			return fmt.Errorf("metrics delta is nil")
		}

		var metricsOld models.Metrics
		err := p.db.QueryRow(queryGet, metricsData.ID, metricsData.MType).Scan(&metricsOld.ID, &metricsOld.MType, &metricsOld.Delta, &metricsOld.Value, &metricsOld.Hash)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("could not check existing metrics: %w", err)
		}

		if !errors.Is(err, sql.ErrNoRows) {
			*metricsData.Delta += *metricsOld.Delta

		}

	default:
		return fmt.Errorf("unknown metric type: %s", metricsData.MType)
	}

	_, err := p.db.Exec(queryUpgrad, metricsData.ID, metricsData.MType, metricsData.Delta, metricsData.Value, metricsData.Hash)
	if err != nil {
		return fmt.Errorf("failed to add metric: %w", err)
	}
	return nil

}
func (p *PostgresDB) ValueMetrics(r io.Reader) ([]byte, bool, error) {
	data := json.NewDecoder(r)
	var metrics models.Metrics
	if err := data.Decode(&metrics); err != nil {
		return nil, false, fmt.Errorf("could not decode metrics: %v", err)
	}
	if (metrics.MType != models.Counter && metrics.MType != models.Gauge) || metrics.ID == "" {
		return nil, false, fmt.Errorf("invalid metric type: %s", metrics.MType)
	}
	if metrics.Value != nil || metrics.Delta != nil {
		return nil, false, fmt.Errorf("invalid metric type: %s", metrics.MType)
	}

	err := p.db.QueryRow(queryGet, metrics.ID, metrics.MType).
		Scan(&metrics.ID, &metrics.MType, &metrics.Delta, &metrics.Value, &metrics.Hash)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, fmt.Errorf("failed to query metrics: %w", err)
	}

	resp, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return nil, false, fmt.Errorf("could not encode metrics: %v", err)
	}
	return resp, true, nil

}

func (p *PostgresDB) CreateMigration() error {
	_, err := p.db.Exec(queryMigration)
	return err
}
