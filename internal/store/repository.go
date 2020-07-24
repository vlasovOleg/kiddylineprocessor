package store

type _repository interface {
	GetCoefficient() (float32, error)
	UpdateCoefficient(coefficient float32) error
}

// BaseballRepository - interface for working with baseball data in store
type BaseballRepository interface {
	_repository
}

// FootballRepository - interface for working with football data in store
type FootballRepository interface {
	_repository
}

// SoccerRepository - interface for working with soccer data in store
type SoccerRepository interface {
	_repository
}
