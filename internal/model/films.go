package model

import "time"

type Film struct {
	Title        string
	Description  string
	DateCreation time.Time
	Rating       int
	ActorList    []*Actor
}
