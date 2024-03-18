package httpserver

import (
	"context"
	"net/http"

	"vk_test/internal/app"
	"vk_test/internal/handlers/actors"
	"vk_test/internal/handlers/films"

	"github.com/sirupsen/logrus"

	_ "vk_test/docs" // без этого работать не будет - нужен путь именно к вашей документации

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type HTTPService struct {
	srv        http.Server
	rolesStore []User
}

type User struct {
	Name     string
	Password string
	Role     string
}

func NewServer(cfg app.Config, mux *http.ServeMux) *HTTPService {
	rStore := make([]User, 0)
	rStore = append(rStore, User{Name: "alex", Password: "12345", Role: "admin"})
	rStore = append(rStore, User{Name: "dave", Password: "1200005", Role: "user"})
	rStore = append(rStore, User{Name: "kirril", Password: "23049", Role: "user"})

	return &HTTPService{
		srv: http.Server{
			Addr:         cfg.HTTPAddress,
			Handler:      mux,
			ReadTimeout:  cfg.HTTPReadTimeout,
			WriteTimeout: cfg.HTTPWriteTimeout,
			IdleTimeout:  cfg.HTTPIdleTimeout,
		},
		rolesStore: rStore,
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
	mux.HandleFunc("/actors", reqIDMiddleware(h.authenticateUser(actorHandler.GetActors)))

	mux.HandleFunc("/actors/create", reqIDMiddleware(h.authenticateAdmin(actorHandler.CreateActor)))
	mux.HandleFunc("/actors/{id}/delete", reqIDMiddleware(h.authenticateAdmin(actorHandler.DeleteActor)))
	mux.HandleFunc("/actors/{id}/update", reqIDMiddleware(h.authenticateAdmin(actorHandler.UpdateActor)))

	// ---FILMS---
	mux.HandleFunc("/films", reqIDMiddleware(h.authenticateUser(filmsHandler.GetFilms)))
	mux.HandleFunc("/films/search", reqIDMiddleware(h.authenticateUser(filmsHandler.SearchFilms)))

	mux.HandleFunc("/films/create", reqIDMiddleware(h.authenticateAdmin(filmsHandler.CreateFilm)))
	mux.HandleFunc("/films/{id}/delete", reqIDMiddleware(h.authenticateAdmin(filmsHandler.DeleteFilm)))
	mux.HandleFunc("/films/{id}/update", reqIDMiddleware(h.authenticateAdmin(filmsHandler.UpdateFilm)))

	mux.HandleFunc("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))
}
