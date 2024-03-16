package model

import "time"

type Actor struct {
	Name      string
	Sex       string
	DateBirth time.Time
	Films     []*Film
}
