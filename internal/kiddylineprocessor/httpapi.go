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
		kp.loger.Panic("httpAPIServer : ListenAndServe : ", err)
	}
}

func (kp *Kiddylineprocessor) httpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kp.loger.Trace("Kiddylineprocessor : httpHandler ready")
		if err := kp.store.PindDB(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if kp.errorBaseball != "" || kp.errorFootball != "" || kp.errorSoccer != "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
