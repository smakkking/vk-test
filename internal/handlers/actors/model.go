package actors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"vk_test/internal/model"
	"vk_test/internal/services/actors"
)

type Handler struct {
	serviceActors *actors.Service
}

func NewHandler(service *actors.Service) *Handler {
	return &Handler{
		serviceActors: service,
	}
}
func (h *Handler) UpdateActor(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPut {
		var err error

		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, fmt.Sprintf("error while updating actor %v", err.Error()), http.StatusInternalServerError)
			return
		}

		oldActor := new(model.ActorPartialUpdate)
		err = json.NewDecoder(req.Body).Decode(&oldActor)
		if err != nil {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		}

		oldActor.NameBool = oldActor.Name != ""
		oldActor.DateBirthBool = !oldActor.DateBirth.IsZero()
		oldActor.SexBool = oldActor.Sex != ""

		err = h.serviceActors.UpdateActor(id, oldActor)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while updating actor %v", err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, fmt.Sprintf("expect method PUT at %s, got %v", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) DeleteActor(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete {
		var err error

		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, fmt.Sprintf("error while deleting actor %v", err.Error()), http.StatusInternalServerError)
			return
		}

		err = h.serviceActors.DeleteActor(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while deleting actor %v", err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, fmt.Sprintf("expect method DELETE at %s, got %v", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
	}
}

type CreateActorResponce struct {
	ID int `json:"id"`
}

func (h *Handler) CreateActor(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		actor := new(model.Actor)
		err := json.NewDecoder(req.Body).Decode(&actor)
		if err != nil {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		}

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

func (h *Handler) GetActors(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		actors, err := h.serviceActors.GetActors()
		if err != nil {
			http.Error(w, fmt.Sprintf("error while getting actor %v", err.Error()), http.StatusInternalServerError)
			return
		}

		sendJSON(w, actors, "error while getting actor %v")
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, fmt.Sprintf("expect method GET at %s, got %v", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
	}
}

func sendJSON(w http.ResponseWriter, payload interface{}, errorMessage string) {
	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf(errorMessage, err.Error()), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, fmt.Sprintf(errorMessage, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
}
