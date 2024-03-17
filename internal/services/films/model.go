package films

import (
	"fmt"
	"vk_test/internal/model"
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
		return 0, fmt.Errorf("can't create film: %w", err)
	}

	return id, nil
}

func (s *Service) DeleteFilm(filmID int) error {
	err := s.filmsStorage.DeleteFilm(filmID)
	if err != nil {
		return fmt.Errorf("can't delete film: %w", err)
	}

	return nil
}

func (s *Service) UpdateFilm(filmID int, film *model.FilmPartialUpdate) error {
	err := s.filmsStorage.UpdateFilm(filmID, film)
	if err != nil {
		return fmt.Errorf("can't update film: %w", err)
	}

	return nil
}

func (s *Service) GetFilmsSorted(sortKey string) ([]*model.FilmWithActors, error) {
	films, err := s.filmsStorage.GetFilmsSorted(sortKey)
	if err != nil {
		return nil, fmt.Errorf("can't get film sorted: %w", err)
	}

	return films, nil
}

func (s *Service) SearchFilmsByActorName(actorName string) ([]*model.FilmWithActors, error) {
	films, err := s.filmsStorage.SearchFilmsByActorName(actorName)
	if err != nil {
		return nil, fmt.Errorf("can't get films by actor name: %w", err)
	}

	return films, nil
}

func (s *Service) SearchFilmByTitle(filmTitle string) ([]*model.FilmWithActors, error) {
	films, err := s.filmsStorage.SearchFilmByTitle(filmTitle)
	if err != nil {
		return nil, fmt.Errorf("can't get films by actor name: %w", err)
	}

	return films, nil
}
