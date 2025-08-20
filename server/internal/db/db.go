package db

import (
	"database/sql"
)

type DB struct {
	conn *sql.DB
}

const connectionString = "postgresql://root:pass@localhost:5433/go-chat?sslmode=disable"

func NewDB() (*DB, error) {
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	if err := db.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (db *DB) GetDB() *sql.DB {
	return db.conn
}
