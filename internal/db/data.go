package db

import "database/sql"

// Repository contains methods for interacting with a database storage
type Data struct {
	*Queries
	db *sql.DB
}

// NewRepository returns an implementation of the Repository interface.
func NewData(db *sql.DB) *Data {
	return &Data{
		Queries: New(db),
		db:      db,
	}
}
