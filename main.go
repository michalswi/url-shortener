package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"

	"github.com/michalswi/url-shortener/home"
	"github.com/michalswi/url-shortener/links"
	"github.com/michalswi/url-shortener/server"
)

var (
	ServiceAddr = os.Getenv("SERVICE_ADDR")
)

func main() {
	logger := log.New(os.Stdout, "shortener ", log.LstdFlags|log.Lshortfile)

	h := home.NewHandlers(logger)
	l := links.NewHandlers(logger, ServiceAddr)

	r := mux.NewRouter()
	h.Routes(r)
	l.Routes(r)
	srv := server.NewServer(r, ServiceAddr)

	logger.Printf("Server starting\n")
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
	logger.Printf("Server stopped\n")
}
