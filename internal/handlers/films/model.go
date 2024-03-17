package films

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"vk_test/internal/handlers"
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

type CreateFilmResponce struct {
	ID int `json:"id"`
}

func (h *Handler) CreateFilm(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		film := new(model.Film)
		err := json.NewDecoder(req.Body).Decode(&film)
		if err != nil {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		}

		film_id, err := h.serviceFilms.CreateFilm(film)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while creating film %v", err.Error()), http.StatusInternalServerError)
			return
		}

		handlers.SendJSON(w, CreateFilmResponce{ID: film_id}, "error while creating film %v")
		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, fmt.Sprintf("expect method POST at %s, got %v", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) DeleteFilm(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete {
		var err error

		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, fmt.Sprintf("error while deleting film %v", err.Error()), http.StatusInternalServerError)
			return
		}

		err = h.serviceFilms.DeleteFilm(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while deleting film %v", err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, fmt.Sprintf("expect method DELETE at %s, got %v", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) UpdateFilm(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPut {
		var err error

		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, fmt.Sprintf("error while updating film %v", err.Error()), http.StatusInternalServerError)
			return
		}

		film := new(model.FilmPartialUpdate)
		err = json.NewDecoder(req.Body).Decode(&film)
		if err != nil {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		}

		film.TitleBool = film.Title != ""
		film.DescriptionBool = film.Description != ""
		film.DateCreationBool = !film.DateCreation.IsZero()
		film.RatingBool = film.Rating != 0
		film.ActorIDListBool = len(film.ActorIDList) != 0

		err = h.serviceFilms.UpdateFilm(id, film)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while deleting film %v", err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, fmt.Sprintf("expect method PUT at %s, got %v", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) GetFilms(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		sortKey := req.URL.Query().Get("sort_key")
		switch sortKey {
		case "title":
		case "date_creation":
		default:
			sortKey = "rating"
		}

		films, err := h.serviceFilms.GetFilmsSorted(sortKey)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while getting film %v", err.Error()), http.StatusInternalServerError)
			return
		}

		handlers.SendJSON(w, films, "error while getting film %v")
		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, fmt.Sprintf("expect method GET at %s, got %v", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) SearchFilms(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		searchWay := req.URL.Query().Get("parametr")
		searchPattern := req.URL.Query().Get("search-query")

		var result []*model.FilmWithActors
		var err error

		switch searchWay {
		case "title":
			result, err = h.serviceFilms.SearchFilmByTitle(searchPattern)
			if err != nil {
				http.Error(w, fmt.Sprintf("error while search films %v", err.Error()), http.StatusInternalServerError)
				return
			}
		case "actor_name":
			result, err = h.serviceFilms.SearchFilmsByActorName(searchPattern)
			if err != nil {
				http.Error(w, fmt.Sprintf("error while search films %v", err.Error()), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "no such search parametr", http.StatusBadRequest)
		}

		handlers.SendJSON(w, result, "error while search films %v")
		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, fmt.Sprintf("expect method GET at %s, got %v", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
	}
}
