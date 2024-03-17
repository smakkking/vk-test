package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"vk_test/internal/model"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) CreateActor(actor *model.Actor) (int, error) {
	var id int

	err := s.db.QueryRow(
		"INSERT INTO Actors(a_id, a_name, a_sex, a_birth_date) VALUES ($1, $2, $3) RETURNING a_id",
		actor.Name,
		actor.Sex,
		actor.DateBirth,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) DeleteActor(actorID int) error {
	_, err := s.db.Exec("DELETE FROM Actors WHERE a_id = $1", actorID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateActor(actorID int, actor *model.Actor) error {
	b := strings.Builder{}
	b.WriteString("UPDATE Actors SET ")

	c := 1
	params := make([]interface{}, 0)

	if actor.Name != "" {
		b.WriteString(fmt.Sprintf("a_name = $%d", c))
		params = append(params, actor.Name)
		c++
	}

	if actor.Sex != "" {
		b.WriteString(fmt.Sprintf("a_sex = $%d", c))
		params = append(params, actor.Sex)
		c++

	}

	// TODO: костыль
	e := time.Time{}
	if actor.DateBirth != e {
		b.WriteString(fmt.Sprintf("a_birth_date = $%d", c))
		params = append(params, actor.DateBirth)
		c++
	}

	b.WriteString(fmt.Sprintf("WHERE a_id = $%d", c))

	_, err := s.db.Exec(b.String(), params...)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateFilm(film *model.Film) (int, error) {
	var id int

	tX, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	// вставка инфы о фильме
	err = tX.QueryRow(
		"INSERT INTO Films(f_id, f_title, f_desc, f_date_creation, f_rating) VALUES ($1, $2, $3, $4) RETURNING f_id",
		film.Title,
		film.Description,
		film.DateCreation,
		film.Rating,
	).Scan(&id)
	if err != nil {
		_ = tX.Rollback()
		return 0, err
	}

	// вставка инфы об актерах, играющих в фильме
	for _, actorID := range film.ActorIDList {
		_, err := tX.Exec(
			"INSERT INTO ActorToFilm(actor_id, film_id) VALUES($1, $2)",
			actorID,
			id,
		)
		if err != nil {
			_ = tX.Rollback()
			return 0, err
		}
	}

	_ = tX.Commit()
	return id, nil
}

func (s *Storage) DeleteFilm(filmID int) error {
	_, err := s.db.Exec("DELETE FROM Films WHERE f_id = $1", filmID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetFilmsSorted(sortKey string) ([]*model.FilmWithActors, error) {
	b := strings.Builder{}
	b.WriteString("SELECT f.f_id, f.f_title, f.f_desc, f.f_date_creation, f.f_rating, a.a_name, a.a_sex, a.a_birth_date ")
	b.WriteString("FROM ")
	b.WriteString("Films AS f ")
	b.WriteString("JOIN ActorToFilm AS atf ON f.f_id = atf.film_id ")
	b.WriteString("JOIN Actors AS a ON a.a_id = atf.actor_id ")

	switch sortKey {
	case "title":
		b.WriteString("ORDER BY f.f_title, f.f_id ")
	case "rating":
		b.WriteString("ORDER BY f.f_rating, f.f_id ")
	case "date_creation":
		b.WriteString("ORDER BY f.f_date_creation, f.f_id ")
	}

	rows, err := s.db.Query(b.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.FilmWithActors
	lastID := -1

	for rows.Next() {
		var filmID, filmRating int
		var filmTitle, filmDescription, actorName, actorSex string
		var filmDateCreation, actorDateBirth time.Time

		err := rows.Scan(&filmID, &filmTitle, &filmDescription, &filmDateCreation, &filmRating, &actorName, &actorSex, &actorDateBirth)
		if err != nil {
			return nil, err
		}

		// если фильм меняется в списке добавляем новый элемент
		if filmID != lastID || lastID == -1 {
			result = append(result, &model.FilmWithActors{
				Title:        filmTitle,
				Description:  filmDescription,
				DateCreation: filmDateCreation,
				Rating:       filmRating,
				ActorList:    make([]*model.Actor, 0),
			})
		}
		// а актера всегда вставляем в последний
		result[len(result)-1].ActorList = append(
			result[len(result)-1].ActorList,
			&model.Actor{
				Name:      actorName,
				Sex:       actorSex,
				DateBirth: actorDateBirth,
			},
		)

		lastID = filmID
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) SearchFilmsByActorName(actorName string) ([]*model.FilmWithActors, error) {
	b := strings.Builder{}
	b.WriteString("SELECT f.f_id, f.f_title, f.f_desc, f.f_date_creation, f.f_rating, a.a_name, a.a_sex, a.a_birth_date ")
	b.WriteString("FROM ")
	b.WriteString("Films AS f ")
	b.WriteString("JOIN ActorToFilm AS atf ON f.f_id = atf.film_id ")
	b.WriteString("JOIN Actors AS a ON a.a_id = atf.actor_id ")
	b.WriteString("WHERE a.a_name LIKE '%$1%' ")
	b.WriteString("ORDER BY f.f_id ")

	rows, err := s.db.Query(b.String(), actorName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.FilmWithActors
	lastID := -1

	for rows.Next() {
		var filmID, filmRating int
		var filmTitle, filmDescription, actorName, actorSex string
		var filmDateCreation, actorDateBirth time.Time

		err := rows.Scan(&filmID, &filmTitle, &filmDescription, &filmDateCreation, &filmRating, &actorName, &actorSex, &actorDateBirth)
		if err != nil {
			return nil, err
		}

		// если фильм меняется в списке добавляем новый элемент
		if filmID != lastID || lastID == -1 {
			result = append(result, &model.FilmWithActors{
				Title:        filmTitle,
				Description:  filmDescription,
				DateCreation: filmDateCreation,
				Rating:       filmRating,
				ActorList:    make([]*model.Actor, 0),
			})
		}
		// а актера всегда вставляем в последний
		result[len(result)-1].ActorList = append(
			result[len(result)-1].ActorList,
			&model.Actor{
				Name:      actorName,
				Sex:       actorSex,
				DateBirth: actorDateBirth,
			},
		)

		lastID = filmID
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) SearchFilmByTitle(filmTitle string) ([]*model.FilmWithActors, error) {
	b := strings.Builder{}
	b.WriteString("SELECT f.f_id, f.f_title, f.f_desc, f.f_date_creation, f.f_rating, a.a_name, a.a_sex, a.a_birth_date ")
	b.WriteString("FROM ")
	b.WriteString("Films AS f ")
	b.WriteString("JOIN ActorToFilm AS atf ON f.f_id = atf.film_id ")
	b.WriteString("JOIN Actors AS a ON a.a_id = atf.actor_id ")
	b.WriteString("WHERE f.f_title LIKE '%$1%' ")
	b.WriteString("ORDER BY f.f_id ")

	rows, err := s.db.Query(b.String(), filmTitle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.FilmWithActors
	lastID := -1

	for rows.Next() {
		var filmID, filmRating int
		var filmTitle, filmDescription, actorName, actorSex string
		var filmDateCreation, actorDateBirth time.Time

		err := rows.Scan(&filmID, &filmTitle, &filmDescription, &filmDateCreation, &filmRating, &actorName, &actorSex, &actorDateBirth)
		if err != nil {
			return nil, err
		}

		// если фильм меняется в списке добавляем новый элемент
		if filmID != lastID || lastID == -1 {
			result = append(result, &model.FilmWithActors{
				Title:        filmTitle,
				Description:  filmDescription,
				DateCreation: filmDateCreation,
				Rating:       filmRating,
				ActorList:    make([]*model.Actor, 0),
			})
		}
		// а актера всегда вставляем в последний
		result[len(result)-1].ActorList = append(
			result[len(result)-1].ActorList,
			&model.Actor{
				Name:      actorName,
				Sex:       actorSex,
				DateBirth: actorDateBirth,
			},
		)

		lastID = filmID
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) GetActorsWithFilms() ([]*model.ActorWithFilms, error) {
	b := strings.Builder{}
	b.WriteString("SELECT a.a_id, a.a_name, a.a_sex, a.a_birth_date, f.f_title, f.f_desc, f.f_date_creation, f.f_rating ")
	b.WriteString("FROM ")
	b.WriteString("Actors as a ")
	b.WriteString("JOIN ActorToFilm AS atf ON a.a_id = atf.actor_id ")
	b.WriteString("JOIN Films as f ON f.f_id = atf.film_id ")
	b.WriteString("ORDER BY a.a_id ")

	rows, err := s.db.Query(b.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.ActorWithFilms
	lastID := -1

	for rows.Next() {
		var actorID, filmRating int
		var filmTitle, filmDescription, actorName, actorSex string
		var filmDateCreation, actorDateBirth time.Time

		err := rows.Scan(&actorID, &actorName, &actorSex, &actorDateBirth, &filmTitle, &filmDescription, &filmDateCreation, &filmRating)
		if err != nil {
			return nil, err
		}

		// если актера меняется в списке добавляем новый элемент
		if actorID != lastID || lastID == -1 {
			result = append(result, &model.ActorWithFilms{
				Name:      actorName,
				Sex:       actorSex,
				DateBirth: actorDateBirth,
				Films:     make([]*model.FilmMinInfo, 0),
			})
		}
		// а фильм всегда вставляем в последний
		result[len(result)-1].Films = append(
			result[len(result)-1].Films,
			&model.FilmMinInfo{
				Title:        filmTitle,
				Description:  filmDescription,
				DateCreation: filmDateCreation,
				Rating:       filmRating,
			},
		)

		lastID = actorID
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}
