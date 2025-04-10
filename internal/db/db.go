package db

import (
	"database/sql"

	"github.com/MagicRodri/go_graphql_service/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var connection *sql.DB

func SetupDB() error {
	db, err := sql.Open("pgx", config.Get().DB.DSN)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	connection = db
	return nil
}

func GetDB() *sql.DB {
	return connection
}
