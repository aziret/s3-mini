package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	migrationsTable := os.Getenv("MIGRATIONS_TABLE")
	storagePath := os.Getenv("STORAGE_PATH")
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")

			return
		}
		panic(err)
	}

	fmt.Println("Migrations applied successfully")
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}
