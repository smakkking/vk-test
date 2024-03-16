package actors

func NewService(storage Storage) *Service {
	return &Service{
		actorStorage: storage,
	}
}
