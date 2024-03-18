package main

import (
	"net/http"
	"os"
	"vk_test/internal/app"
	handlerActors "vk_test/internal/handlers/actors"
	handlerFilms "vk_test/internal/handlers/films"
	"vk_test/internal/httpserver"
	"vk_test/internal/infrastucture/postgres"
	serviceActors "vk_test/internal/services/actors"
	serviceFilms "vk_test/internal/services/films"

	"github.com/sirupsen/logrus"
)

const (
	configPath = "./config/config.yaml"
)

func main() {
	// logger setup
	logrus.SetFormatter(
		&logrus.JSONFormatter{
			PrettyPrint:     true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
	)
	logrus.SetOutput(os.Stdout)
	logrus.Infoln("service started...")

	// init config
	config, err := app.NewConfig(configPath)
	if err != nil {
		logrus.Errorf("error reading config: %s", err.Error())
		os.Exit(1)
	}

	// init db
	mainStorage, err := postgres.NewStorage(config)
	if err != nil {
		logrus.Errorf("error connecting to db: %s", err.Error())
		os.Exit(1)
	}
	logrus.Infoln("connected to DB...")

	// init services
	actorService := serviceActors.NewService(mainStorage)
	filmService := serviceFilms.NewService(mainStorage)
	logrus.Infoln("services initialized...")

	// init handlers
	actorHandler := handlerActors.NewHandler(actorService)
	filmsHandler := handlerFilms.NewHandler(filmService)

	mux := http.NewServeMux()
	logrus.Infoln("handlers initialized...")

	// start server
	srv := httpserver.NewServer(config, mux)
	srv.SetupHTTPService(mux, actorHandler, filmsHandler)

	srv.Run()
}
