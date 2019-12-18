package links

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
)

const postDir = "/links"

var (
	shortUrlID  string
	redirectURL string
)

type Handlers struct {
	logger        *log.Logger
	serverAddress string
	storeAddress  string
	dnsName       string
}

type getUrl struct {
	LongUrl string `json:"longUrl"`
}

type dataUrl struct {
	Id        string `json:"id"`
	LongUrl   string `json:"longUrl"`
	ShortUrl  string `json:"shortUrl"`
	CreatedAt string `json:"createdAt"`
}

type healthState struct {
	State              string   `json:"state"`
	UrlErrorMessages   []string `json:"urlerrormessages"`
	Redis              string   `json:"redis"`
	RedisErrorMessages []string `json:"rediserrormessages"`
}

// keep Id + longUrl
var keepData = make(map[string]string)

func (l *Handlers) Links(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var geturl getUrl

	err := json.NewDecoder(r.Body).Decode(&geturl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	current := time.Now()
	urlID := genID(4)
	redirectURL := geturl.LongUrl
	genShortURL := fmt.Sprintf("http://%s:%s/%s", l.dnsName, l.serverAddress, urlID)

	// cache
	keepData[urlID] = redirectURL

	dataurl := &dataUrl{
		Id:        urlID,
		LongUrl:   redirectURL,
		ShortUrl:  genShortURL,
		CreatedAt: current.Format("2006-01-02 15:04:0512"),
	}

	b, err := json.Marshal(dataurl)
	if err != nil {
		l.logger.Printf("Marshaling failed: %v\n", err)
	}

	w.Write([]byte(b))
}

func (l *Handlers) Proxy(w http.ResponseWriter, r *http.Request) {
	l.logger.Println("Proxy request processed")
	for k, v := range keepData {
		if r.URL.Path == "/"+k {
			http.Redirect(w, r, v, http.StatusMovedPermanently)
		}
	}
}

func (l *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer l.logger.Printf("Links request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

func NewHandlers(logger *log.Logger, serverAddress string, storeAddress string, dnsName string) *Handlers {
	return &Handlers{
		logger:        logger,
		serverAddress: serverAddress,
		storeAddress:  storeAddress,
		dnsName:       dnsName,
	}
}

func (l *Handlers) Routes(mux *mux.Router) {
	mux.HandleFunc("/links", l.Logger(l.Links)).Methods("POST")
	mux.HandleFunc("/health", l.Logger(l.healthcheck)).Methods("GET")
	mux.HandleFunc("/{shortUrlID}", l.Logger(l.Proxy)).Methods("GET")
}

func (l *Handlers) healthcheck(w http.ResponseWriter, r *http.Request) {
	l.logger.Println("Healthcheck request processed")

	hs := &healthState{}

	// urlshortener
	resWeb, err := http.Get(fmt.Sprintf("http://localhost:%s", l.serverAddress))
	if err != nil {
		l.logger.Printf("Check failed: %v\n", err)
	} else {
		defer resWeb.Body.Close()

		// if err != nil || resWeb.StatusCode != 200 {
		if resWeb.StatusCode != 200 {
			l.logger.Printf("Status code error: %d, Status error: %s\n", resWeb.StatusCode, resWeb.Status)
			hs.UrlErrorMessages = append(hs.UrlErrorMessages, fmt.Sprintf("HealthError: %s", resWeb.Status))
		}

		if len(hs.UrlErrorMessages) > 0 {
			hs.State = "NOK"
		} else {
			hs.State = "OK"
		}
	}

	// redis
	resRedis, err := http.Get(fmt.Sprintf("http://localhost:%s", l.storeAddress))
	if err != nil {
		l.logger.Printf("Check failed: %v", err)
		hs.Redis = "NOK"
		hs.RedisErrorMessages = append(hs.RedisErrorMessages, fmt.Sprintf("HealthError: %v", err))
	} else {
		defer resRedis.Body.Close()

		if resRedis.StatusCode != 200 {
			l.logger.Printf("Status code error: %d, Status error: %s\n", resRedis.StatusCode, resRedis.Status)
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
		log.Fatalf("Marshaling failed: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(b))
}

// generate random ID
func genID(length int) (randomID string) {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
