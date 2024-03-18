package teststore

import (
	"time"
	"vk_test/internal/model"
)

type Storage struct {
	actors_count int
	films_count  int

	Actors map[int]*model.Actor
	Films  map[int]*model.Film
}

func NewStorage() *Storage {
	return &Storage{
		Actors: map[int]*model.Actor{
			1: &model.Actor{
				Name:      "Leo Dicaprio",
				Sex:       "мужчина",
				DateBirth: model.CivilTime(time.Date(2006, 1, 23, 0, 0, 0, 0, time.UTC)),
			},
			2: &model.Actor{
				Name:      "Brad Pit",
				Sex:       "мужчина",
				DateBirth: model.CivilTime(time.Date(1996, 4, 14, 0, 0, 0, 0, time.UTC)),
			},
			3: &model.Actor{
				Name:      "Monica Levinsky",
				Sex:       "женщина",
				DateBirth: model.CivilTime(time.Date(1970, 6, 14, 0, 0, 0, 0, time.UTC)),
			},
		},
		Films: map[int]*model.Film{
			1: &model.Film{
				Title:        "Однажды в Голливуде",
				Description:  "",
				DateCreation: model.CivilTime(time.Date(1970, 6, 14, 0, 0, 0, 0, time.UTC)),
				Rating:       6,
				ActorIDList:  []int{1, 2},
			},
			2: &model.Film{
				Title:        "Семья из Бразилии",
				Description:  "",
				DateCreation: model.CivilTime(time.Date(1975, 7, 24, 0, 0, 0, 0, time.UTC)),
				Rating:       1,
				ActorIDList:  []int{1, 2},
			},
			3: &model.Film{
				Title:        "Убить Билла",
				Description:  "",
				DateCreation: model.CivilTime(time.Date(2000, 10, 2, 0, 0, 0, 0, time.UTC)),
				Rating:       10,
				ActorIDList:  []int{1, 2},
			},
		},
	}
}

func (s *Storage) CreateActor(actor *model.Actor) (int, error) {
	s.actors_count++
	s.Actors[s.actors_count] = actor
	return s.actors_count, nil
}

func (s *Storage) DeleteActor(actorID int) error {
	delete(s.Actors, actorID)
	return nil
}

func (s *Storage) UpdateActor(actorID int, actor *model.ActorPartialUpdate) error {
	if actor.NameBool {
		s.Actors[actorID].Name = actor.Name
	}
	if actor.SexBool {
		s.Actors[actorID].Sex = actor.Sex
	}
	if actor.DateBirthBool {
		s.Actors[actorID].DateBirth = actor.DateBirth
	}
	return nil
}

func (s *Storage) CreateFilm(film *model.Film) (int, error) {
	s.films_count++
	s.Films[s.films_count] = film
	return s.films_count, nil
}

func (s *Storage) DeleteFilm(filmID int) error {
	delete(s.Films, filmID)
	return nil
}

func (s *Storage) UpdateFilm(filmID int, film *model.FilmPartialUpdate) error {
	if film.TitleBool {
		s.Films[filmID].Title = film.Title
	}
	if film.DateCreationBool {
		s.Films[filmID].DateCreation = film.DateCreation
	}
	if film.ActorIDListBool {
		s.Films[filmID].ActorIDList = film.ActorIDList
	}
	if film.RatingBool {
		s.Films[filmID].Rating = film.Rating
	}
	if film.DescriptionBool {
		s.Films[filmID].Description = film.Description
	}

	return nil
}

func (s *Storage) GetFilmsSorted(sortKey string) ([]*model.FilmWithActors, error) {
	if sortKey == "rating" {
		return []*model.FilmWithActors{
			s.convert(2, s.Films[2]),
			s.convert(1, s.Films[1]),
			s.convert(3, s.Films[3]),
		}, nil
	} else if sortKey == "title" {
		return []*model.FilmWithActors{
			s.convert(1, s.Films[1]),
			s.convert(2, s.Films[2]),
			s.convert(3, s.Films[3]),
		}, nil
	}

	return nil, nil
}

func (s *Storage) SearchFilmsByActorName(actorName string) ([]*model.FilmWithActors, error) {
	return nil, nil
}

func (s *Storage) SearchFilmByTitle(filmTitle string) ([]*model.FilmWithActors, error) {
	return nil, nil
}

func (s *Storage) GetActorsWithFilms() ([]*model.ActorWithFilms, error) {
	return nil, nil
}

func (s *Storage) convert(filmID int, film *model.Film) *model.FilmWithActors {
	actorList := make([]*model.Actor, 0)
	for _, actorID := range s.Films[filmID].ActorIDList {
		actorList = append(actorList, s.Actors[actorID])
	}

	return &model.FilmWithActors{
		Title:        film.Title,
		Description:  film.Description,
		DateCreation: film.DateCreation,
		Rating:       film.Rating,
		ActorList:    actorList,
	}
}
