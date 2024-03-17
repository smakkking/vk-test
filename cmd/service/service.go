package main

import (
	"log"
	"net/http"
	"os"
	"vk_test/internal/app"
	handlerActors "vk_test/internal/handlers/actors"
	handlerFilms "vk_test/internal/handlers/films"
	"vk_test/internal/httpserver"
	"vk_test/internal/infrastucture/inmemory"
	serviceActors "vk_test/internal/services/actors"
	serviceFilms "vk_test/internal/services/films"
)

const (
	configPath = "./config/config.yaml"
)

func main() {
	// init logger
	logger := log.Default()
	logger.Println("service started ...")

	// init config
	config, err := app.NewConfig(configPath)
	if err != nil {
		logger.Fatalf("error reading config: %s", err.Error())
		os.Exit(1)
	}

	// init db
	mainStorage, err := inmemory.NewStorage()
	if err != nil {
		logger.Fatalf("error connecting to db: %s", err.Error())
		os.Exit(1)
	}

	// init services
	actorService := serviceActors.NewService(mainStorage)
	filmService := serviceFilms.NewService(mainStorage)

	// init handlers
	actorHandler := handlerActors.NewHandler(actorService)
	filmsHandler := handlerFilms.NewHandler(filmService)

	mux := http.NewServeMux()

	// start server
	srv := httpserver.NewServer(config, logger, mux)
	srv.SetupHTTPService(mux, actorHandler, filmsHandler)

	srv.Run()
}
