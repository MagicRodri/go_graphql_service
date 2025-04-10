package main

import (
	"flag"

	"github.com/MagicRodri/go_graphql_service/internal/api"
	"github.com/MagicRodri/go_graphql_service/internal/config"
	"github.com/MagicRodri/go_graphql_service/internal/db"
	"github.com/MagicRodri/go_graphql_service/internal/logging"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config.Init(*configPath)
	logging.Init()
	logging.Get().Infof("Config loaded: %s", config.Get())
	err := db.SetupDB()
	if err != nil {
		logging.Get().Fatal("Database connection failed:", err)
	}
	logging.Get().Info("Database connection established")

	if err := api.InitGraphQLTables(db.GetDB()); err != nil {
		logging.Get().Fatal("Failed to initialize GraphQL tables:", err)
	}
	api.SetupServer(config.Get().HTTP.Host, config.Get().HTTP.Port)
}
