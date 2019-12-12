package proxy

import (
	"fmt"
	"github.com/michalswi/url-shortener/links"
	"net/http"
)

func Proxy(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(links.ShortUrlID)
	fmt.Printf(links.RedirectURL)
	http.Redirect(w, r, "http://www.google.com", 301)
}
