package inmemory

import (
	"vk_test/internal/model"

	"github.com/google/uuid"
)

type Storage struct {
	Actors map[string]*model.Actor
	Films  map[string]*model.Film
}

func (s *Storage) Create(*model.Actor) error {
	return nil
}

func (s *Storage) Update(uuid.UUID) error {
	return nil
}

func (s *Storage) Delete(uuid.UUID) error {
	return nil
}
