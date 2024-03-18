package postgres

import (
	"database/sql"
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
		"INSERT INTO Actors(a_name, a_sex, a_birth_date) VALUES ($1, $2, $3) RETURNING a_id",
		actor.Name,
		actor.Sex,
		time.Time(actor.DateBirth),
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

func (s *Storage) UpdateActor(actorID int, actor *model.ActorPartialUpdate) error {
	b := strings.Builder{}
	// костыль
	if actor.Sex == "" {
		actor.Sex = "мужчина"
	}

	b.WriteString(`
		UPDATE Actors 
		SET 
			a_name = CASE WHEN $1::boolean THEN $2::TEXT ELSE a_name END,
			a_sex = CASE WHEN $3::boolean THEN $4::SEX ELSE a_sex END,
			a_birth_date = CASE WHEN $5::boolean THEN $6::DATE ELSE a_birth_date END
		WHERE a_id = $7;
	`)

	_, err := s.db.Exec(
		b.String(),
		actor.NameBool, actor.Name,
		actor.SexBool, actor.Sex,
		actor.DateBirthBool, time.Time(actor.DateBirth),
		actorID,
	)
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
		"INSERT INTO Films(f_title, f_desc, f_date_creation, f_rating) VALUES ($1, $2, $3, $4) RETURNING f_id",
		film.Title,
		film.Description,
		time.Time(film.DateCreation),
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

func (s *Storage) UpdateFilm(filmID int, film *model.FilmPartialUpdate) error {
	b := strings.Builder{}
	b.WriteString(`
		UPDATE Films 
		SET 
		f_title = CASE WHEN $1::boolean THEN $2::VARCHAR(150) ELSE f_title END,
		f_desc = CASE WHEN $3::boolean THEN $4::VARCHAR(1000) ELSE f_desc END,
		f_date_creation = CASE WHEN $5::boolean THEN $6::DATE ELSE f_date_creation END,
		f_rating = CASE WHEN $7::boolean THEN $8::INT ELSE f_rating END
		
		WHERE a_id = $9;
	`)

	// делаем в транзакции
	tX, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tX.Exec(
		b.String(),
		film.TitleBool, film.Title,
		film.DescriptionBool, film.Description,
		film.DateCreationBool, time.Time(film.DateCreation),
		film.RatingBool, film.Rating,
		filmID,
	)
	if err != nil {
		_ = tX.Rollback()
		return err
	}

	if film.ActorIDListBool {
		_, err = tX.Exec("DELETE FROM ActorToFilm WHERE film_id = $1", filmID)
		if err != nil {
			_ = tX.Rollback()
			return err
		}

		for _, actorID := range film.ActorIDList {
			_, err := tX.Exec(
				"INSERT INTO ActorToFilm(actor_id, film_id) VALUES($1, $2)",
				actorID,
				filmID,
			)
			if err != nil {
				_ = tX.Rollback()
				return err
			}
		}
	}
	_ = tX.Commit()

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
				DateCreation: model.CivilTime(filmDateCreation),
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
				DateBirth: model.CivilTime(actorDateBirth),
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
	b.WriteString(`
		SELECT f.f_id, f.f_title, f.f_desc, f.f_date_creation, f.f_rating, a.a_name, a.a_sex, a.a_birth_date
		FROM
			Films AS f
				JOIN ActorToFilm AS atf ON f.f_id = atf.film_id
				JOIN Actors AS a ON a.a_id = atf.actor_id
		WHERE a.a_name LIKE '%$1%'
		ORDER BY f.f_id;
	`)

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
				DateCreation: model.CivilTime(filmDateCreation),
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
				DateBirth: model.CivilTime(actorDateBirth),
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
	b.WriteString(`
		SELECT f.f_id, f.f_title, f.f_desc, f.f_date_creation, f.f_rating, a.a_name, a.a_sex, a.a_birth_date
		FROM
			Films AS f
				JOIN ActorToFilm AS atf ON f.f_id = atf.film_id
				JOIN Actors AS a ON a.a_id = atf.actor_id
		WHERE f.f_title LIKE '%$1%'
		ORDER BY f.f_id;
	`)

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
				DateCreation: model.CivilTime(filmDateCreation),
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
				DateBirth: model.CivilTime(actorDateBirth),
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
	b.WriteString(`
		SELECT a.a_id, a.a_name, a.a_sex, a.a_birth_date, f.f_title, f.f_desc, f.f_date_creation, f.f_rating
		FROM
			Actors as a
				JOIN ActorToFilm AS atf ON a.a_id = atf.actor_id
				JOIN Films as f ON f.f_id = atf.film_id
		ORDER BY a.a_id;
	`)

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
				DateBirth: model.CivilTime(actorDateBirth),
				Films:     make([]*model.FilmMinInfo, 0),
			})
		}
		// а фильм всегда вставляем в последний
		result[len(result)-1].Films = append(
			result[len(result)-1].Films,
			&model.FilmMinInfo{
				Title:        filmTitle,
				Description:  filmDescription,
				DateCreation: model.CivilTime(filmDateCreation),
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
