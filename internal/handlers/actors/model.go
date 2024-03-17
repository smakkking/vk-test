package actors

import (
	"fmt"
	"net/http"
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

func (h *Handler) CreateActor(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		_, err := h.serviceActors.CreateActor()
		if err != nil {
			http.Error(w, fmt.Sprintf("expect method POST at %s, got %v", req.URL.Path, req.Method), http.StatusAccepted)
		}
	} else {
		http.Error(w, fmt.Sprintf("expect method POST at %s, got %v", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
	}
}
