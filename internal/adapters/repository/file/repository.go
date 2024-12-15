package file

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type repository struct {
	DB *sql.DB
}

func NewRepository(storagePath string) (*repository, error) {
	const op = "repository.NewRepository"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &repository{
		DB: db,
	}, nil
}
