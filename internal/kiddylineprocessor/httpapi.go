package kiddylineprocessor

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func (kp *Kiddylineprocessor) httpAPIServer(ctx context.Context, wg *sync.WaitGroup) {
	router := mux.NewRouter()
	router.HandleFunc("/ready", kp.httpHandler()).Methods("GET")
	server := http.Server{
		Addr:         kp.config.HTTPAPI.Address,
		Handler:      router,
		ReadTimeout:  kp.config.HTTPAPI.ReadTimeout,
		WriteTimeout: kp.config.HTTPAPI.WriteTimeout,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			kp.loger.Info("httpAPIServer : ListenAndServe : ", err)
		}
	}()

	<-ctx.Done()
	server.Shutdown(context.TODO())
	wg.Done()
}

func (kp *Kiddylineprocessor) httpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kp.loger.Trace("Kiddylineprocessor : httpHandler ready")
		if err := kp.store.PindDB(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(kp.errors.data) != 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
