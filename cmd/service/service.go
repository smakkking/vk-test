package main

import (
	"log"
	"os"
	"vk_test/internal/app"
	"vk_test/internal/infrastucture/inmemory"
	"vk_test/internal/services/actors"
	"vk_test/internal/services/films"
)

const (
	configPath = "./config/config.yaml"
)

func main() {
	// init logger
	logger := log.Default()

	// init config
	config, err := app.NewConfig(configPath)
	if err != nil {
		os.Exit(1)
		logger.Fatalf("error reading config: %s", err.Error())
	}

	// init db
	mainStorage := inmemory.NewStorage()

	// init services
	actorService := actors.NewService(mainStorage)
	filmService := films.NewService(mainStorage)

	// init handlers

	// start server
}
