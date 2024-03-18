package actors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"vk_test/internal/handlers"
	"vk_test/internal/model"
	"vk_test/internal/services/actors"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	serviceActors *actors.Service
}

var (
	ErrIncorrectID  = errors.New("некорректный идентификатор")
	ErrInvalidData  = errors.New("неправильные данные: несовпадение по параметрам")
	ErrServiceError = errors.New("при обработке произошла ошибка")
)

func NewHandler(service *actors.Service) *Handler {
	return &Handler{
		serviceActors: service,
	}
}

// @Summary UpdateActor
// @Security BasicAuth
// @Tags actors
// @Description обновление информации об актере
// @Accept json
// @Param actor_id path int true "идентификатор актера"
// @Param input body model.Actor true "данные актера"
// @Success 200
// @Failure 400
// @Failure 405
// @Failure 500
// @Router /actors/{actor_id}/update [put]
func (h *Handler) UpdateActor(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPut {
		var err error

		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			logrus.Errorf("handlers.UpdateActor: %v", err.Error())
			http.Error(w, ErrIncorrectID.Error(), http.StatusBadRequest)
			return
		}

		oldActor := new(model.ActorPartialUpdate)
		err = json.NewDecoder(req.Body).Decode(&oldActor)
		if err != nil {
			logrus.Errorf("handlers.UpdateActor: %v", err.Error())
			http.Error(w, ErrInvalidData.Error(), http.StatusBadRequest)
			return
		}

		oldActor.NameBool = oldActor.Name != ""
		oldActor.DateBirthBool = !oldActor.DateBirth.IsZero()
		oldActor.SexBool = oldActor.Sex != ""

		err = h.serviceActors.UpdateActor(id, oldActor)
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

// @Summary DeleteActor
// @Security BasicAuth
// @Tags actors
// @Description удаление актера
// @Param actor_id path int true "идентификатор актера"
// @Success 200
// @Failure 400
// @Failure 405
// @Failure 500
// @Router /actors/{actor_id}/delete [delete]
func (h *Handler) DeleteActor(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete {
		var err error

		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			logrus.Errorf("handlers.DeleteActor: %v", err.Error())
			http.Error(w, ErrIncorrectID.Error(), http.StatusBadRequest)
			return
		}

		err = h.serviceActors.DeleteActor(id)
		if err != nil {
			http.Error(w, ErrServiceError.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		logrus.Errorf("expect method DELETE at %s, got %v", req.URL.Path, req.Method)
	}
}

type CreateActorResponce struct {
	ID int `json:"id"`
}

// @Summary CreateActor
// @Security BasicAuth
// @Tags actors
// @Description создание актера
// @Accept json
// @Produce json
// @Param input body model.Actor true "данные актера"
// @Success 200 {integer} integer 1
// @Failure 400
// @Failure 405
// @Failure 500
// @Router /actors/create [post]
func (h *Handler) CreateActor(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		actor := new(model.Actor)
		err := json.NewDecoder(req.Body).Decode(&actor)
		if err != nil {
			logrus.Errorf("handlers.CreateActor: %v", err.Error())
			http.Error(w, ErrInvalidData.Error(), http.StatusBadRequest)
			return
		}

		actor_id, err := h.serviceActors.CreateActor(actor)
		if err != nil {
			http.Error(w, ErrServiceError.Error(), http.StatusInternalServerError)
			return
		}

		handlers.SendJSON(w, CreateActorResponce{ID: actor_id})
		w.WriteHeader(http.StatusOK)

	} else {
		logrus.Errorf("expect method POST at %s, got %v", req.URL.Path, req.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// @Summary GetActors
// @Tags actors
// @Description получить список всех актеров, с фильмами где они снимались
// @Produce json
// @Success 200 {array} model.ActorWithFilms
// @Failure 400
// @Failure 405
// @Failure 500
// @Router /actors [get]
func (h *Handler) GetActors(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		actors, err := h.serviceActors.GetActors()
		if err != nil {
			http.Error(w, fmt.Sprintf("error while getting actor %v", err.Error()), http.StatusInternalServerError)
			return
		}

		handlers.SendJSON(w, actors)
		w.WriteHeader(http.StatusOK)
	} else {
		logrus.Errorf("expect method GET at %s, got %v", req.URL.Path, req.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
