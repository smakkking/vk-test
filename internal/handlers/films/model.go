package films

import (
	"fmt"
	"net/http"
	"time"
	"vk_test/internal/model"
	"vk_test/internal/services/films"
)

type Handler struct {
	serviceFilms *films.Service
}

func NewHandler(service *films.Service) *Handler {
	return &Handler{
		serviceFilms: service,
	}
}

func (h *Handler) CreateFilm(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var err error

		film := new(model.Film)
		film.Title = req.PostFormValue("title")
		film.Description = req.PostFormValue("desc")
		film.DateCreation, err = time.Parse("02-01-2006", req.PostFormValue("date_creation"))
		if err != nil {
			http.Error(w, "invalid date format", http.StatusInternalServerError)
			return
		}
		film.ActorIDList

		actor_id, err := h.serviceActors.CreateActor(actor)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while creating actor %v", err.Error()), http.StatusInternalServerError)
			return
		}

		sendJSON(w, CreateActorResponce{ID: actor_id}, "error while creating actor %v")
		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, fmt.Sprintf("expect method POST at %s, got %v", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
	}
}
