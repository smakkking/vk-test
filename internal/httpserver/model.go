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
	srv http.Server
}

func NewServer(cfg app.Config, mux *http.ServeMux) *HTTPService {
	return &HTTPService{
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
	mux.HandleFunc("/actors", actorHandler.GetActors)

	mux.HandleFunc("/actors/create", actorHandler.CreateActor)
	mux.HandleFunc("/actors/{id}/delete", actorHandler.DeleteActor)
	mux.HandleFunc("/actors/{id}/update", actorHandler.UpdateActor)

	// ---FILMS---
	mux.HandleFunc("/films", filmsHandler.GetFilms)
	mux.HandleFunc("/films/search", filmsHandler.SearchFilms)

	mux.HandleFunc("/films/create", filmsHandler.CreateFilm)
	mux.HandleFunc("/films/{id}/delete", filmsHandler.DeleteFilm)
	mux.HandleFunc("/films/{id}/update", filmsHandler.UpdateFilm)
}
