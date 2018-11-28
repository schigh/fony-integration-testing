package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()
	router.Get(`/{words:.*}`, func(rw http.ResponseWriter, r *http.Request) {
		sum := md5.Sum([]byte(chi.URLParam(r, "words")))
		bytesOut := make([]byte, md5.Size*2)
		hex.Encode(bytesOut, sum[:])

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = rw.Write([]byte(fmt.Sprintf(`{"value":"%s"}`, string(bytesOut))))
	})
	_ = http.ListenAndServe("0.0.0.0:8883", router)
}
