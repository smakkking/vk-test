package films

import (
	"vk_test/internal/model"

	"github.com/sirupsen/logrus"
)

type Service struct {
	filmsStorage Storage
}

type Storage interface {
	CreateFilm(*model.Film) (int, error)
	DeleteFilm(int) error
	UpdateFilm(int, *model.FilmPartialUpdate) error

	GetFilmsSorted(sortKey string) ([]*model.FilmWithActors, error)
	SearchFilmsByActorName(actorName string) ([]*model.FilmWithActors, error)
	SearchFilmByTitle(filmTitle string) ([]*model.FilmWithActors, error)
}

func (s *Service) CreateFilm(film *model.Film) (int, error) {
	id, err := s.filmsStorage.CreateFilm(film)
	if err != nil {
		logrus.Errorf("can't create film: %v", err)
		return 0, err
	}

	return id, nil
}

func (s *Service) DeleteFilm(filmID int) error {
	err := s.filmsStorage.DeleteFilm(filmID)
	if err != nil {
		logrus.Errorf("can't delete film: %v", err)
		return err
	}

	return nil
}

func (s *Service) UpdateFilm(filmID int, film *model.FilmPartialUpdate) error {
	err := s.filmsStorage.UpdateFilm(filmID, film)
	if err != nil {
		logrus.Errorf("can't update film: %v", err)
		return err
	}

	return nil
}

func (s *Service) GetFilmsSorted(sortKey string) ([]*model.FilmWithActors, error) {
	films, err := s.filmsStorage.GetFilmsSorted(sortKey)
	if err != nil {
		logrus.Errorf("can't dget film sorted: %v", err)
		return nil, err
	}

	return films, nil
}

func (s *Service) SearchFilmsByActorName(actorName string) ([]*model.FilmWithActors, error) {
	films, err := s.filmsStorage.SearchFilmsByActorName(actorName)
	if err != nil {
		logrus.Errorf("can't get films by actor nam: %v", err)
		return nil, err
	}

	return films, nil
}

func (s *Service) SearchFilmByTitle(filmTitle string) ([]*model.FilmWithActors, error) {
	films, err := s.filmsStorage.SearchFilmByTitle(filmTitle)
	if err != nil {
		logrus.Errorf("can't get films by actor name: %v", err)
		return nil, err
	}

	return films, nil
}
