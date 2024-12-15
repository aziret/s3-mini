package file

import "database/sql"

type repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		DB: db,
	}
}
