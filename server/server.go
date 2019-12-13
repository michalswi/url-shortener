package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func NewServer(r *mux.Router, serverAddress string) *http.Server {

	srv := &http.Server{
		Addr:         serverAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      r,
	}
	return srv
}
