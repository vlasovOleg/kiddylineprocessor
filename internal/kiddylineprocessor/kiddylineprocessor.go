package kiddylineprocessor

import (
	"database/sql"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/vlasovoleg/kiddyLineProcessor/internal/store"
	"github.com/vlasovoleg/kiddyLineProcessor/internal/store/sqlstore"
)

// Kiddylineprocessor ...
type Kiddylineprocessor struct {
	config *Config
	store  store.Store
	loger  *logrus.Logger
}

// New kiddylineprocessor
func New(config *Config) *Kiddylineprocessor {
	kp := &Kiddylineprocessor{}
	kp.config = config

	loger := logrus.New()
	loger.SetFormatter(
		&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "02 15:04:05",
		},
	)
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		log.Panic(err)
	}
	loger.SetLevel(level)
	kp.loger = loger

	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		loger.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		loger.Panic(err)
	}
	kp.store = sqlstore.New(db)

	return kp
}
