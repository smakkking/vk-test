package actors

import (
	"fmt"
	"net/http"
	"time"
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

func (h *Handler) DeleteActor(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete {
		var err error

		actorName := req.PostFormValue("name")
		err = h.serviceActors.DeleteActor(actorName)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while deleting actor %v", err.Error()), http.StatusAccepted)
			return
		}

	} else {
		http.Error(w, fmt.Sprintf("expect method DELETE at %s, got %v", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) CreateActor(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var err error

		actor := new(model.Actor)
		actor.Name = req.PostFormValue("name")
		actor.Sex = req.PostFormValue("sex")
		actor.DateBirth, err = time.Parse("02-01-2006", req.PostFormValue("date_birth"))
		if err != nil {
			http.Error(w, "invalid date format", http.StatusInternalServerError)
			return
		}

		err = h.serviceActors.CreateActor(actor)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while creating actor %v", err.Error()), http.StatusAccepted)
			return
		}

		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, fmt.Sprintf("expect method POST at %s, got %v", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
	}
}
