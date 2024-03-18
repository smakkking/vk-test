package films

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"vk_test/internal/handlers"
	"vk_test/internal/model"
	"vk_test/internal/services/films"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	serviceFilms *films.Service
}

func NewHandler(service *films.Service) *Handler {
	return &Handler{
		serviceFilms: service,
	}
}

var (
	ErrIncorrectID  = errors.New("некорректный идентификатор")
	ErrInvalidData  = errors.New("неправильные данные: несовпадение по параметрам")
	ErrServiceError = errors.New("при обработке произошла ошибка")
)

type CreateFilmResponce struct {
	ID int `json:"id"`
}

// @Summary CreateFilm
// @Security BasicAuth
// @Tags films
// @Description создание фильма
// @Accept json
// @Produce json
// @Param input body model.Film true "данные фильма"
// @Success 200 {integer} integer 1
// @Failure 400
// @Failure 405
// @Failure 500
// @Router /films/create [post]
func (h *Handler) CreateFilm(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		film := new(model.Film)
		err := json.NewDecoder(req.Body).Decode(&film)
		if err != nil {
			logrus.Errorf("handlers.CreateFilm: %v", err.Error())
			http.Error(w, ErrInvalidData.Error(), http.StatusBadRequest)
			return
		}

		film_id, err := h.serviceFilms.CreateFilm(film)
		if err != nil {
			http.Error(w, ErrServiceError.Error(), http.StatusInternalServerError)
			return
		}

		handlers.SendJSON(w, CreateFilmResponce{ID: film_id})
		w.WriteHeader(http.StatusOK)

	} else {
		logrus.Errorf("expect method POST at %s, got %v", req.URL.Path, req.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// @Summary DeleteFilm
// @Security BasicAuth
// @Tags films
// @Description удаление фильма
// @Param film_id path int true "идентификатор фильма"
// @Success 200
// @Failure 400
// @Failure 405
// @Failure 500
// @Router /films/{film_id}/delete [delete]
func (h *Handler) DeleteFilm(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete {
		var err error

		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			logrus.Errorf("handlers.DeleteFilm: %v", err.Error())
			http.Error(w, ErrIncorrectID.Error(), http.StatusBadRequest)
			return
		}

		err = h.serviceFilms.DeleteFilm(id)
		if err != nil {
			http.Error(w, ErrServiceError.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		logrus.Errorf("expect method DELETE at %s, got %v", req.URL.Path, req.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// @Summary UpdateFilm
// @Security BasicAuth
// @Tags films
// @Description обновление информации о фильме
// @Accept json
// @Param film_id path int true "идентификатор фильма"
// @Param input body model.Film true "данные фильма"
// @Success 200
// @Failure 400
// @Failure 405
// @Failure 500
// @Router /films/{film_id}/update [put]
func (h *Handler) UpdateFilm(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPut {
		var err error

		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			logrus.Errorf("handlers.UpdateFilm: %v", err.Error())
			http.Error(w, ErrIncorrectID.Error(), http.StatusBadRequest)
			return
		}

		film := new(model.FilmPartialUpdate)
		err = json.NewDecoder(req.Body).Decode(&film)
		if err != nil {
			logrus.Errorf("handlers.UpdateFilm: %v", err.Error())
			http.Error(w, ErrInvalidData.Error(), http.StatusBadRequest)
			return
		}

		film.TitleBool = film.Title != ""
		film.DescriptionBool = film.Description != ""
		film.DateCreationBool = !film.DateCreation.IsZero()
		film.RatingBool = film.Rating != 0
		film.ActorIDListBool = len(film.ActorIDList) != 0

		err = h.serviceFilms.UpdateFilm(id, film)
		if err != nil {
			http.Error(w, ErrServiceError.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	} else {
		logrus.Errorf("expect method PUT at %s, got %v", req.URL.Path, req.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// @Summary GetFilm
// @Tags films
// @Description получить список всех фильмов
// @Produce json
// @Param sort_key query string false "ключ сортировки результата" Enums(title, date_creation, rating)
// @Success 200 {array} model.FilmWithActors
// @Failure 400
// @Failure 405
// @Failure 500
// @Router /films [get]
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
			http.Error(w, ErrServiceError.Error(), http.StatusInternalServerError)
			return
		}

		handlers.SendJSON(w, films)
		w.WriteHeader(http.StatusOK)

	} else {
		logrus.Errorf("expect method GET at %s, got %v", req.URL.Path, req.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// @Summary SearchFilms
// @Tags films
// @Description найти фильмы по фрагменту названия или имени актера
// @Produce json
// @Param parametr query string true "по какому параметру будет идти поиск" Enums(title, actor_name)
// @Param search-query query string true "фрагмент поиска"
// @Success 200 {array} model.FilmWithActors
// @Failure 400
// @Failure 405
// @Failure 500
// @Router /films/search [get]
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
				http.Error(w, ErrServiceError.Error(), http.StatusInternalServerError)
				return
			}
		case "actor_name":
			result, err = h.serviceFilms.SearchFilmsByActorName(searchPattern)
			if err != nil {
				http.Error(w, ErrServiceError.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "no such search parametr", http.StatusBadRequest)
		}

		handlers.SendJSON(w, result)
		w.WriteHeader(http.StatusOK)

	} else {
		logrus.Errorf("expect method GET at %s, got %v", req.URL.Path, req.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
