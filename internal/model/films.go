package model

import "time"

type Film struct {
	Title        string    `json:"title"`
	Description  string    `json:"desc"`
	DateCreation time.Time `json:"date_creation"`
	Rating       int       `json:"rating"`
	ActorIDList  []int     `json:"actors"`
}

type FilmPartialUpdate struct {
	Title     string `json:"title,omitempty"`
	TitleBool bool   `json:"-"`

	Description     string `json:"desc,omitempty"`
	DescriptionBool bool   `json:"-"`

	DateCreation     time.Time `json:"date_creation,omitempty"`
	DateCreationBool bool      `json:"-"`

	Rating     int  `json:"rating,omitempty"`
	RatingBool bool `json:"-"`

	ActorIDList     []int `json:"actors,omitempty"`
	ActorIDListBool bool  `json:"-"`
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
