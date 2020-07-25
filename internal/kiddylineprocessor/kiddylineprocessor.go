package kiddylineprocessor

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vlasovoleg/kiddyLineProcessor/internal/store"
	"github.com/vlasovoleg/kiddyLineProcessor/internal/store/sqlstore"
)

// Kiddylineprocessor ...
type Kiddylineprocessor struct {
	config        *Config
	store         store.Store
	httpClient    *http.Client
	loger         *logrus.Logger
	ready         bool
	errorBaseball string
	errorFootball string
	errorSoccer   string
}

// New kiddylineprocessor
func New(config *Config) *Kiddylineprocessor {
	kp := &Kiddylineprocessor{}
	kp.config = config
	kp.config.LinesProviderAddress += "/api/v1/lines/"

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

	kp.httpClient = &http.Client{
		Timeout: config.LinesProviderRequestsTimeout * time.Second,
	}

	kp.errorBaseball = "waiting sync"
	kp.errorFootball = "waiting sync"
	kp.errorSoccer = "waiting sync"
	return kp
}

// Start ...
func (kp Kiddylineprocessor) Start() {
	kp.loger.Debug("Kiddylineprocessor : Start")

	go kp.updaterByLineProviderBaseball()
	go kp.updaterByLineProviderFootball()
	go kp.updaterByLineProviderSoccer()

	go kp.httpAPIServer()
	go kp.NewGRPS(&kp.store, kp.loger)
}
