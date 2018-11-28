package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ddliu/go-httpclient"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/schigh/taskchain"
)

// Thingy is a thing-y
type Thingy struct {
	Value string `json:"value"`
}

// AllThings is all the thingies
type AllThings struct {
	Caps    *Thingy `json:"caps"`
	Foo     *Thingy `json:"foo"`
	Md5     *Thingy `json:"md5"`
	Reverse *Thingy `json:"reverse"`
}

func addTask(tg *taskchain.TaskGroup, words string, env string, key string) error {
	endpoint := fmt.Sprintf("%s/%s", os.Getenv(env), words)
	resp, err := httpclient.Get(endpoint)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("%s-get", key))
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("%s-readall", key))
	}
	t := &Thingy{}
	if err := json.Unmarshal(b, t); err != nil {
		return errors.Wrap(err, fmt.Sprintf("%s-unmarshal", key))
	}

	tg.Set(key, t)
	return nil
}

func main() {
	// this is a no-op in a container environment
	_ = godotenv.Load()

	router := chi.NewRouter()
	router.Get(`/{words:.*}`, func(rw http.ResponseWriter, r *http.Request) {
		words := chi.URLParam(r, "words")

		group := taskchain.TaskGroup{}
		group.Add(func(tg *taskchain.TaskGroup) error {
			return addTask(tg, words, "CAPS_SERVICE_URL", "CAPS")
		})
		group.Add(func(tg *taskchain.TaskGroup) error {
			return addTask(tg, words, "FOO_SERVICE_URL", "FOO")
		})
		group.Add(func(tg *taskchain.TaskGroup) error {
			return addTask(tg, words, "MD5_SERVICE_URL", "MD5")
		})
		group.Add(func(tg *taskchain.TaskGroup) error {
			return addTask(tg, words, "REVERSE_SERVICE_URL", "REVERSE")
		})

		if err := group.Exec(); err != nil {
			fmt.Printf("%v\n", err)
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		caps := group.Get("CAPS", &Thingy{}).(*Thingy)
		foo := group.Get("FOO", &Thingy{}).(*Thingy)
		md5 := group.Get("MD5", &Thingy{}).(*Thingy)
		reverse := group.Get("REVERSE", &Thingy{}).(*Thingy)

		allThings := &AllThings{
			Caps:caps,
			Foo:foo,
			Md5:md5,
			Reverse:reverse,
		}

		b, err := json.Marshal(allThings)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json; charset=utf8")
		_, _ = rw.Write(b)
	})
	_ = http.ListenAndServe("0.0.0.0:80", router)
}
