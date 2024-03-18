package httpserver

import (
	"context"
	"net/http"

	"vk_test/internal/app"
	"vk_test/internal/handlers/actors"
	"vk_test/internal/handlers/films"

	"github.com/sirupsen/logrus"
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
		logrus.Errorf("cannot start server: %s", err.Error())
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
	mux.HandleFunc("/actors", reqIDMiddleware(actorHandler.GetActors))

	mux.HandleFunc("/actors/create", reqIDMiddleware(actorHandler.CreateActor))
	mux.HandleFunc("/actors/{id}/delete", reqIDMiddleware(actorHandler.DeleteActor))
	mux.HandleFunc("/actors/{id}/update", reqIDMiddleware(actorHandler.UpdateActor))

	// ---FILMS---
	mux.HandleFunc("/films", reqIDMiddleware(filmsHandler.GetFilms))
	mux.HandleFunc("/films/search", reqIDMiddleware(filmsHandler.SearchFilms))

	mux.HandleFunc("/films/create", reqIDMiddleware(filmsHandler.CreateFilm))
	mux.HandleFunc("/films/{id}/delete", reqIDMiddleware(filmsHandler.DeleteFilm))
	mux.HandleFunc("/films/{id}/update", reqIDMiddleware(filmsHandler.UpdateFilm))
}
