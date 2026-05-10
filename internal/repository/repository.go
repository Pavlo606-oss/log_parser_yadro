package repository

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(Db *sql.DB) *Repository {
	return &Repository{db: Db}
}
