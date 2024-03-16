package actors

import "vk_test/internal/services/actors"

type Handler struct {
	serviceActors *actors.Service
}

func NewHandler(service *actors.Service) *Handler {
	return &Handler{
		serviceActors: service,
	}
}
