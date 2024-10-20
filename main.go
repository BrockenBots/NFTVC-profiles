package main

import (
	"nftvc-profile/pkg/config"
	"nftvc-profile/pkg/logger"
	"nftvc-profile/pkg/server"
)

const pathToConfig = "pkg/config/config.yml"

func main() {
	cfg, err := config.LoadConfig(pathToConfig)
	if err != nil {
		panic("Failed to load config")
	}
	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.Fatal(server.NewServer(appLogger, cfg).Run())
}
