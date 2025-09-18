package db

import "database/sql"

func NewPostgresDB(PostgresURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", PostgresURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
