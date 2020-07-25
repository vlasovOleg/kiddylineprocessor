package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq" // ...
	"github.com/vlasovoleg/kiddyLineProcessor/internal/store"
)

// Store ...
type Store struct {
	db                 *sql.DB
	baseballRepository *BaseballRepository
	soccerRepository   *SoccerRepository
	footballRepository *FootballRepository
}

// New sore
func New(db *sql.DB) *Store {
	s := &Store{}
	s.db = db
	s.baseballRepository = &BaseballRepository{store: s}
	s.soccerRepository = &SoccerRepository{store: s}
	s.footballRepository = &FootballRepository{store: s}
	return s
}

// PindDB - check db connection
func (s *Store) PindDB() error {
	return s.db.Ping()
}

// BaseballRepository ...
func (s *Store) BaseballRepository() store.BaseballRepository {
	return s.baseballRepository
}

// SoccerRepository ...
func (s *Store) SoccerRepository() store.SoccerRepository {
	return s.soccerRepository
}

// FootballRepository ...
func (s *Store) FootballRepository() store.FootballRepository {
	return s.footballRepository
}
