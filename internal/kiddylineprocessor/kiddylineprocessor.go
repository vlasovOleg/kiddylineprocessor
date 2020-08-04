package kiddylineprocessor

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/vlasovoleg/kiddyLineProcessor/internal/store"
	"github.com/vlasovoleg/kiddyLineProcessor/internal/store/sqlstore"
)

// Kiddylineprocessor main struct
type Kiddylineprocessor struct {
	config        *Config
	store         store.Store
	httpClient    *http.Client
	loger         *logrus.Logger
	errorBaseball string
	errorFootball string
	errorSoccer   string
}

// New kiddylineprocessor
func New(config *Config) *Kiddylineprocessor {
	fmt.Printf("%+v\n", config)
	kp := &Kiddylineprocessor{}
	kp.config = config
	kp.config.LinesProvider.Address += "/api/v1/lines/"

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
		loger.Panic("kiddylineprocessor : New : sql.Open : ", err)
	}
	kp.store = sqlstore.New(db)

	kp.httpClient = &http.Client{
		Timeout: config.LinesProvider.RequestsTimeout,
	}

	kp.errorBaseball = "waiting sync"
	kp.errorFootball = "waiting sync"
	kp.errorSoccer = "waiting sync"
	return kp
}

// Start updater by line provider for baseball, football, soccer.
// Start http api and grpc
func (kp *Kiddylineprocessor) Start() {
	kp.loger.Debug("Kiddylineprocessor : Start")

	_, stop := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	go kp.updaterByLineProviderBaseball()
	go kp.updaterByLineProviderFootball()
	go kp.updaterByLineProviderSoccer()

	go kp.httpAPIServer()
	go kp.NewGRPC(&kp.store, kp.loger)

	sigCh := make(chan os.Signal, 10)
	signal.Notify(sigCh, os.Interrupt)
	for {
		if <-sigCh == os.Interrupt {
			stop()
			wg.Wait()
			return
		}
		continue
	}
}
