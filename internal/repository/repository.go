package repository

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func NewRepository(Db *sql.DB) *Repository {
	return &Repository{db: Db}
}
