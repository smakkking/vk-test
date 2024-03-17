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
}

func (s *Service) CreateActor(actor *model.Actor) error {
	err := s.actorStorage.Create(actor)
	if err != nil {
		return fmt.Errorf("can't create actor: %w", err)
	}
	return nil
}
