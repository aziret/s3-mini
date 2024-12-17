package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
		pgUser         = os.Getenv("PG_USER")
		pgPass         = os.Getenv("PG_PASS")
		pgHost         = os.Getenv("PG_HOST")
		pgPort         = os.Getenv("PG_PORT")
		pgDB           = os.Getenv("PG_DB")
		pgSSLMode      = os.Getenv("PG_SSL_MODE")
	)
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", pgUser, pgPass, pgHost, pgPort, pgDB, pgSSLMode)

	m, err := migrate.New(
		"file://"+migrationsPath,
		databaseURL,
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
