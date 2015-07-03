package main

import (
	"fmt"
	"github.com/yihungjen/registry-driver/driver/redis"
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

	domain := r.Form.Get("domain")
	tier := r.Form.Get("tier")
	app := r.Form.Get("app")

	key := fmt.Sprintf("%s:%s:%s", tier, domain, app)
	if _, ok := Expect[key]; !ok {
		Expect[key] = make(chan []byte, 1)
	}
	defer func() { delete(Expect, key) }()

	exp, _ := time.ParseDuration(r.Form.Get("expire"))
	if err := client.SetNX(key, "", exp).Err(); err != nil {
		panic(err)
	}

	if response, ok := <-Expect[key]; !ok {
		panic(fmt.Errorf("Very bad thing happend for registry key: %s", key))
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
