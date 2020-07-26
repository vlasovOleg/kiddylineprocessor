package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq" // ...
	"github.com/vlasovoleg/kiddyLineProcessor/internal/store"
)

// Store for working with postgresql
type Store struct {
	db                 *sql.DB
	baseballRepository *baseballRepository
	soccerRepository   *soccerRepository
	footballRepository *footballRepository
}

// New sore
func New(db *sql.DB) *Store {
	s := &Store{}
	s.db = db
	s.baseballRepository = &baseballRepository{store: s}
	s.soccerRepository = &soccerRepository{store: s}
	s.footballRepository = &footballRepository{store: s}
	return s
}

// PindDB - check db connection
func (s *Store) PindDB() error {
	return s.db.Ping()
}

// BaseballRepository return interface for baseball data repository
func (s *Store) BaseballRepository() store.BaseballRepository {
	return s.baseballRepository
}

// SoccerRepository return interface for soccer data repository
func (s *Store) SoccerRepository() store.SoccerRepository {
	return s.soccerRepository
}

// FootballRepository return interface for football data repository
func (s *Store) FootballRepository() store.FootballRepository {
	return s.footballRepository
}
