package actors

import (
	"fmt"
	"vk_test/internal/model"
)

type Service struct {
	actorStorage Storage
}

type Storage interface {
	Create(*model.Actor) error
	Delete(string) error
}

func (s *Service) CreateActor(actor *model.Actor) error {
	err := s.actorStorage.Create(actor)
	if err != nil {
		return fmt.Errorf("can't create actor: %w", err)
	}
	return nil
}

func (s *Service) DeleteActor(actorName string) error {
	err := s.actorStorage.Delete(actorName)
	if err != nil {
		return fmt.Errorf("can't create actor: %w", err)
	}
	return nil
}
