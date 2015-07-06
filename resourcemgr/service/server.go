package service

import (
	_ "github.com/yihungjen/bigobject-registry/resourcemgr/core"
	"net/http"
)

var (
	// Request endpoint multiplexer at PORT
	Server = http.NewServeMux()
)

func createResource(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowd", 405)
		return
	}
	http.Error(w, "Not yet implemented", 503)
}

func removeResource(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowd", 405)
		return
	}
	http.Error(w, "Not yet implemented", 503)
}

func init() {
	Server.HandleFunc("/create/", createResource)
	Server.HandleFunc("/remove/", removeResource)
}
