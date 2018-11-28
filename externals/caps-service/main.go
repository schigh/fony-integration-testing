package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()
	router.Get(`/{words:.*}`, func(rw http.ResponseWriter, r *http.Request) {
		words := chi.URLParam(r, "words")
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = rw.Write([]byte(fmt.Sprintf(`{"value":"%s"}`, strings.ToUpper(words))))
	})
	_ = http.ListenAndServe("0.0.0.0:8881", router)
}
