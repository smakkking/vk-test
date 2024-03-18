package actors

import (
	"vk_test/internal/model"

	"github.com/sirupsen/logrus"
)

type Service struct {
	actorStorage Storage
}

type Storage interface {
	CreateActor(*model.Actor) (int, error)
	DeleteActor(int) error
	UpdateActor(int, *model.ActorPartialUpdate) error
	GetActorsWithFilms() ([]*model.ActorWithFilms, error)
}

func (s *Service) CreateActor(actor *model.Actor) (int, error) {
	id, err := s.actorStorage.CreateActor(actor)
	if err != nil {
		logrus.Errorf("can't create actor: %v", err)
		return 0, err
	}

	return id, nil
}

func (s *Service) DeleteActor(actorID int) error {
	err := s.actorStorage.DeleteActor(actorID)
	if err != nil {
		logrus.Errorf("can't delete actor: %v", err)
		return err
	}

	return nil
}

func (s *Service) UpdateActor(id int, newActor *model.ActorPartialUpdate) error {
	err := s.actorStorage.UpdateActor(id, newActor)
	if err != nil {
		logrus.Errorf("can't update actor: %v", err)
		return err
	}

	return nil
}

func (s *Service) GetActors() ([]*model.ActorWithFilms, error) {
	result, err := s.actorStorage.GetActorsWithFilms()
	if err != nil {
		logrus.Errorf("can't get actor: %v", err)
		return nil, err
	}

	return result, nil
}
