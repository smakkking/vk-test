package films

import "vk_test/internal/services/films"

type Handler struct {
	serviceFilms *films.Service
}

func NewHandler(service *films.Service) *Handler {
	return &Handler{
		serviceFilms: service,
	}
}
