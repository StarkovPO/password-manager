package store

import "github.com/jmoiron/sqlx"

type Store struct {
	store *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{store: db}
}
