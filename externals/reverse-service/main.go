package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()
	router.Get(`/{words:.*}`, func(rw http.ResponseWriter, r *http.Request) {
		runes := []rune(chi.URLParam(r, "words"))
		l := len(runes)
		for i, j := 0, l-1; i < l/2; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		rw.Header().Set("Content-Type", "application/json; charset=utf8")
		_, _ = rw.Write([]byte(fmt.Sprintf(`{"value":"%s"}`, string(runes))))
	})
	_ = http.ListenAndServe("0.0.0.0:8884", router)
}
