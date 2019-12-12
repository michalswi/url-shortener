package main

import (
	"log"
	"net/http"
	"os"

	"github.com/michalswi/url-shortener/home"
	"github.com/michalswi/url-shortener/links"
	"github.com/michalswi/url-shortener/proxy"
	"github.com/michalswi/url-shortener/server"
)

// SERVICE_ADDR=:8080 go run main.go
// curl -H "Content-Type: application/json" -X POST -d '{"longUrl":"https://golang.org/doc/effective_go.html"}' localhost:8080/links

var (
	ServiceAddr = os.Getenv("SERVICE_ADDR")
)

func main() {
	logger := log.New(os.Stdout, "shortener ", log.LstdFlags|log.Lshortfile)

	h := home.NewHandlers(logger)
	l := links.NewHandlers(logger, ServiceAddr)

	mux := http.NewServeMux()
	// mux.HandleFunc("/", h.Home) >>OR>> mux.NewRouter() >>OR>>:
	h.Routes(mux)
	l.Routes(mux)

	mux.HandleFunc("/test", proxy.Proxy)

	srv := server.NewServer(mux, ServiceAddr)

	logger.Printf("Server starting\n")
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
	logger.Printf("Server stopped\n")
}
