package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/michalswi/url-shortener/home"
	"github.com/michalswi/url-shortener/links"
	"github.com/michalswi/url-shortener/server"
)

var version = "0.0.1"

func main() {
	logger := log.New(os.Stdout, "shortener ", log.LstdFlags|log.Lshortfile)

	ServiceAddr := os.Getenv("SERVICE_ADDR")
	PprofAddr := os.Getenv("PPROF_ADDR")
	StoreAddr := os.Getenv("STORE_ADDR")
	DnsName := os.Getenv("DNS_NAME")

	h := home.NewHandlers(logger, version)
	l := links.NewHandlers(logger, ServiceAddr, StoreAddr, DnsName)

	r := mux.NewRouter()
	h.LinkRoutes(r)
	l.LinkRoutes(r)
	srv := server.NewServer(r, ServiceAddr)

	// start server
	go func() {
		logger.Printf("Starting server on port %s \n", ServiceAddr)
		err := srv.ListenAndServe()
		if err != nil {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// start pprof
	if PprofAddr != "" {
		go func() {
			logger.Printf("Starting pprof server on port %s \n", PprofAddr)
			if err := http.ListenAndServe(fmt.Sprintf(":%v", PprofAddr), pprofRouter()); err != nil {
				logger.Fatalf("Pprof server failed to start on %s: %v\n", PprofAddr, err)
			}
		}()
	}

	// shutdown server
	gracefulShutdown(srv, logger)
}

// graceful shutdown
func gracefulShutdown(srv *http.Server, logger *log.Logger) {

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}

	logger.Printf("Shutting down the server...\n")
	os.Exit(0)
}
