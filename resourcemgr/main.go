package main

import (
	"github.com/yihungjen/bigobject-registry/resourcemgr/service"
	"log"
	"net/http"
)

var (
	Api     = make(chan bool, 1)
	Qworker = make(chan bool, 1)
	Status  = make(chan bool, 1)
)

func main() {
	go func() {
		s := &http.Server{Addr: ":80", Handler: service.Server}
		log.Println(s.ListenAndServe())
		Api <- true
	}()

	log.Println("begin serving Resource in Registry")
	select {
	case <-Api:
		log.Fatalln("Service API worker halt")
	case <-Qworker:
		log.Fatalln("Service queue worker halt")
	case <-Status:
		log.Fatalln("Service upkeep status worker halt")
	}
}
