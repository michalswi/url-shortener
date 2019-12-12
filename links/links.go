package links

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/golang/gddo/httputil/header"
)

const postDir = "/links"

var (
	ShortUrlID string
	RedirectURL string
)

type Handlers struct {
	logger        *log.Logger
	serverAddress string
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

func (l *Handlers) Links(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != postDir {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "403 Forbidden", http.StatusNotFound)
		return
	}

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
	ShortUrlID = genID()
	RedirectURL = geturl.LongUrl
	genShortURL := fmt.Sprintf("http://localhost%s/%s", l.serverAddress, ShortUrlID)

	dataurl := &dataUrl{
		Id:        ShortUrlID,
		LongUrl:   RedirectURL,
		ShortUrl:  genShortURL,
		CreatedAt: current.Format("2006-01-02 15:04:0512"),
	}

	b, err := json.Marshal(dataurl)
	if err != nil {
		l.logger.Printf("Marshaling failed: %v\n", err)
	}

	w.Write([]byte(b))
}

func (l *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer l.logger.Printf("Links request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

func NewHandlers(logger *log.Logger, serverAddress string) *Handlers {
	return &Handlers{
		logger:        logger,
		serverAddress: serverAddress,
	}
}

func (l *Handlers) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/links", l.Logger(l.Links))
}

func genID() (randomID string) {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	randomID = b.String()
	return randomID
}
