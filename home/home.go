package home

import (
	"log"
	"net/http"
	"time"
)

const message = "url-shortener"
const homeDir = "/"

type Handlers struct {
	logger *log.Logger
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {

	// h.logger.Println("Home request processed")

	// if 'mux.NewRouter' then 404 by default
	if r.URL.Path != homeDir {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	// instead of:
	// h.logger.Println("Home request processed")
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("Home request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}

func (h *Handlers) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Logger(h.Home))
}
