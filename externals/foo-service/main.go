package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-chi/chi"
)

var re *regexp.Regexp

func init() {
	re = regexp.MustCompile(`(?i)foo`)
}

func main() {
	router := chi.NewRouter()
	router.Get(`/{words:.*}`, func(rw http.ResponseWriter, r *http.Request) {
		words := chi.URLParam(r, "words")
		unfoo := re.ReplaceAllString(words, "bar")
		rw.Header().Set("Content-Type", "application/json; charset=utf8")
		_, _ = rw.Write([]byte(fmt.Sprintf(`{"value":"%s"}`, unfoo)))
	})
	_ = http.ListenAndServe("0.0.0.0:8882", router)
}
