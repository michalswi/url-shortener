package home

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Handlers struct {
	logger  *log.Logger
	version string
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	// h.logger.Println("Home request processed")

	message := "url-shortener"

	// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(message))

	hostname, err := os.Hostname()
	if err != nil {
		h.logger.Printf("Get 'host name' failed: %v\n", err)
	}

	version := h.version
	w.WriteHeader(http.StatusOK)

	var html = `
	<html>
	<h1>%s</h1>
	<p><b>Hostname</b>: %s; <b>Version</b>: %s</p>
	</html>
	`
	fmt.Fprintf(w, html, message, hostname, version)
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

func NewHandlers(logger *log.Logger, version string) *Handlers {
	return &Handlers{
		logger:  logger,
		version: version,
	}
}

func (h *Handlers) LinkRoutes(mux *mux.Router) {
	mux.HandleFunc("/home", h.Logger(h.Home))
}
