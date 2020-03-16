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

var (
	shortUrlID  string
	redirectURL string
)

type Handlers struct {
	logger        *log.Logger
	serverAddress string
	storeAddress  string
	dnsName       string
	apipath       string
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

	genShortURL := fmt.Sprintf("http://%s:%s%s/%s", l.dnsName, l.serverAddress, l.apipath, urlID)
	// var genShortURL string
	// if l.serverAddress == "80" {
	// 	genShortURL = fmt.Sprintf("http://%s%s/%s", l.dnsName, l.apipath, urlID)
	// } else {
	// 	genShortURL = fmt.Sprintf("http://%s:%s%s/%s", l.dnsName, l.serverAddress, l.apipath, urlID)
	// }

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

// Proxy traffic to url base on shortUrl
func (l *Handlers) Proxy(w http.ResponseWriter, r *http.Request) {
	l.logger.Println("Proxy request processed")
	for k, v := range keepData {
		if r.URL.Path == l.apipath+"/"+k {
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

func NewHandlers(logger *log.Logger, serverAddress string, storeAddress string, dnsName string, apipath string) *Handlers {
	return &Handlers{
		logger:        logger,
		serverAddress: serverAddress,
		storeAddress:  storeAddress,
		dnsName:       dnsName,
		apipath:       apipath,
	}
}

func (l *Handlers) LinkRoutes(mux *mux.Router) {
	mux.HandleFunc("/links", l.Logger(l.Links)).Methods("POST")
	mux.HandleFunc("/{shortUrlID}", l.Logger(l.Proxy)).Methods("GET")
}

// Generate random ID
func genID(length int) (randomID string) {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
