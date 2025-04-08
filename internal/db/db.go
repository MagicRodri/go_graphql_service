package db

import (
	"database/sql"

	"github.com/MagicRodri/go_graphql_service/internal/config"
	"github.com/MagicRodri/go_graphql_service/internal/logging"
)

func setupDB() (*sql.DB, error) {
	connection, err := sql.Open("pgx", config.Get().DB.DSN)
	if err != nil {
		return connection, err
	}

	if err := connection.Ping(); err != nil {
		return connection, err
	}

	return connection, nil
}
func GetDB() *sql.DB {
	connection, err := setupDB()
	if err != nil {
		logging.Get().Fatal("Database connection failed:", err)
	}
	return connection
}

// func Close() error {
// 	if connection != nil {
// 		return connection.Close()
// 	}
// 	return nil
// }
