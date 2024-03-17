package actors

import (
	"fmt"
	"vk_test/internal/model"
)

type Service struct {
	actorStorage Storage
}

type Storage interface {
	CreateActor(*model.Actor) (int, error)
	DeleteActor(int) error
	UpdateActor(int, *model.Actor) error
	GetActorsWithFilms() ([]*model.ActorWithFilms, error)
}

func (s *Service) CreateActor(actor *model.Actor) (int, error) {
	id, err := s.actorStorage.CreateActor(actor)
	if err != nil {
		return id, fmt.Errorf("can't create actor: %w", err)
	}
	return 0, nil
}

func (s *Service) DeleteActor(actorID int) error {
	err := s.actorStorage.DeleteActor(actorID)
	if err != nil {
		return fmt.Errorf("can't delete actor: %w", err)
	}

	return nil
}

func (s *Service) UpdateActor(id int, newActor *model.Actor) error {
	err := s.actorStorage.UpdateActor(id, newActor)
	if err != nil {
		return fmt.Errorf("can't update actor: %w", err)
	}

	return nil
}

func (s *Service) GetActors() ([]*model.ActorWithFilms, error) {
	result, err := s.actorStorage.GetActorsWithFilms()
	if err != nil {
		return nil, fmt.Errorf("can't get actors: %w", err)
	}

	return result, nil
}
