package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()
	router.Get(`/{idx}`, func(rw http.ResponseWriter, r *http.Request) {
		_idx := chi.URLParam(r, "idx")
		idx, _ := strconv.Atoi(_idx)

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = rw.Write([]byte(fmt.Sprintf(`{"value":"your index is %d"}`, idx)))
	})
	_ = http.ListenAndServe("0.0.0.0:8885", router)
}
