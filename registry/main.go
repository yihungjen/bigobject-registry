package main

import (
	"github.com/yihungjen/registry-driver/driver/redis"
	"log"
	"net/http"
)

var (
	// map completion mapper for update key
	Expect map[string]chan []byte

	// Request endpoint multiplexer at PORT 9090
	Server = http.NewServeMux()
)

func main() {
	log.Println("begin serving Registry...")
	s := &http.Server{Addr: ":80", Handler: Server}
	log.Fatalln(s.ListenAndServe())
}

func init() {
	Expect = make(map[string]chan []byte)

	// spawn worker to monitor key space change; KEY UPDATE
	go func() {
		regevent := redis.KeySpaceEventLoop()

		client := redis.NewClient()
		defer client.Close()

		for event := range regevent {
			target, ok := Expect[event.Key]
			if !ok {
				continue
			}
			switch event.Action {
			case "expired":
				close(target)
				break
			case "append":
				val, err := client.Get(event.Key).Result()
				if err != nil {
					panic(err)
				}
				target <- []byte(val)
				break
			}
		}
	}()

	Server.HandleFunc("/provision", Provision)
	Server.HandleFunc("/deprovision", Deprovision)
}
