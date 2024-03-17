package inmemory

import "vk_test/internal/model"

func NewStorage() (*Storage, error) {
	return &Storage{
		Actors: make(map[string]*model.Actor),
		Films:  make(map[string]*model.Film),
	}, nil
}
