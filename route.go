package main

import (
	"net/http/pprof"

	"github.com/gorilla/mux"
)

// pprof server
func pprofRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	return router
}
