package actors

import (
	"vk_test/internal/model"

	"github.com/google/uuid"
)

type Service struct {
	actorStorage Storage
}

type Storage interface {
	Create(*model.Actor) (uuid.UUID, error)
}

func (s *Service) CreateActor() (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
