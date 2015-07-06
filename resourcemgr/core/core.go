package core

import (
	"fmt"
	"log"
	"runtime"
)

const (
	RESOURCE_REQ_BACKLOG = 10000
)

var (
	hub = make(chan *ResourceReq, RESOURCE_REQ_BACKLOG)
)

func create(name string) *ResourceInfo {
	return &ResourceInfo{fmt.Errorf("Not yet implemented")}
}

func remove(name string) error {
	return nil
}

func coreLoop() {
	log.Println("number of CPU:", runtime.NumCPU())
	log.Println("number of vCPU:", runtime.GOMAXPROCS(-1))
	go func() {
		for onereq := range hub {
			log.Println(onereq)
			switch onereq.Action {
			case "CREATE":
				onereq.Resp <- create(onereq.Name)
				break
			case "REMOVE":
				if err := remove(onereq.Name); err != nil {
					log.Println(err)
				}
				break
			}
		}
	}()
	return
}

func init() {
	go coreLoop() // run a long a long
}
