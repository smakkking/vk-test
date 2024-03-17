package httpserver

import (
	"context"
	"log"
	"net/http"

	"vk_test/internal/app"
	"vk_test/internal/handlers/actors"
	"vk_test/internal/handlers/films"
)

type HTTPService struct {
	srv    http.Server
	logger *log.Logger
}

func NewServer(cfg app.Config, logger *log.Logger, mux *http.ServeMux) *HTTPService {
	return &HTTPService{
		logger: logger,
		srv: http.Server{
			Addr:         cfg.HTTPAddress,
			Handler:      mux,
			ReadTimeout:  cfg.HTTPReadTimeout,
			WriteTimeout: cfg.HTTPWriteTimeout,
			IdleTimeout:  cfg.HTTPIdleTimeout,
		},
	}
}

func (h *HTTPService) Run() {
	err := h.srv.ListenAndServe()
	if err != nil {
		log.Fatalf("cannot start server: %s", err)
	}
}

func (h *HTTPService) Shutdown(ctx context.Context) error {
	err := h.srv.Shutdown(ctx)
	return err
}

func (h *HTTPService) SetupHTTPService(
	mux *http.ServeMux,
	actorHandler *actors.Handler,
	filmsHandler *films.Handler,
) {
	// ---ACTORS---
	mux.HandleFunc("/actors/create", actorHandler.CreateActor)
	mux.HandleFunc("/actors/delete", actorHandler.DeleteActor)

	// ---FILMS---
}
