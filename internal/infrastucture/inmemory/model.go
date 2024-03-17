package inmemory

import (
	"vk_test/internal/model"
)

type Storage struct {
	Actors map[string]*model.Actor
	Films  map[string]*model.Film
}

func (s *Storage) Create(actor *model.Actor) error {
	s.Actors[actor.Name] = actor
	return nil
}

func (s *Storage) Delete(actorName string) error {
	delete(s.Actors, actorName)
	return nil
}
