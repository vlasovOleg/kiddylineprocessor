package kiddylineprocessor

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (kp *Kiddylineprocessor) httpAPIServer() {
	router := mux.NewRouter()
	router.HandleFunc("/ready", kp.httpHandler()).Methods("GET")
	server := http.Server{
		Addr:         kp.config.HTTPserverAddress,
		Handler:      router,
		ReadTimeout:  kp.config.HTTPserverReadTimeout * time.Second,
		WriteTimeout: kp.config.HTTPserverWriteTimeout * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		kp.loger.Fatal("httpAPIServer : ListenAndServe : ", err)
	}
}

func (kp *Kiddylineprocessor) httpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := kp.store.PindDB(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(kp.errorBaseball) != 0 || len(kp.errorFootball) != 0 || len(kp.errorSoccer) != 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
