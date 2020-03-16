package health

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Handlers struct {
	logger        *log.Logger
	serverAddress string
	apiprefix     string
	storeAddress  string
}

type healthState struct {
	State              string   `json:"state"`
	UrlErrorMessages   []string `json:"urlerrormessages"`
	Redis              string   `json:"redis"`
	RedisErrorMessages []string `json:"rediserrormessages"`
}

func NewHandlers(logger *log.Logger, serverAddress string, apiprefix string, storeAddress string) *Handlers {
	return &Handlers{
		logger:        logger,
		serverAddress: serverAddress,
		apiprefix:     apiprefix,
		storeAddress:  storeAddress,
	}
}

func (h *Handlers) LinkRoutes(mux *mux.Router) {
	mux.HandleFunc("/ok", h.Logger(h.statusOK)).Methods("GET")
	mux.HandleFunc("/health", h.Logger(h.healthcheck)).Methods("GET")
	mux.HandleFunc("/healthz", h.Logger(h.healthz)).Methods("GET")
}

func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	// instead of:
	// h.logger.Println("Home request processed")
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("Health request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

// Manage 'healthcheck' endpoint, return json
func (h *Handlers) healthcheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("healthcheck request processed")

	hs := &healthState{}

	// urlshortener
	resWeb, err := http.Get(fmt.Sprintf("http://localhost:%s%s/ok", h.serverAddress, h.apiprefix))
	if err != nil {
		h.logger.Printf("Check failed: %v\n", err)
	} else {
		defer resWeb.Body.Close()
		// if err != nil || resWeb.StatusCode != 200 {
		if resWeb.StatusCode != 200 {
			h.logger.Printf("Status code error: %d, Status error: %s\n", resWeb.StatusCode, resWeb.Status)
			hs.UrlErrorMessages = append(hs.UrlErrorMessages, fmt.Sprintf("HealthError: %s", resWeb.Status))
		}
		if len(hs.UrlErrorMessages) > 0 {
			hs.State = "NOK"
		} else {
			hs.State = "OK"
		}
	}

	// redis
	resRedis, err := http.Get(fmt.Sprintf("http://localhost:%s", h.storeAddress))
	if err != nil {
		h.logger.Printf("Check failed: %v", err)
		hs.Redis = "NOK"
		hs.RedisErrorMessages = append(hs.RedisErrorMessages, fmt.Sprintf("HealthError: %v", err))
	} else {
		defer resRedis.Body.Close()
		if resRedis.StatusCode != 200 {
			h.logger.Printf("Status code error: %d, Status error: %s\n", resRedis.StatusCode, resRedis.Status)
			hs.RedisErrorMessages = append(hs.RedisErrorMessages, fmt.Sprintf("HealthError: %s", resRedis.Status))
		}
		if len(hs.RedisErrorMessages) > 0 {
			hs.Redis = "NOK"
		} else {
			hs.Redis = "OK"
		}
	}
	// both
	b, err := json.Marshal(hs)
	if err != nil {
		h.logger.Printf("Marshaling failed: %v\n", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(b))
}

// Manage 'healthz' endpoint, return byte
func (h *Handlers) healthz(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("healthz request processed")
	resWeb, err := http.Get(fmt.Sprintf("http://localhost:%s%s/ok", h.serverAddress, h.apiprefix))
	if err != nil {
		h.logger.Printf("Check failed: %v\n", err)
	} else {
		defer resWeb.Body.Close()
		// if err != nil || resWeb.StatusCode != 200 {
		if resWeb.StatusCode != 200 {
			h.logger.Printf("Status code error: %d, Status error: %s\n", resWeb.StatusCode, resWeb.Status)
			w.Write([]byte("NOK"))
		} else {
			w.Write([]byte("OK"))
		}
	}
}

// Endpoint for health checks
func (h *Handlers) statusOK(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("statusOK processed\n")
	w.WriteHeader(http.StatusOK)
}
