package store

// Store - interface for working with database
type Store interface {
	BaseballRepository() BaseballRepository
	SoccerRepository() SoccerRepository
	FootballRepository() FootballRepository
}
