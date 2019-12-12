package main

import (
	"log"
	"net/http"
	"os"

	"github.com/michalswi/url-shortener/home"
	"github.com/michalswi/url-shortener/server"
)

// SERVICE_ADDR=:8080 go run main.go

var (
	ServiceAddr = os.Getenv("SERVICE_ADDR")
)

func main() {
	logger := log.New(os.Stdout, "shortener ", log.LstdFlags|log.Lshortfile)

	h := home.NewHandlers(logger)

	mux := http.NewServeMux()
	// mux.HandleFunc("/", h.Home) >>OR>> mux.NewRouter() >>OR>>:
	h.Routes(mux)

	srv := server.NewServer(mux, ServiceAddr)

	logger.Println("Server starting")
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
