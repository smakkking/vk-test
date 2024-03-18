package model

type Actor struct {
	Name      string    `json:"name"`
	Sex       string    `json:"sex"`
	DateBirth CivilTime `json:"date_birth"`
}

type ActorPartialUpdate struct {
	Name     string `json:"name,omitempty"`
	NameBool bool   `json:"-"`

	Sex     string `json:"sex,omitempty"`
	SexBool bool   `json:"-"`

	DateBirth     CivilTime `json:"date_birth,omitempty"`
	DateBirthBool bool      `json:"-"`
}

type ActorWithFilms struct {
	Name      string         `json:"name"`
	Sex       string         `json:"sex"`
	DateBirth CivilTime      `json:"date_birth"`
	Films     []*FilmMinInfo `json:"films"`
}
