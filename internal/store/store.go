package store

// Store - interface for working with database
type Store interface {
	PindDB() error
	BaseballRepository() BaseballRepository
	SoccerRepository() SoccerRepository
	FootballRepository() FootballRepository
}
