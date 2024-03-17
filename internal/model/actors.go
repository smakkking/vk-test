package model

import "time"

type Actor struct {
	Name      string    `json:"name"`
	Sex       string    `json:"sex"`
	DateBirth time.Time `json:"date_birth"`
}

type ActorWithFilms struct {
	Name      string         `json:"name"`
	Sex       string         `json:"sex"`
	DateBirth time.Time      `json:"date_birth"`
	Films     []*FilmMinInfo `json:"films"`
}
