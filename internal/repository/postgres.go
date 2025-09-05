package repository

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

func NewPostgresDB() (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
		"disable")

	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
