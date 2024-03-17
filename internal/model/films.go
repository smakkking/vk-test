package model

import "time"

type Film struct {
	Title        string    `json:"title"`
	Description  string    `json:"desc"`
	DateCreation time.Time `json:"date_creation"`
	Rating       int       `json:"rating"`
	ActorIDList  []int     `json:"actors"`
}

type FilmMinInfo struct {
	Title        string    `json:"title"`
	Description  string    `json:"desc"`
	DateCreation time.Time `json:"date_creation"`
	Rating       int       `json:"rating"`
}

type FilmWithActors struct {
	Title        string    `json:"title"`
	Description  string    `json:"desc"`
	DateCreation time.Time `json:"date_creation"`
	Rating       int       `json:"rating"`
	ActorList    []*Actor  `json:"actors"`
}
