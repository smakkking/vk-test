package model

import "time"

type Actor struct {
	Name      string // id
	Sex       string
	DateBirth time.Time
}

type ActorWithFilms struct {
	Name      string
	Sex       string
	DateBirth time.Time
	Films     []*FilmMinInfo
}
