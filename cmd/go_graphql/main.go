package main

import (
	"flag"

	"github.com/MagicRodri/go_graphql_service/internal/api"
	"github.com/MagicRodri/go_graphql_service/internal/config"
	"github.com/MagicRodri/go_graphql_service/internal/logging"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config.Init(*configPath)
	logging.Init()
	logging.Get().Infof("Config loaded: %s", config.Get())
	// var err error
	// err = db.Init(config.Get().DB.DSN)
	// if err != nil {
	// 	logging.Get().Fatal("Database connection failed:", err)
	// }
	// defer db.Close()

	api.SetupServer(config.Get().HTTP.Host, config.Get().HTTP.Port)
}
