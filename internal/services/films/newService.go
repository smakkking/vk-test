package films

func NewService(storage Storage) *Service {
	return &Service{
		filmsStorage: storage,
	}
}
