package postgres

import (
	"database/sql"
	"vk_test/internal/model"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Create(actor *model.Actor) error {
	_, err := s.db.Exec(
		"INSERT INTO Actors(a_name, a_sex, a_birth_date) VALUES ($1, $2, $3)",
		actor.Name,
		actor.Sex,
		actor.DateBirth,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Delete(actorName string) error {
	_, err := s.db.Exec("DELETE FROM Actors WHERE a_name = $1)", actorName)
	if err != nil {
		return err
	}
	return nil
}
