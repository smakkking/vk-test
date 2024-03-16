package actors

import (
	"vk_test/internal/model"

	"github.com/google/uuid"
)

type Storage interface {
	Create(*model.Actor) error
	Update(uuid.UUID) error
	Delete(uuid.UUID) error
}
