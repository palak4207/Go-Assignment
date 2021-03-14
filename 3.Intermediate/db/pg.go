package db

import (
	"github.com/jackc/pgx/v4"
)

// PGStore ...
type PGStore struct {
	db *pgx.Conn
}

// NewPGStore creates PGStore that implements the 'Store' interface
func NewPGStore(db *pgx.Conn) Store {
	return &PGStore{
		db: db,
	}
}
