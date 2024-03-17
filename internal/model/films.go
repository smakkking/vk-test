package model

import "time"

type Film struct {
	Title        string
	Description  string
	DateCreation time.Time
	Rating       int
	ActorIDList  []int
}

type FilmMinInfo struct {
	Title        string
	Description  string
	DateCreation time.Time
	Rating       int
}

type FilmWithActors struct {
	Title        string // id
	Description  string
	DateCreation time.Time
	Rating       int
	ActorList    []*Actor
}
