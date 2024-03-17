package inmemory

import (
	"vk_test/internal/model"
)

type Storage struct {
	Actors map[string]*model.Actor
	Films  map[string]*model.Film
}

func (s *Storage) Create(*model.Actor) error {
	return nil
}
