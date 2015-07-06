package main

import (
	"fmt"
	"github.com/yihungjen/registry-driver/driver/redis"
	"log"
	"net/http"
	"time"
)

func Provision(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowd", 405)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	client := redis.NewClient()
	defer client.Close()

	domain, tier, app := r.Form.Get("domain"), r.Form.Get("tier"), r.Form.Get("app")
	if domain == "" || tier == "" || app == "" {
		http.Error(w, "Bad request", 400)
		return
	}

	key := fmt.Sprintf("%s:%s:%s", tier, domain, app)
	if _, ok := Expect[key]; !ok {
		Expect[key] = make(chan []byte)
	}
	defer func() {
		delete(Expect, key)
		if err := recover(); err != nil {
			log.Println(err)
			switch t := err.(type) {
			default:
				http.Error(w, "Unexpected Error", 503)
				break
			case error:
				http.Error(w, t.Error(), 503)
				break
			case string:
				http.Error(w, t, 503)
				break
			}
		}
	}()

	exp, _ := time.ParseDuration(r.Form.Get("expire"))
	if state, err := client.SetNX(key, "", exp).Result(); err != nil {
		panic(err)
	} else if !state {
		panic(fmt.Errorf("KEY exist: %s", key))
	}

	if response, ok := <-Expect[key]; !ok {
		panic(fmt.Errorf("Key expired before response %s", key))
	} else {
		headers := w.Header()
		headers.Add("Content-Type", "application/json")
		w.Write(response)
	}
}

func Deprovision(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unable to delete resource", 503)
	return
}
