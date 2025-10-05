// Package db
package db

import (
	"context"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

func OpenSQLite(dsn string) (*sql.DB, func() error, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, nil, err
	}

	// Reasonable defaults for SQLite
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	return db, db.Close, nil
}

func InitSchema(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	const ddl = `CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.ExecContext(ctx, ddl)
	return err
}
